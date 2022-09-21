package tracks

import (
	"testing"
)

func TestTrackConstruction(t *testing.T) {
	ts := NewTrackSystem()
	l := ts.Layers[""]
	track := NewG1(l, 1)
	track.
		Connect(NewG1(l, 2)).
		Connect(NewG1(l, 3)).
		Connect(NewG1(l, 4)).
		Connect(NewR6Left(l, 5)).
		Connect(NewR6Left(l, 6)).
		Connect(NewR6Left(l, 7)).
		Connect(NewR6Right(l, 5)).
		Connect(NewR6Right(l, 6)).
		Connect(NewR6Right(l, 7))
}

func TestTrackConstruction2(t *testing.T) {
	InitRoco()
	ts := NewTrackSystem()
	l := ts.Layers[""]
	track := l.NewTrack("G1")
	if track == nil {
		t.Fatal("Unknown track type")
	}
	track.Connection(0).AddMark("Start")
	track2 := l.NewTrack("G1")
	if track2 == nil {
		t.Fatal("Unknown track type")
	}
	track.Connection(1).AddMark("End")
	track2.AddMark(0, "start")
	track2.AddMark(1, "end")
	track.Connect(track2)
}
