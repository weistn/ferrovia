package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/weistn/ferrovia/errlog"
	"github.com/weistn/ferrovia/interpreter"
	"github.com/weistn/ferrovia/parser"
	"github.com/weistn/ferrovia/tracks"
	"github.com/weistn/ferrovia/view2d"
	"github.com/weistn/goui"
)

/*
var data string = `{"statements":[
	{"way":{
		"exp": [
			{"anchor":{"x": 100, "y": 100, "z": 0, "a": 120}},
			{"con":{"name": "A"}},
			{"track":{"t": "G1"}},
			{"track":{"t": "R6"}},
			{"track":{"t": "G1"}},
			{"track":{"t": "R5_L"}},
			{"track":{"t": "G1"}},
			{"track":{"t": "W15R", "j1":[{"track":{"t": "DG1"}}, {"track":{"t": "DG1"}}, {"track":{"t": "R10_L"}}, {"track":{"t": "G1"}}]}},
			{"track":{"t": "W15L", "j1":[{"track":{"t": "DG1"}}, {"track":{"t": "DG1"}}, {"track":{"t": "R10"}}]}},
			{"track":{"t": "G1"}},
			{"track":{"t": "G1"}},
			{"con":{"name": "B"}}
		]
	}}
]}`

var data2 string = `
railway (
    "Einfahrt"
    @(120 mm, 0 mm, 0 mm, 180 deg)
    G1
    G1
    "Ausfahrt"
)

railway (
	"Ausfahrt"
    3 * L6
	WL15 /-> G4
	G1 ->/ WL15
	"weiter" <-/ WL15
	WL10 /-> G1
	G1 <-/ WR10
)

railway (
	"weiter"
	G1
	G1 ->/ DKW15 /-> G1
	G1 ->/ DKW10 /-> G1
)

railway (
    @(900 mm, 0 mm, 0 mm, 180 deg)
	BWR5 /-> 2 * R5
	2 * R5
)

railway (
    @(1500 mm, 0 mm, 0 mm, 180 deg)
	BWL5 /-> 2 * L5
	2 * L5
)

railway (
    @(1200 mm, 0 mm, 0 mm, 180 deg)
	L5
	BWL5 /-> "Innen"
	L6
)

railway (
	L5
	L5
	"Innen" ->/ BWR5
)

railway (
    @(2000 mm, 0 mm, 0 mm, 180 deg)
	2 * L9
	BWL9 /-> "Innen2"
	2 * L10
)

railway (
	4 * L9
	"Innen2" ->/ BWR9
)

ground {
    Top: 0 mm
    Left: 0 mm
    Width: 200 cm
    Height: 194 cm
}

ground {
    Top: 100 mm
    Left: 230 cm
    Width: 250 cm
    Height: 184 cm
}

ground {
	Top: 0 cm
	Left: 200 cm
	Polygon: (0 cm, 0 cm) (30 cm, 10 cm) (270 cm, 10 cm) (270 cm, 194 cm) (0 cm, 194 cm)
}
`
*/

type ServerAPI struct {
}

var server *goui.WebUI

func loadFile(name string) (*tracks.TrackSystem, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	log := errlog.NewErrorLog()
	fileId := log.AddFile(errlog.NewSourceFile(name))
	p := parser.NewParser(log)
	file := p.Parse(fileId, string(data))
	if log.HasErrors() {
		log.Print()
		return nil, errors.New("parsing Error")
	}

	b := interpreter.NewInterpreter(log)
	ts := b.Process(file)
	if log.HasErrors() {
		log.Print()
		return nil, errors.New("interpreter error")
	}

	return ts, nil
}

func showFile(filename string) error {
	ts, err := loadFile(filename)
	if err != nil {
		return err
	}
	ts.Name = "Demo"
	canvas := view2d.Render(ts)

	if err := server.SendEvent("canvas", canvas); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return err
	}
	return nil
}

func main() {
	//
	// Parse command lines
	//
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprint(os.Stderr, "Missing command line argument\n")
		return
	}
	filename := flag.Arg(0)

	tracks.InitRoco()

	/*file, err := parser.Load([]byte(data))
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}*/
	/*
		log := errlog.NewErrorLog()
		fileId := log.AddFile(errlog.NewSourceFile("data2"))
		p := parser.NewParser(log)
		file := p.Parse(fileId, data2)
		if log.HasErrors() {
			log.Print()
			return
		}

		b := interpreter.NewInterpreter(log)
		ts := b.Process(file)
		if log.HasErrors() {
			log.Print()
			return
		} */
	// ts, err := loadFile(filename)

	//
	// Open UI in browser
	//

	mime.AddExtensionType(".css", "text/css")


	remote := &ServerAPI{}

	server = goui.NewWebUI("/", remote, nil)
	server.Handle("/", http.FileServer(http.Dir("./ui")))
	err := server.Start()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}

	//
	// Watch file
	//
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprint(os.Stderr, "Could not watch file", err)
		return
	}
	defer watcher.Close()
	err = watcher.Add(filename)
	if err != nil {
		fmt.Fprint(os.Stderr, "Could not watch file "+filename, err)
	}

	go func() {
		showFile(filename)
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// fmt.Fprintf(os.Stderr, "%s %s\n", event.Name, event.Op)
				if event.Op == fsnotify.Create || event.Op == fsnotify.Write {
					showFile(filename)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			}
		}

	}()

	server.Wait()
}
