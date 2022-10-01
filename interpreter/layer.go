package interpreter

import (
	"github.com/weistn/ferrovia/errlog"
	"github.com/weistn/ferrovia/model/tracks"
)

// Implements IContext
type LayerContext struct {
	layer *tracks.TrackLayer
}

// Returns (nil, false, nil) if the call could not be made, because the function is unknown.
func (c *LayerContext) Call(b *Interpreter, loc errlog.LocationRange, name string, args ...*ExprValue) (*ExprValue, bool, *errlog.Error) {
	var err *errlog.Error
	switch name {
	case "name":
		if len(args) != 1 {
			return nil, true, errlog.NewError(errlog.ErrorArgumentCountMismatch, loc, "1")
		}
		c.layer.Name, err = args[0].ToString(loc)
		return nil, true, err
	case "color":
		if len(args) != 1 {
			return nil, true, errlog.NewError(errlog.ErrorArgumentCountMismatch, loc, "1")
		}
		c.layer.Color, err = args[0].ToString(loc)
		return nil, true, err
	}
	return nil, false, nil
}
