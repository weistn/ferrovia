package interpreter

import (
	"github.com/weistn/ferrovia/errlog"
	"github.com/weistn/ferrovia/parser"
)

func evalExpression(b *Interpreter, ast parser.IExpression) (*ExprValue, *errlog.Error) {
	switch t := ast.(type) {
	case *parser.BinaryExpression:
		return evalBinaryExpression(b, t)
	case *parser.DimensionExpression:
		return evalDimensionExpression(b, t)
	default:
		panic("TODO")
	}
}

func evalDimensionExpression(b *Interpreter, ast *parser.DimensionExpression) (*ExprValue, *errlog.Error) {
	left, err := evalExpression(b, ast.Value)
	if err != nil {
		return nil, err
	}
	switch ast.Dimension.StringValue {
	case "mm", "deg":
		return left, nil
	case "cm":
		if left.Type == numberType {
			result := &ExprValue{Type: numberType}
			result.NumberValue *= 10
		} else {
			return nil, errlog.NewError(errlog.ErrorTypeMismtach, ast.Dimension.Location)
		}
	case "m":
		if left.Type == numberType {
			result := &ExprValue{Type: numberType}
			result.NumberValue *= 1000
		} else {
			return nil, errlog.NewError(errlog.ErrorTypeMismtach, ast.Dimension.Location)
		}
	}
	panic("Oooops")
}

func evalBinaryExpression(b *Interpreter, ast *parser.BinaryExpression) (*ExprValue, *errlog.Error) {
	left, err := evalExpression(b, ast.Left)
	if err != nil {
		return nil, err
	}
	switch ast.Op.Kind {
	case parser.TokenLogicalAnd:
		b, err := left.ToBool(ast.Op.Location)
		if err != nil {
			return nil, err
		}
		if !b {
			return &ExprValue{Type: numberType, NumberValue: 0}, nil
		}
	case parser.TokenLogicalOr:
		b, err := left.ToBool(ast.Op.Location)
		if err != nil {
			return nil, err
		}
		if b {
			return &ExprValue{Type: numberType, NumberValue: 1}, nil
		}
	}

	right, err := evalExpression(b, ast.Right)
	if err != nil {
		return nil, err
	}

	switch ast.Op.Kind {
	case parser.TokenLogicalAnd:
		return left.LogicalAnd(right, ast.Op.Location)
	case parser.TokenLogicalOr:
		return left.LogicalOr(right, ast.Op.Location)
	case parser.TokenEqual:
		return left.Equal(right, ast.Op.Location)
	case parser.TokenNotEqual:
		return left.NotEqual(right, ast.Op.Location)
	case parser.TokenLessOrEqual:
		return left.LessOrEqual(right, ast.Op.Location)
	case parser.TokenGreaterOrEqual:
		return left.GreaterOrEqual(right, ast.Op.Location)
	case parser.TokenLess:
		return left.Less(right, ast.Op.Location)
	case parser.TokenGreater:
		return left.Greater(right, ast.Op.Location)
	case parser.TokenPlus:
		return left.Plus(right, ast.Op.Location)
	case parser.TokenDash:
		return left.Minus(right, ast.Op.Location)
	case parser.TokenAsterisk:
		return left.Mul(right, ast.Op.Location)
	case parser.TokenSlash:
		return left.Div(right, ast.Op.Location)
	case parser.TokenPercent:
		return left.Rem(right, ast.Op.Location)
	case parser.TokenAmpersand:
		return left.BinaryAnd(right, ast.Op.Location)
	case parser.TokenPipe:
		return left.BinaryOr(right, ast.Op.Location)
	case parser.TokenCaret:
		return left.BinaryXor(right, ast.Op.Location)
	case parser.TokenBitClear:
		return left.BinaryAndNot(right, ast.Op.Location)
	case parser.TokenShiftLeft:
		return left.Lsh(right, ast.Op.Location)
	case parser.TokenShiftRight:
		return left.Rsh(right, ast.Op.Location)
	}
	panic("Oooops")
}

func evalFloatExpression(b *Interpreter, ast parser.IExpression, loc errlog.LocationRange) (float64, *errlog.Error) {
	v, err := evalExpression(b, ast)
	if err != nil {
		return 0, err
	}
	return v.ToFloat(loc)
}

func evalBoolExpression(b *Interpreter, ast parser.IExpression, loc errlog.LocationRange) (bool, *errlog.Error) {
	v, err := evalExpression(b, ast)
	if err != nil {
		return false, err
	}
	return v.ToBool(loc)
}

func evalVectorExpression(b *Interpreter, ast parser.IExpression, loc errlog.LocationRange) ([]*ExprValue, *errlog.Error) {
	v, err := evalExpression(b, ast)
	if err != nil {
		return nil, err
	}
	if v.Type != vectorType {
		return nil, errlog.NewError(errlog.ErrorTypeMismtach, loc)
	}
	return v.VectorValue, nil
}
