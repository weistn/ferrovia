package interpreter

import (
	"testing"

	"github.com/weistn/ferrovia/errlog"
	"github.com/weistn/ferrovia/model/tracks"
	"github.com/weistn/ferrovia/parser"
)

var data string = `{"statements":[
	{"way":{
		"exp": [
			{"anchor":{"x": 100, "y": 100, "z": 0, "a": 0}},
			{"con":{"name": "A"}},
			{"track":{"t": "G1"}},
			{"track":{"t": "G4"}},
			{"con":{"name": "B"}}
		]
	}},
	{"way":{
		"exp": [
			{"con":{"name": "B"}},
			{"repeat":{"n": 6, "exp":[
				{"track":{"t": "R6"}}
			]}},
			{"track":{"t": "G1"}},
			{"track":{"t": "G1"}},
			{"repeat":{"n": 6, "exp":[
				{"track":{"t": "R6"}}
			]}},
			{"con":{"name": "A"}}
		]
	}}
]}`

/*
var data2 string = `{"statements":[
	{"way":{
		"exp": [
			{"anchor":{"x": 100, "y": 100, "z": 0, "a": 0}},
			{"con":{"name": "A"}},
			{"track":{"t": "G1"}},
			{"track":{"t": "G4"}},
			{"con":{"name": "B"}}
		]
	}}
]}`
*/

func TestInterpreter(t *testing.T) {
	tracks.InitRoco()
	file, err := parser.Load([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	e := errlog.NewErrorLog()
	b := NewInterpreter(e)
	ts, _ := b.Process(file)
	if ts == nil {
		t.Fatal("No track system")
	}
	if e.HasErrors() {
		e.Print()
		t.Fail()
	}
	if len(ts.Layers[""].Tracks) != 16 {
		println(len(ts.Layers[""].Tracks))
		t.Fatal("Tracks missing")
	}
	if ts.GetMark("A") == nil {
		t.Fatal("Missing mark A")
	}
	if ts.GetMark("B") == nil {
		t.Fatal("Missing mark B")
	}
	if !ts.GetMark("A").Connection.IsConnected() {
		t.Fatal("Mark A not connected")
	}
	if !ts.GetMark("B").Connection.IsConnected() {
		t.Fatal("Mark B not connected")
	}
}
