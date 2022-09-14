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

type WindowAPI struct {
}

var server *goui.Window

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

	//
	// Open UI in browser
	//

	mime.AddExtensionType(".css", "text/css")

	remote := &WindowAPI{}

	server = goui.NewWindow("/", remote, nil)
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
