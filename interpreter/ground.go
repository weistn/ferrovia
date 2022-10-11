package interpreter

import (
	"github.com/weistn/ferrovia/errlog"
	"github.com/weistn/ferrovia/model"
	"github.com/weistn/ferrovia/parser"
)

// Implements IContext
type GroundContext struct {
	Ground  *model.GroundPlate
	top     *FuncValue
	left    *FuncValue
	width   *FuncValue
	height  *FuncValue
	polygon *FuncValue
}

func NewGroundContext(ground *model.GroundPlate) *GroundContext {
	ctx := &GroundContext{Ground: ground}
	ctx.top = &FuncValue{
		Name: "top",
		Func: func(b *Interpreter, c []IContext, loc errlog.LocationRange, args ...parser.IExpression) (*ExprValue, *errlog.Error) {
			if len(args) != 1 {
				return nil, b.errlog.LogError(errlog.ErrorArgumentCountMismatch, loc, "1")
			}
			arg, err := b.evalExpression(c, args[0])
			if arg != nil {
				return nil, err
			}
			ctx.Ground.Top, err = b.ToFloat(arg, loc)
			return nil, err
		},
	}
	ctx.left = &FuncValue{
		Name: "left",
		Func: func(b *Interpreter, c []IContext, loc errlog.LocationRange, args ...parser.IExpression) (*ExprValue, *errlog.Error) {
			if len(args) != 1 {
				return nil, b.errlog.LogError(errlog.ErrorArgumentCountMismatch, loc, "1")
			}
			arg, err := b.evalExpression(c, args[0])
			if arg != nil {
				return nil, err
			}
			ctx.Ground.Left, err = b.ToFloat(arg, loc)
			return nil, err
		},
	}
	ctx.width = &FuncValue{
		Name: "width",
		Func: func(b *Interpreter, c []IContext, loc errlog.LocationRange, args ...parser.IExpression) (*ExprValue, *errlog.Error) {
			if len(args) != 1 {
				return nil, b.errlog.LogError(errlog.ErrorArgumentCountMismatch, loc, "1")
			}
			arg, err := b.evalExpression(c, args[0])
			if arg != nil {
				return nil, err
			}
			ctx.Ground.Width, err = b.ToFloat(arg, loc)
			return nil, err
		},
	}
	ctx.height = &FuncValue{
		Name: "height",
		Func: func(b *Interpreter, c []IContext, loc errlog.LocationRange, args ...parser.IExpression) (*ExprValue, *errlog.Error) {
			if len(args) != 1 {
				return nil, b.errlog.LogError(errlog.ErrorArgumentCountMismatch, loc, "1")
			}
			arg, err := b.evalExpression(c, args[0])
			if arg != nil {
				return nil, err
			}
			ctx.Ground.Height, err = b.ToFloat(arg, loc)
			return nil, err
		},
	}
	ctx.polygon = &FuncValue{
		Name: "polygon",
		Func: func(b *Interpreter, c []IContext, loc errlog.LocationRange, args ...parser.IExpression) (*ExprValue, *errlog.Error) {
			if len(args) < 3 {
				return nil, b.errlog.LogError(errlog.ErrorArgumentCountMismatch, loc, "1")
			}
			for _, argexpr := range args {
				arg, err := b.evalExpression(c, argexpr)
				if arg != nil {
					return nil, err
				}
				vector, err := b.ToVector(arg, loc)
				if err != nil {
					return nil, err
				}
				if len(vector) != 2 {
					return nil, b.errlog.LogError(errlog.ErrorArgumentCountMismatch, loc, "1")
				}
				x, err := b.ToFloat(vector[0], loc)
				if err != nil {
					return nil, err
				}
				y, err := b.ToFloat(vector[1], loc)
				if err != nil {
					return nil, err
				}
				ctx.Ground.Polygon = append(ctx.Ground.Polygon, model.GroundPoint{X: x, Y: y})
			}
			return nil, nil
		},
	}
	return ctx
}

func (c *GroundContext) Lookup(b *Interpreter, loc errlog.LocationRange, name string) (*ExprValue, *errlog.Error) {
	switch name {
	case "top":
		return &ExprValue{Type: funcType, FuncValue: c.top}, nil
	case "left":
		return &ExprValue{Type: funcType, FuncValue: c.left}, nil
	case "width":
		return &ExprValue{Type: funcType, FuncValue: c.width}, nil
	case "height":
		return &ExprValue{Type: funcType, FuncValue: c.height}, nil
	case "polygon":
		return &ExprValue{Type: funcType, FuncValue: c.polygon}, nil
	}
	return nil, nil
}
