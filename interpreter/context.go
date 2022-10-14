package interpreter

import (
	"github.com/weistn/ferrovia/errlog"
	"github.com/weistn/ferrovia/model/tracks"
)

type IContext interface {
	//	Call(b *Interpreter, loc errlog.LocationRange, name string, args ...*ExprValue) (*ExprValue, bool, *errlog.Error)
	Lookup(b *Interpreter, loc errlog.LocationRange, name string) (*ExprValue, *errlog.Error)
	Process(b *Interpreter, loc errlog.LocationRange, value *ExprValue) *errlog.Error
	Close(b *Interpreter) *errlog.Error
}

// Implements IContext
type GlobalContext struct {
	identifiers map[string]interface{}
}

// Implements IContext
type ValueContext struct {
	Value interface{}
}

func NewGlobalContext() *GlobalContext {
	return &GlobalContext{identifiers: make(map[string]interface{})}
}

func (ctx *GlobalContext) Lookup(b *Interpreter, loc errlog.LocationRange, name string) (*ExprValue, *errlog.Error) {
	ident, ok := ctx.identifiers[name]
	if ok {
		switch t := ident.(type) {
		case *TracksContext:
			return &ExprValue{Type: contextType, Context: t}, nil
		case *LayerContext:
			return &ExprValue{Type: contextType, Context: t}, nil
		default:
			panic("Ooooops")
		}
	}
	return nil, nil
}

func (ctx *GlobalContext) Process(b *Interpreter, loc errlog.LocationRange, value *ExprValue) *errlog.Error {
	return b.errlog.LogError(errlog.ErrorIllegalInThisContext, loc)
}

func (ctx *GlobalContext) Close(b *Interpreter) *errlog.Error {
	return nil
}

func (ctx *GlobalContext) RegisterTracks(b *Interpreter, loc errlog.LocationRange, name string) (*TracksContext, *errlog.Error) {
	ident, ok := ctx.identifiers[name]
	if ok {
		if t, ok := ident.(*TracksContext); ok {
			return t, nil
		}
		return nil, b.errlog.LogError(errlog.ErrorDuplicateIdentifier, loc, name)
	}
	t := NewTracksContext(b.model.Tracks.Layers[""])
	ctx.identifiers[name] = t
	return t, nil
}

func (ctx *GlobalContext) RegisterLayer(b *Interpreter, loc errlog.LocationRange, name string) (*LayerContext, *errlog.Error) {
	if _, ok := ctx.identifiers[name]; ok {
		return nil, b.errlog.LogError(errlog.ErrorDuplicateIdentifier, loc, name)
	}
	l := NewLayerContext(&tracks.TrackLayer{Name: name})
	ctx.identifiers[name] = l
	return l, nil
}

func (ctx *ValueContext) Lookup(b *Interpreter, loc errlog.LocationRange, name string) (*ExprValue, *errlog.Error) {
	return nil, nil
}

func (ctx *ValueContext) Process(b *Interpreter, loc errlog.LocationRange, value *ExprValue) *errlog.Error {
	return b.errlog.LogError(errlog.ErrorIllegalInThisContext, loc)
}

func (ctx *ValueContext) Close(b *Interpreter) *errlog.Error {
	return nil
}
