package parser

import (
	"testing"

	"github.com/weistn/ferrovia/errlog"
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

func TestParser(t *testing.T) {
	println(data)
	log := errlog.NewErrorLog()
	fileId := log.AddFile(errlog.NewSourceFile("data"))
	p := NewParser(log)
	p.Parse(fileId, data)
	if log.HasErrors() {
		log.Print()
		t.Fail()
	}
}
