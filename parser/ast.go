package parser

import (
	"github.com/weistn/ferrovia/errlog"
)

type File struct {
	Statements []IStatement
	Location   errlog.Location
}

type IStatement interface {
}

// Implements IStatement
type GroundPlate struct {
	Expressions []IExpression
	Location    errlog.LocationRange
}

// Implements IStatement
type Switchboard struct {
	Name          *Token
	RawText       string
	LocationToken errlog.LocationRange
	LocationText  errlog.LocationRange
}

// Implements IStatement
type Layer struct {
	Name        *Token
	Expressions []IExpression
	Location    errlog.LocationRange
}

// Implements IStatement
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
type VectorExpression struct {
	Values   []IExpression
	Location errlog.LocationRange
}
