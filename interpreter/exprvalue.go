package interpreter

import (
	"math/big"

	"github.com/weistn/ferrovia/errlog"
)

type ExprValue struct {
	Type         IType
	StringValue  string
	IntegerValue *big.Int
	FloatValue   *big.Float
	BoolValue    bool
	VectorValue  []*ExprValue
}

func (e *ExprValue) LogicalOr(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != boolType || p.Type != boolType {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	return &ExprValue{Type: boolType, BoolValue: e.BoolValue || p.BoolValue}, nil
}

func (e *ExprValue) LogicalAnd(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != boolType && p.Type != boolType {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	return &ExprValue{Type: boolType, BoolValue: e.BoolValue || p.BoolValue}, nil
}

func (e *ExprValue) Equal(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != p.Type {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	result := &ExprValue{Type: boolType}
	switch e.Type {
	case stringType:
		result.BoolValue = (e.StringValue == p.StringValue)
	case intType:
		result.BoolValue = (e.IntegerValue.Cmp(p.IntegerValue) == 0)
	case floatType:
		result.BoolValue = (e.FloatValue.Cmp(p.FloatValue) == 0)
	case boolType:
		result.BoolValue = (e.BoolValue == p.BoolValue)
	default:
		panic("TODO")
	}
	return result, nil
}

func (e *ExprValue) NotEqual(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != p.Type {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	result := &ExprValue{Type: boolType}
	switch e.Type {
	case stringType:
		result.BoolValue = (e.StringValue != p.StringValue)
	case intType:
		result.BoolValue = (e.IntegerValue.Cmp(p.IntegerValue) != 0)
	case floatType:
		result.BoolValue = (e.FloatValue.Cmp(p.FloatValue) != 0)
	case boolType:
		result.BoolValue = (e.BoolValue != p.BoolValue)
	default:
		panic("TODO")
	}
	return result, nil
}

func (e *ExprValue) LessOrEqual(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != p.Type {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	result := &ExprValue{Type: boolType}
	switch e.Type {
	case stringType:
		result.BoolValue = (e.StringValue <= p.StringValue)
	case intType:
		result.BoolValue = (e.IntegerValue.Cmp(p.IntegerValue) <= 0)
	case floatType:
		result.BoolValue = (e.FloatValue.Cmp(p.FloatValue) <= 0)
	case boolType:
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	default:
		panic("TODO")
	}
	return result, nil
}

func (e *ExprValue) GreaterOrEqual(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != p.Type {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	result := &ExprValue{Type: boolType}
	switch e.Type {
	case stringType:
		result.BoolValue = (e.StringValue >= p.StringValue)
	case intType:
		result.BoolValue = (e.IntegerValue.Cmp(p.IntegerValue) >= 0)
	case floatType:
		result.BoolValue = (e.FloatValue.Cmp(p.FloatValue) >= 0)
	case boolType:
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	default:
		panic("TODO")
	}
	return result, nil
}

func (e *ExprValue) Less(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != p.Type {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	result := &ExprValue{Type: boolType}
	switch e.Type {
	case stringType:
		result.BoolValue = (e.StringValue < p.StringValue)
	case intType:
		result.BoolValue = (e.IntegerValue.Cmp(p.IntegerValue) < 0)
	case floatType:
		result.BoolValue = (e.FloatValue.Cmp(p.FloatValue) < 0)
	case boolType:
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	default:
		panic("TODO")
	}
	return result, nil
}

func (e *ExprValue) Greater(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != p.Type {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	result := &ExprValue{Type: boolType}
	switch e.Type {
	case stringType:
		result.BoolValue = (e.StringValue > p.StringValue)
	case intType:
		result.BoolValue = (e.IntegerValue.Cmp(p.IntegerValue) > 0)
	case floatType:
		result.BoolValue = (e.FloatValue.Cmp(p.FloatValue) > 0)
	case boolType:
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	default:
		panic("TODO")
	}
	return result, nil
}

func (e *ExprValue) Plus(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != p.Type {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	result := &ExprValue{Type: e.Type}
	switch e.Type {
	case stringType:
		result.StringValue = e.StringValue + p.StringValue
	case intType:
		result.IntegerValue = big.NewInt(0)
		result.IntegerValue.Add(e.IntegerValue, p.IntegerValue)
	case floatType:
		result.FloatValue = big.NewFloat(0)
		result.FloatValue.Add(e.FloatValue, p.FloatValue)
	case boolType:
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	default:
		panic("TODO")
	}
	return result, nil
}

func (e *ExprValue) Minus(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != p.Type {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	result := &ExprValue{Type: e.Type}
	switch e.Type {
	case stringType:
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	case intType:
		result.IntegerValue = big.NewInt(0)
		result.IntegerValue.Sub(e.IntegerValue, p.IntegerValue)
	case floatType:
		result.FloatValue = big.NewFloat(0)
		result.FloatValue.Sub(e.FloatValue, p.FloatValue)
	case boolType:
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	default:
		panic("TODO")
	}
	return result, nil
}

func (e *ExprValue) Mul(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != p.Type {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	result := &ExprValue{Type: e.Type}
	switch e.Type {
	case stringType:
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	case intType:
		result.IntegerValue = big.NewInt(0)
		result.IntegerValue.Mul(e.IntegerValue, p.IntegerValue)
	case floatType:
		result.FloatValue = big.NewFloat(0)
		result.FloatValue.Mul(e.FloatValue, p.FloatValue)
	case boolType:
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	default:
		panic("TODO")
	}
	return result, nil
}

func (e *ExprValue) Div(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != p.Type {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	result := &ExprValue{Type: e.Type}
	switch e.Type {
	case stringType:
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	case intType:
		result.IntegerValue = big.NewInt(0)
		result.IntegerValue.Div(e.IntegerValue, p.IntegerValue)
	case floatType:
		result.FloatValue = big.NewFloat(0)
		result.FloatValue.Quo(e.FloatValue, p.FloatValue)
	case boolType:
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	default:
		panic("TODO")
	}
	return result, nil
}

func (e *ExprValue) BinaryOr(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != intType || p.Type != intType {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	result := &ExprValue{}
	result.IntegerValue = big.NewInt(0)
	result.IntegerValue.Or(e.IntegerValue, p.IntegerValue)
	return result, nil
}

func (e *ExprValue) BinaryAnd(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != intType || p.Type != intType {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	result := &ExprValue{}
	result.IntegerValue = big.NewInt(0)
	result.IntegerValue.And(e.IntegerValue, p.IntegerValue)
	return result, nil
}

func (e *ExprValue) BinaryAndNot(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != intType || p.Type != intType {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	result := &ExprValue{}
	result.IntegerValue = big.NewInt(0)
	result.IntegerValue.AndNot(e.IntegerValue, p.IntegerValue)
	return result, nil
}

func (e *ExprValue) BinaryXor(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != intType || p.Type != intType {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	result := &ExprValue{}
	result.IntegerValue = big.NewInt(0)
	result.IntegerValue.Xor(e.IntegerValue, p.IntegerValue)
	return result, nil
}

func (e *ExprValue) Rem(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != intType || p.Type != intType {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	result := &ExprValue{}
	result.IntegerValue = big.NewInt(0)
	result.IntegerValue.Rem(e.IntegerValue, p.IntegerValue)
	return result, nil
}

func (e *ExprValue) Lsh(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != intType || p.Type != intType {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	result := &ExprValue{}
	result.IntegerValue = big.NewInt(0)
	result.IntegerValue.Lsh(e.IntegerValue, uint(p.IntegerValue.Uint64()))
	return result, nil
}

func (e *ExprValue) Rsh(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != intType || p.Type != intType {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	result := &ExprValue{}
	result.IntegerValue = big.NewInt(0)
	result.IntegerValue.Rsh(e.IntegerValue, uint(p.IntegerValue.Uint64()))
	return result, nil
}

func (e *ExprValue) ToFloat(loc errlog.LocationRange) (float64, *errlog.Error) {
	if e.Type == floatType {
		v, _ := e.FloatValue.Float64()
		return v, nil
	}
	if e.Type == intType {
		var f = big.NewFloat(0)
		f.SetInt(e.IntegerValue)
		v, _ := e.FloatValue.Float64()
		return v, nil
	}
	return 0, errlog.NewError(errlog.ErrorTypeMismtach, loc)
}
