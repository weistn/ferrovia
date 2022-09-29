package interpreter

import (
	"github.com/weistn/ferrovia/errlog"
	"github.com/weistn/ferrovia/model"
)

// Implements IContext
type GroundContext struct {
	Ground *model.GroundPlate
}

// Returns (nil, false, nil) if the call could not be made, because the function is unknown.
func (c *GroundContext) Call(b *Interpreter, loc errlog.LocationRange, name string, args ...*ExprValue) (*ExprValue, bool, *errlog.Error) {
	var err *errlog.Error
	switch name {
	case "top":
		if len(args) != 1 {
			b.errlog.AddError(errlog.NewError(errlog.ErrorArgumentCountMismatch, loc, "1"))
		}
		c.Ground.Top, err = args[0].ToFloat(loc)
		return nil, true, err
	case "left":
		if len(args) != 1 {
			b.errlog.AddError(errlog.NewError(errlog.ErrorArgumentCountMismatch, loc, "1"))
		}
		c.Ground.Left, err = args[0].ToFloat(loc)
		return nil, true, err
	case "width":
		if len(args) != 1 {
			b.errlog.AddError(errlog.NewError(errlog.ErrorArgumentCountMismatch, loc, "1"))
		}
		c.Ground.Width, err = args[0].ToFloat(loc)
		return nil, true, err
	case "height":
		if len(args) != 1 {
			b.errlog.AddError(errlog.NewError(errlog.ErrorArgumentCountMismatch, loc, "1"))
		}
		c.Ground.Height, err = args[0].ToFloat(loc)
		return nil, true, err
	case "polygon":
		if len(args) < 3 {
			b.errlog.AddError(errlog.NewError(errlog.ErrorArgumentCountMismatch, loc, "1"))
		}
		for _, arg := range args {
			vector, err := arg.ToVector(loc)
			if err != nil {
				return nil, true, err
			}
			if len(vector) != 2 {
				b.errlog.AddError(errlog.NewError(errlog.ErrorArgumentCountMismatch, loc, "1"))
			}
			x, err := vector[0].ToFloat(loc)
			if err != nil {
				return nil, true, err
			}
			y, err := vector[1].ToFloat(loc)
			if err != nil {
				return nil, true, err
			}
			c.Ground.Polygon = append(c.Ground.Polygon, model.GroundPoint{X: x, Y: y})
		}
		return nil, true, nil
	}
	return nil, false, nil
}
