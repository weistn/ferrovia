package parser

import (
	"github.com/weistn/ferrovia/errlog"
)

// Parser ...
type Parser struct {
	l          *Lexer
	savedToken *Token
	log        *errlog.ErrorLog
}

// NewParser creates a new parser that reports to the given error log.
func NewParser(log *errlog.ErrorLog) *Parser {
	return &Parser{log: log}
}

// Parse a file. For efficieny, files are identified by a number instead of their path.
func (p *Parser) Parse(fileId int, str string) *File {
	p.l = NewLexer(fileId, str, p.log)
	f := &File{Location: errlog.EncodeLocation(fileId, 0, 0)}
	p.parseFile(f)
	return f
}

func (p *Parser) parseFile(f *File) {
	for {
		t, err := p.expectMulti(TokenRailway, TokenGround, TokenEOF, TokenNewline)
		if err != nil {
			p.log.AddError(err)
			return
		}
		if t.Kind == TokenNewline {
			continue
		}
		if t.Kind == TokenRailway {
			rail, err := p.parseRailway(t)
			if err != nil {
				p.log.AddError(err)
				return
			}
			f.Statements = append(f.Statements, &Statement{RailWay: rail})
		} else if t.Kind == TokenLayer {
			l, err := p.parseLayer(t)
			if err != nil {
				p.log.AddError(err)
				return
			}
			f.Statements = append(f.Statements, &Statement{Layer: l})
		} else if t.Kind == TokenGround {
			ground, err := p.parseGround(t)
			if err != nil {
				p.log.AddError(err)
				return
			}
			f.Statements = append(f.Statements, &Statement{GroundPlate: ground})
		} else if t.Kind == TokenEOF {
			break
		} else {
			panic("Oooops")
		}
	}
}

func (p *Parser) parseLayer(t *Token) (*Layer, *errlog.Error) {
	l := &Layer{Location: t.Location}

	if _, ok := p.optional(TokenOpenBraces); ok {
		if _, err := p.expect(TokenNewline); err != nil {
			return nil, err
		}
		for {
			if _, ok := p.optional(TokenCloseBraces); ok {
				break
			}
			t, err := p.expect(TokenIdentifier)
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(TokenColon); err != nil {
				return nil, err
			}
			if t.StringValue == "Color" {
				// TODO
			} else {
				return nil, errlog.NewError(errlog.ErrorIllegalProperty, t.Location, t.StringValue)
			}
			if _, err := p.expect(TokenNewline); err != nil {
				return nil, err
			}
		}
	}
	return l, nil
}

func (p *Parser) parseRailway(t *Token) (*RailWay, *errlog.Error) {
	rail := &RailWay{Location: t.Location}

	if _, ok := p.optional(TokenOpenBraces); ok {
		if _, err := p.expect(TokenNewline); err != nil {
			return nil, err
		}
		for {
			if _, ok := p.optional(TokenCloseBraces); ok {
				break
			}
			t, err := p.expect(TokenIdentifier)
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(TokenColon); err != nil {
				return nil, err
			}
			if t.StringValue == "Layer" {
				t, err = p.expect(TokenIdentifier)
				if err != nil {
					return nil, err
				}
				rail.Layer = t.Raw
			} else {
				return nil, errlog.NewError(errlog.ErrorIllegalProperty, t.Location, t.StringValue)
			}
			if _, err := p.expect(TokenNewline); err != nil {
				return nil, err
			}
		}
	}
	if p.peek(TokenOpenParanthesis) {
		var err *errlog.Error
		rail.Expressions, err = p.parseExpression()
		return rail, err
	}
	return nil, p.expectError(TokenOpenParanthesis)
}

func (p *Parser) parseExpressionList(t *Token) ([]*Expression, *errlog.Error) {
	var expressions []*Expression
	p.optional(TokenNewline)
	for {
		if _, ok := p.optional(TokenCloseParanthesis); ok {
			break
		}
		exp, err := p.parseExpression()
		if err != nil {
			return expressions, err
		}
		expressions = append(expressions, exp...)
	}
	return expressions, nil
}

func (p *Parser) parseExpression() ([]*Expression, *errlog.Error) {
	t, err := p.expectMulti(TokenIdentifier, TokenAt, TokenString, TokenInteger, TokenOpenParanthesis)
	if err != nil {
		return nil, err
	}
	if t.Kind == TokenOpenParanthesis {
		exp, err := p.parseExpressionList(t)
		if err != nil {
			return nil, err
		}
		p.optional(TokenNewline)
		return p.parseLeftJunctionExpression(exp)
	}
	if t.Kind == TokenIdentifier {
		var junctionsRight []*JunctionExpression
		for p.peek(TokenArrowRightIn) || p.peek(TokenArrowRightOut) {
			j, err := p.parseRightJunctionExpression()
			if err != nil {
				return nil, err
			}
			junctionsRight = append(junctionsRight, j)
		}
		exp := &Expression{Track: &TrackExpression{Type: t.StringValue, JunctionsOnRight: junctionsRight, Location: t.Location}}
		if len(junctionsRight) == 0 {
			p.optional(TokenNewline)
			return p.parseLeftJunctionExpression([]*Expression{exp})
		}
		return []*Expression{exp}, nil
	} else if t.Kind == TokenString {
		p.optional(TokenNewline)
		exp := &Expression{ConnectionMark: &ConnectionMarkExpression{Name: t.StringValue, Location: t.Location}}
		return p.parseLeftJunctionExpression([]*Expression{exp})
	} else if t.Kind == TokenAt {
		if _, err := p.expect(TokenOpenParanthesis); err != nil {
			return nil, err
		}
		x, err := p.expectConstantWithLengthUnit()
		if err != nil {
			return nil, err
		}
		if _, err = p.expect(TokenComma); err != nil {
			return nil, err
		}
		y, err := p.expectConstantWithLengthUnit()
		if err != nil {
			return nil, err
		}
		if _, err = p.expect(TokenComma); err != nil {
			return nil, err
		}
		z, err := p.expectConstantWithLengthUnit()
		if err != nil {
			return nil, err
		}
		if _, err = p.expect(TokenComma); err != nil {
			return nil, err
		}
		a, err := p.expectConstantWithAngleUnit()
		if err != nil {
			return nil, err
		}
		t2, err := p.expect(TokenCloseParanthesis)
		if err != nil {
			return nil, err
		}
		p.optional(TokenNewline)
		exp := &Expression{Anchor: &AnchorExpression{X: x, Y: y, Z: z, Angle: a, Location: t.Location.Join(t2.Location)}}
		return []*Expression{exp}, nil
	} else if t.Kind == TokenInteger {
		if !t.IntegerValue.IsInt64() {
			return nil, errlog.NewError(errlog.ErrorIllegalNumber, t.Location)
		}
		count := t.IntegerValue.Int64()
		if count < 0 || count >= 100 {
			return nil, errlog.NewError(errlog.ErrorIllegalNumber, t.Location)
		}
		if _, err := p.expect(TokenAsterisk); err != nil {
			return nil, err
		}
		exp, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		p.optional(TokenNewline)
		rep := &Expression{Repeat: &RepeatExpression{Count: int(count), TrackExpressions: exp}}
		return p.parseLeftJunctionExpression([]*Expression{rep})
	}
	panic("Oooops")
}

func (p *Parser) parseRightJunctionExpression() (*JunctionExpression, *errlog.Error) {
	t, err := p.expectMulti(TokenArrowRightIn, TokenArrowRightOut)
	if err != nil {
		return nil, err
	}
	j := &JunctionExpression{Arrow: t.Kind}
	j.Expressions, err = p.parseExpression()
	return j, err
}

func (p *Parser) parseLeftJunctionExpression(exp []*Expression) ([]*Expression, *errlog.Error) {
	t, ok := p.optionalMulti(TokenArrowLeftIn, TokenArrowLeftOut)
	if !ok {
		p.optional(TokenNewline)
		return exp, nil
	}
	trackexp, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	if len(trackexp) != 1 || trackexp[0].Track == nil {
		return nil, p.log.LogError(errlog.ErrorIllegalJunction, t.Location)
	}
	j := &JunctionExpression{Arrow: t.Kind, Expressions: exp}
	trackexp[0].Track.JunctionsOnLeft = append([]*JunctionExpression{j}, trackexp[0].Track.JunctionsOnLeft...)
	return trackexp, err
}

func (p *Parser) parseGround(t *Token) (*GroundPlate, *errlog.Error) {
	if _, err := p.expect(TokenOpenBraces); err != nil {
		return nil, err
	}
	if _, err := p.expect(TokenNewline); err != nil {
		return nil, err
	}
	ground := &GroundPlate{Location: t.Location}
	for {
		if _, ok := p.optional(TokenCloseBraces); ok {
			break
		}
		t, err := p.expect(TokenIdentifier)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(TokenColon); err != nil {
			return nil, err
		}
		if t.StringValue == "Top" {
			ground.Top, err = p.expectConstantWithLengthUnit()
			if err != nil {
				return nil, err
			}
		} else if t.StringValue == "Left" {
			ground.Left, err = p.expectConstantWithLengthUnit()
			if err != nil {
				return nil, err
			}
		} else if t.StringValue == "Width" {
			ground.Width, err = p.expectConstantWithLengthUnit()
			if err != nil {
				return nil, err
			}
		} else if t.StringValue == "Height" {
			ground.Height, err = p.expectConstantWithLengthUnit()
			if err != nil {
				return nil, err
			}
		} else if t.StringValue == "Polygon" {
			ground.Polygon, err = p.parseGroundPolygon()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errlog.NewError(errlog.ErrorIllegalProperty, t.Location, t.StringValue)
		}
		if _, err := p.expect(TokenNewline); err != nil {
			return nil, err
		}
	}
	return ground, nil
}

func (p *Parser) parseGroundPolygon() ([]GroundPoint, *errlog.Error) {
	var points []GroundPoint
	for {
		var pnt GroundPoint
		if _, ok := p.optional(TokenOpenParanthesis); !ok {
			return points, nil
		}
		var err *errlog.Error
		pnt.X, err = p.expectConstantWithLengthUnit()
		if err != nil {
			return points, err
		}
		_, err = p.expect(TokenComma)
		if err != nil {
			return points, nil
		}
		pnt.Y, err = p.expectConstantWithLengthUnit()
		if err != nil {
			return points, err
		}
		_, err = p.expect(TokenCloseParanthesis)
		if err != nil {
			return points, nil
		}
		points = append(points, pnt)
	}
}

func (p *Parser) expect(tokenKind TokenKind) (*Token, *errlog.Error) {
	t := p.scan()
	if t.Kind != tokenKind {
		err := p.log.LogError(errlog.ErrorExpectedToken, t.Location, t.Raw, p.l.TokenKindToString(tokenKind))
		return nil, err
	}
	return t, nil
}

func (p *Parser) expectMulti(tokenKind ...TokenKind) (*Token, *errlog.Error) {
	t := p.scan()
	for _, k := range tokenKind {
		if t.Kind == k {
			return t, nil
		}
	}
	var str = []string{t.Raw}
	for _, kind := range tokenKind {
		str = append(str, p.l.TokenKindToString(kind))
	}
	err := p.log.LogError(errlog.ErrorExpectedToken, t.Location, str...)
	return nil, err
}

func (p *Parser) optional(tokenKind TokenKind) (*Token, bool) {
	t := p.scan()
	if t.Kind != tokenKind {
		p.savedToken = t
		return nil, false
	}
	return t, true
}

func (p *Parser) optionalMulti(tokenKind ...TokenKind) (*Token, bool) {
	t := p.scan()
	for _, k := range tokenKind {
		if t.Kind == k {
			return t, true
		}
	}
	p.savedToken = t
	return nil, false
}

func (p *Parser) peek(tokenKind TokenKind) bool {
	t := p.scan()
	p.savedToken = t
	return t.Kind == tokenKind
}

func (p *Parser) scan() *Token {
	if p.savedToken != nil {
		t := p.savedToken
		p.savedToken = nil
		return t
	}
	return p.l.Scan()
}

func (p *Parser) expectError(tokenKind ...TokenKind) *errlog.Error {
	t := p.scan()
	var str = []string{t.Raw}
	for _, kind := range tokenKind {
		str = append(str, p.l.TokenKindToString(kind))
	}
	err := p.log.LogError(errlog.ErrorExpectedToken, t.Location, str...)
	return err
}

func (p *Parser) expectConstantWithLengthUnit() (value float64, err *errlog.Error) {
	t, err := p.expectMulti(TokenInteger, TokenFloat)
	if err != nil {
		return 0, err
	}
	if t.Kind == TokenInteger {
		if !t.IntegerValue.IsInt64() {
			return 0, errlog.NewError(errlog.ErrorIllegalNumber, t.Location)
		}
		value = float64(t.IntegerValue.Int64())
	} else {
		value, _ = t.FloatValue.Float64()
	}
	unit, err := p.expect(TokenUnit)
	if err != nil {
		return 0, nil
	}
	if unit.StringValue == "mm" {
		// Do nothing
	} else if unit.StringValue == "cm" {
		value *= 10
	} else if unit.StringValue == "m" {
		value *= 1000
	} else {
		return 0, errlog.NewError(errlog.ErrorIllegalUnit, unit.Location, unit.StringValue)
	}
	return
}

func (p *Parser) expectConstantWithAngleUnit() (value float64, err *errlog.Error) {
	t, err := p.expectMulti(TokenInteger, TokenFloat)
	if err != nil {
		return 0, err
	}
	if t.Kind == TokenInteger {
		if !t.IntegerValue.IsInt64() {
			return 0, errlog.NewError(errlog.ErrorIllegalNumber, t.Location)
		}
		value = float64(t.IntegerValue.Int64())
	} else {
		value, _ = t.FloatValue.Float64()
	}
	unit, err := p.expect(TokenUnit)
	if err != nil {
		return 0, nil
	}
	if unit.StringValue == "deg" {
		// Do nothing
	} else {
		return 0, errlog.NewError(errlog.ErrorIllegalUnit, unit.Location, unit.StringValue)
	}
	return
}
