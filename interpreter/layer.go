package interpreter

import (
	"github.com/weistn/ferrovia/errlog"
	"github.com/weistn/ferrovia/model/tracks"
	"github.com/weistn/ferrovia/parser"
)

// Implements IContext
type LayerContext struct {
	layer *tracks.TrackLayer
	color *FuncValue
}

func NewLayerContext(layer *tracks.TrackLayer) *LayerContext {
	ctx := &LayerContext{layer: layer}
	ctx.color = &FuncValue{
		Name: "color",
		Func: func(b *Interpreter, c []IContext, loc errlog.LocationRange, args ...parser.IExpression) (*ExprValue, *errlog.Error) {
			if len(args) != 1 {
				return nil, b.errlog.LogError(errlog.ErrorArgumentCountMismatch, loc, "1")
			}
			arg, err := b.evalExpression(c, args[0])
			if arg != nil {
				return nil, err
			}
			ctx.layer.Color, err = b.ToString(arg, loc)
			return nil, err
		},
	}
	return ctx
}

func (c *LayerContext) Lookup(b *Interpreter, loc errlog.LocationRange, name string) (*ExprValue, *errlog.Error) {
	switch name {
	case "color":
		return &ExprValue{Type: funcType, FuncValue: c.color}, nil
	}
	return nil, nil
}

func (ctx *LayerContext) Process(b *Interpreter, loc errlog.LocationRange, value *ExprValue) *errlog.Error {
	return b.errlog.LogError(errlog.ErrorIllegalInThisContext, loc)
}

func (c *LayerContext) Close(b *Interpreter) *errlog.Error {
	return nil
}
