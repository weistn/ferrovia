package interpreter

import (
	"github.com/weistn/ferrovia/errlog"
	"github.com/weistn/ferrovia/model/tracks"
)

// Implements IContext
type TrackContext struct {
	track *tracks.Track
	first *tracks.TrackConnection
	last  *tracks.TrackConnection
}

// Returns (nil, false, nil) if the call could not be made, because the function is unknown.
func (c *TrackContext) Call(b *Interpreter, loc errlog.LocationRange, name string, args ...*ExprValue) (*ExprValue, bool, *errlog.Error) {
	panic("TODO")
}
