package interpreter

import "github.com/weistn/ferrovia/errlog"

type IContext interface {
	Call(b *Interpreter, loc errlog.LocationRange, name string, args ...*ExprValue) (*ExprValue, bool, *errlog.Error)
}
