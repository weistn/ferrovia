package interpreter

type IType interface {
}

type BasicType struct {
	Name string
}

var numberType = &BasicType{Name: "number"}
var stringType = &BasicType{Name: "string"}
var vectorType = &BasicType{Name: "vector"}
