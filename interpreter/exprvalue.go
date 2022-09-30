package interpreter

import (
	"github.com/weistn/ferrovia/errlog"
)

type ExprValue struct {
	Type        IType
	StringValue string
	NumberValue float64
	VectorValue []*ExprValue
	Context     IContext
}

func (e *ExprValue) LogicalOr(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != numberType || p.Type != numberType {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	if (e.NumberValue != 0) || (p.NumberValue != 0) {
		return &ExprValue{Type: numberType, NumberValue: 1}, nil
	}
	return &ExprValue{Type: numberType, NumberValue: 0}, nil
}

func (e *ExprValue) LogicalAnd(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != numberType && p.Type != numberType {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	if (e.NumberValue != 0) && (p.NumberValue != 0) {
		return &ExprValue{Type: numberType, NumberValue: 1}, nil
	}
	return &ExprValue{Type: numberType, NumberValue: 0}, nil
}

func (e *ExprValue) Equal(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != p.Type {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	result := &ExprValue{Type: numberType}
	switch e.Type {
	case stringType:
		if e.StringValue == p.StringValue {
			result.NumberValue = 1
		} else {
			result.NumberValue = 0
		}
	case numberType:
		if e.NumberValue == p.NumberValue {
			result.NumberValue = 1
		} else {
			result.NumberValue = 0
		}
	case vectorType:
		panic("TODO")
	default:
		panic("TODO")
	}
	return result, nil
}

func (e *ExprValue) NotEqual(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != p.Type {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	result := &ExprValue{Type: numberType}
	switch e.Type {
	case stringType:
		if e.StringValue != p.StringValue {
			result.NumberValue = 1
		} else {
			result.NumberValue = 0
		}
	case numberType:
		if e.NumberValue != p.NumberValue {
			result.NumberValue = 1
		} else {
			result.NumberValue = 0
		}
	case vectorType:
		panic("TODO")
	default:
		panic("TODO")
	}
	return result, nil
}

func (e *ExprValue) LessOrEqual(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != p.Type {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	result := &ExprValue{Type: numberType}
	switch e.Type {
	case stringType:
		if e.StringValue <= p.StringValue {
			result.NumberValue = 1
		} else {
			result.NumberValue = 0
		}
	case numberType:
		if e.NumberValue <= p.NumberValue {
			result.NumberValue = 1
		} else {
			result.NumberValue = 0
		}
	case vectorType:
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
	result := &ExprValue{Type: numberType}
	switch e.Type {
	case stringType:
		if e.StringValue >= p.StringValue {
			result.NumberValue = 1
		} else {
			result.NumberValue = 0
		}
	case numberType:
		if e.NumberValue >= p.NumberValue {
			result.NumberValue = 1
		} else {
			result.NumberValue = 0
		}
	case vectorType:
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
	result := &ExprValue{Type: numberType}
	switch e.Type {
	case stringType:
		if e.StringValue < p.StringValue {
			result.NumberValue = 1
		} else {
			result.NumberValue = 0
		}
	case numberType:
		if e.NumberValue < p.NumberValue {
			result.NumberValue = 1
		} else {
			result.NumberValue = 0
		}
	case vectorType:
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
	result := &ExprValue{Type: numberType}
	switch e.Type {
	case stringType:
		if e.StringValue > p.StringValue {
			result.NumberValue = 1
		} else {
			result.NumberValue = 0
		}
	case numberType:
		if e.NumberValue > p.NumberValue {
			result.NumberValue = 1
		} else {
			result.NumberValue = 0
		}
	case vectorType:
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
	case numberType:
		result.NumberValue = e.NumberValue + p.NumberValue
	case vectorType:
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
	case numberType:
		result.NumberValue = e.NumberValue - p.NumberValue
	case vectorType:
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
	case numberType:
		result.NumberValue = e.NumberValue * p.NumberValue
	case vectorType:
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
	case numberType:
		result.NumberValue = e.NumberValue / p.NumberValue
	case vectorType:
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	default:
		panic("TODO")
	}
	return result, nil
}

func (e *ExprValue) BinaryOr(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != numberType || p.Type != numberType {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	return &ExprValue{Type: numberType, NumberValue: float64(uint64(e.NumberValue) | uint64(p.NumberValue))}, nil
}

func (e *ExprValue) BinaryAnd(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != numberType || p.Type != numberType {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	return &ExprValue{Type: numberType, NumberValue: float64(uint64(e.NumberValue) & uint64(p.NumberValue))}, nil
}

func (e *ExprValue) BinaryAndNot(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != numberType || p.Type != numberType {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	return &ExprValue{Type: numberType, NumberValue: float64(uint64(e.NumberValue) &^ uint64(p.NumberValue))}, nil
}

func (e *ExprValue) BinaryXor(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != numberType || p.Type != numberType {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	return &ExprValue{Type: numberType, NumberValue: float64(uint64(e.NumberValue) ^ uint64(p.NumberValue))}, nil
}

func (e *ExprValue) Rem(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != numberType || p.Type != numberType {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	return &ExprValue{Type: numberType, NumberValue: float64(uint64(e.NumberValue) % uint64(p.NumberValue))}, nil
}

func (e *ExprValue) Lsh(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != numberType || p.Type != numberType {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	return &ExprValue{Type: numberType, NumberValue: float64(uint64(e.NumberValue) << uint64(p.NumberValue))}, nil
}

func (e *ExprValue) Rsh(p *ExprValue, loc errlog.LocationRange) (*ExprValue, *errlog.Error) {
	if e.Type != numberType || p.Type != numberType {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	return &ExprValue{Type: numberType, NumberValue: float64(uint64(e.NumberValue) >> uint64(p.NumberValue))}, nil
}

func (e *ExprValue) ToFloat(loc errlog.LocationRange) (float64, *errlog.Error) {
	if e.Type == numberType {
		return e.NumberValue, nil
	}
	return 0, errlog.NewError(errlog.ErrorTypeMismtach, loc)
}

func (e *ExprValue) ToBool(loc errlog.LocationRange) (bool, *errlog.Error) {
	if e.Type == numberType {
		return e.NumberValue != 0, nil
	}
	if e.Type == stringType {
		return e.StringValue != "", nil
	}
	if e.Type == vectorType {
		return len(e.VectorValue) != 0, nil
	}
	return false, errlog.NewError(errlog.ErrorTypeMismtach, loc)
}

func (e *ExprValue) ToVector(loc errlog.LocationRange) ([]*ExprValue, *errlog.Error) {
	if e.Type == vectorType {
		return e.VectorValue, nil
	}
	return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
}

func (e *ExprValue) ToString(loc errlog.LocationRange) (string, *errlog.Error) {
	if e.Type == stringType {
		return e.StringValue, nil
	}
	return "", errlog.NewError(errlog.ErrorTypeMismtach, loc)
}
