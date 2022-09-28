package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/weistn/ferrovia/errlog"
	"github.com/weistn/ferrovia/interpreter"
	"github.com/weistn/ferrovia/model"
	"github.com/weistn/ferrovia/model/tracks"
	"github.com/weistn/ferrovia/parser"
	"github.com/weistn/ferrovia/view/switchboard"
	"github.com/weistn/ferrovia/view/tracks2d"
	"github.com/weistn/goui"

	"embed"
)

type WindowAPI struct {
}

/*
var demodata string = `{
"tracks": [
	{"c": 5, "r": 5, "kind": 20},
	{"c": 6, "r": 5, "kind": 20},
	{"c": 7, "r": 5, "kind": 20}
],
"columns": 20,
"rows": 20
}`
*/

//go:embed view/*.html view/*.css view/*.js view/fonts view/switchboard/*.js view/switchboard/*.css view/tracks2d/*.js view/tracks2d/*.css
var uiFS embed.FS

var window *goui.Window

func loadFile(name string) (*model.Model, error) {
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
	m := b.ProcessStatics(file)
	if log.HasErrors() {
		log.Print()
		return nil, errors.New("interpreter error")
	}

	return m, nil
}

func showFile(filename string) error {
	model, err := loadFile(filename)
	if err != nil {
		return err
	}
	model.Name = "Demo"

	canvas := tracks2d.Render(model)
	if err := window.SendEvent("canvas", canvas); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return err
	}

	layout := switchboard.Render(model.Switchboards)
	if err = window.SendEvent("layout", layout); err != nil {
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

	//
	// Open UI in browser
	//

	remote := &WindowAPI{}

	window = goui.NewWindow("/", remote, nil)
	subfs, err := fs.Sub(uiFS, "view")
	if err != nil {
		panic("Embedding failed")
	}
	window.Handle("/", http.FileServer(http.FS(subfs)))
	err = window.Start()
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

	window.Wait()
}
