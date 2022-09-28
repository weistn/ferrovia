package interpreter

type IContext interface {
	Call(b *Interpreter, name string, args ...*ExprValue) *ExprValue
}
