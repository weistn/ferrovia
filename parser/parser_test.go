package parser

import (
	"testing"

	"github.com/weistn/ferrovia/errlog"
)

var data string = `
railway {
    "Einfahrt"
    @(120 mm, 120 mm, 0 mm, 180 deg)
    G1
    G1
    3 * R6
    "Ausfahrt"
}

ground {
    Top: 0 mm
    Left: 0 mm
    Width: 470 cm
    Height: 194 cm
}`

func TestParser(t *testing.T) {
	log := errlog.NewErrorLog()
	fileId := log.AddFile(errlog.NewSourceFile("data"))
	p := NewParser(log)
	p.Parse(fileId, data)
	if log.HasErrors() {
		log.Print()
		t.Fail()
	}
}
