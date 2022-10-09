package parser

import (
	"github.com/weistn/ferrovia/errlog"
)

type File struct {
	Statements []IDirective
	Location   errlog.Location
}

type IDirective interface {
}

// Implements IDirective
type GroundPlate struct {
	Expressions []IExpression
	Location    errlog.LocationRange
}

// Implements IDirective
type Switchboard struct {
	Name          *Token
	RawText       string
	LocationToken errlog.LocationRange
	LocationText  errlog.LocationRange
}

// Implements IDirective
type Layer struct {
	Name        *Token
	Expressions []IExpression
	Location    errlog.LocationRange
}

// Implements IDirective
type Tracks struct {
	Name        *Token
	Parameters  []*Parameter
	Expressions []IExpression
	Location    errlog.LocationRange
}

type Parameter struct {
	Name *Token
}

type IExpression interface {
}

// Implements IExpression
type BinaryExpression struct {
	Left  IExpression
	Op    *Token
	Right IExpression
}

// Implements IExpression
type IdentifierExpression struct {
	Identifier *Token
}

// Implements IExpression
type ConstantExpression struct {
	Value *Token
}

// Implements IExpression
type DimensionExpression struct {
	Value     IExpression
	Dimension *Token
}

// Implements IExpression
type CallExpression struct {
	Func      IExpression
	Arguments []IExpression
}

// Implements IExpression
type DotExpression struct {
	Context    IExpression
	Identifier *Token
	Arguments  []IExpression
}

type ContextExpression struct {
	Object     IExpression
	Statements []IExpression
}

// Implements IExpression
type VectorExpression struct {
	Values   []IExpression
	Location errlog.LocationRange
}
