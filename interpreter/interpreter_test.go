package interpreter

import (
	"testing"

	"github.com/weistn/ferrovia/errlog"
	"github.com/weistn/ferrovia/model/tracks"
	"github.com/weistn/ferrovia/parser"
)

var data string = `
tracks {
	mountain
    Ausfahrt
    @(120 mm, 120 mm, 0 mm, 180 deg)
    G1
    G1
    3 * R6
    ` + "`" + `Gleis 1` + "`" + `
}

tracks Spindel(radius) {
	6 * radius
}

layer mountain {
	color("red")
}

ground {
	top(0 cm)
	left(0 cm)
	polygon([530 cm, 335 cm], [530 cm, 385 cm], [0 cm, 385 cm], [0 cm, 335 cm])
}`

func TestInterpreter(t *testing.T) {
	tracks.InitRoco()
	e := errlog.NewErrorLog()

	// Parse
	fileId := e.AddFile(errlog.NewSourceFile("data"))
	p := parser.NewParser(e)
	file := p.Parse(fileId, data)
	if e.HasErrors() {
		e.Print()
		t.Fail()
	}

	// Interpret
	b := NewInterpreter(e)
	model := b.Process(file)
	if model == nil {
		t.Fatal("No model")
	}
	if e.HasErrors() {
		e.Print()
		t.Fail()
	}
	if len(model.Tracks.Layers[""].Tracks) != 16 {
		println(len(model.Tracks.Layers[""].Tracks))
		t.Fatal("Tracks missing")
	}
	if model.Tracks.GetMark("A") == nil {
		t.Fatal("Missing mark A")
	}
	if model.Tracks.GetMark("B") == nil {
		t.Fatal("Missing mark B")
	}
	if !model.Tracks.GetMark("A").Connection.IsConnected() {
		t.Fatal("Mark A not connected")
	}
	if !model.Tracks.GetMark("B").Connection.IsConnected() {
		t.Fatal("Mark B not connected")
	}
}
