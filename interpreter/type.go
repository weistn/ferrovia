package interpreter

type IType interface {
}

type BasicType struct {
	Name string
}

var intType = &BasicType{Name: "int"}
var floatType = &BasicType{Name: "float"}
var boolType = &BasicType{Name: "bool"}
var stringType = &BasicType{Name: "string"}
var vectorType = &BasicType{Name: "vector"}
