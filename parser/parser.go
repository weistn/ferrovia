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
		t, err := p.expectMulti(TokenRailway, TokenGround, TokenEOF, TokenNewline, TokenIdentifier)
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
		} else if t.Kind == TokenIdentifier {
			if t.StringValue == "layout" {
				l, err := p.parseLayout(t)
				if err != nil {
					return
				}
				f.Statements = append(f.Statements, &Statement{Layout: l})
			} else {
				p.log.LogError(errlog.ErrorUnknownDirective, t.Location, t.StringValue)
			}
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

	// Parse optional name
	if t, ok := p.optional(TokenIdentifier); ok {
		rail.Name = t.StringValue
	}

	// Parse configutation values, e.g. "railway { ... }"
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
	// Parse rows, e.g. "railway ( ... )"
	if _, err := p.expect(TokenOpenParanthesis); err != nil {
		return nil, err
	}
	if _, err := p.expect(TokenNewline); err != nil {
		return nil, err
	}
	for {
		if _, ok := p.optional(TokenCloseParanthesis); ok {
			break
		}
		row, err := p.parseRow()
		if err != nil {
			return nil, err
		}
		if row != nil {
			rail.Rows = append(rail.Rows, row)
		}
	}
	return rail, nil
}

func (p *Parser) parseRow() (*ExpressionRow, *errlog.Error) {
	row := &ExpressionRow{}
	for {
		if _, ok := p.optional(TokenNewline); ok {
			break
		}
		exp, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		row.Expressions = append(row.Expressions, exp)
	}
	if len(row.Expressions) == 0 {
		return nil, nil
	}
	return row, nil
}

/*
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
*/

func (p *Parser) parseExpression() (*Expression, *errlog.Error) {
	if p.peek(TokenAt) || p.peek(TokenEnd) || p.peek(TokenEllipsis) || p.peek(TokenPipe) {
		return p.parseSimpleExpression()
	}
	//
	// RepeatExpression
	//
	if t, ok := p.optional(TokenInteger); ok {
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
		exp, err := p.parseSimpleExpression()
		if err != nil {
			return nil, err
		}
		rep := &Expression{Repeat: &RepeatExpression{Count: int(count), TrackExpression: exp}, Location: t.Location}
		return rep, nil
	}
	//
	// TrackExpression, TrackTerminationExpression or SwitchExpression
	//
	t, switchok := p.optionalMulti(TokenSlash, TokenBackslash)
	if switchok {
		// Consume dashes
		for {
			if _, ok := p.optional(TokenDash); !ok {
				break
			}
		}
	}
	exp, err := p.parseSimpleExpression()
	if err != nil {
		return nil, err
	}
	t2, switchok2 := p.optionalMulti(TokenSlash, TokenBackslash, TokenDash)
	if switchok2 && t2.Kind == TokenDash {
		// Consume dashes
		for {
			t2, err = p.expectMulti(TokenSlash, TokenBackslash, TokenDash)
			if err != nil {
				return nil, err
			}
			if t2.Kind != TokenDash {
				break
			}
		}
	}
	if switchok || switchok2 {
		// It is a switch expression
		if exp.Track == nil {
			p.log.LogError(errlog.ErrorNoSwitchInSwitchExpression, exp.Location)
		}
		s := &SwitchExpression{TrackExpression: *exp.Track}
		if switchok {
			s.PositionLeft = t.Location.Position()
			if t.Kind == TokenSlash {
				s.SplitLeft = true
			} else {
				s.JoinLeft = true
			}
		}
		if switchok2 {
			s.PositionRight = t2.Location.Position()
			if t2.Kind == TokenSlash {
				s.JoinRight = true
			} else {
				s.SplitRight = true
			}
		}
		return &Expression{Switch: s, Location: exp.Location}, nil
	}
	return exp, nil
}

func (p *Parser) parseSimpleExpression() (*Expression, *errlog.Error) {
	t, err := p.expectMulti(TokenIdentifier, TokenEnd, TokenEllipsis, TokenAt, TokenString, TokenPipe, TokenInteger /*, TokenOpenParanthesis*/)
	if err != nil {
		return nil, err
	}
	/*
		if t.Kind == TokenOpenParanthesis {
			return p.parseExpressionList(t)
		}
	*/
	if t.Kind == TokenPipe {
		exp := &Expression{Placeholder: true, Location: t.Location}
		return exp, nil
	}
	if t.Kind == TokenEllipsis {
		ident, err := p.expect(TokenIdentifier)
		if err != nil {
			return nil, err
		}
		exp := &Expression{TrackTermination: &TrackTerminationExpression{Name: ident.StringValue, EllipsisLeft: true}, Location: t.Location}
		return exp, nil
	}
	if t.Kind == TokenIdentifier {
		exp := &Expression{Track: &TrackExpression{Type: t.StringValue}, Location: t.Location}
		if _, ok := p.optional(TokenOpenParanthesis); ok {
			for i := 0; ; i++ {
				if _, ok := p.optional(TokenCloseParanthesis); ok {
					break
				}
				if i != 0 {
					if _, err := p.expect(TokenComma); err != nil {
						return nil, err
					}
				}
				p, err := p.parseParameter()
				if err != nil {
					return nil, err
				}
				exp.Track.Parameters = append(exp.Track.Parameters, p)
			}
		} else if _, ok := p.optional(TokenEllipsis); ok {
			exp.TrackTermination = &TrackTerminationExpression{Name: exp.Track.Type, EllipsisRight: true}
			exp.Track = nil
		}
		return exp, nil
	} else if t.Kind == TokenEnd {
		exp := &Expression{TrackTermination: &TrackTerminationExpression{}, Location: t.Location}
		return exp, nil
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
		exp := &Expression{Anchor: &AnchorExpression{X: x, Y: y, Z: z, Angle: a}, Location: t.Location.Join(t2.Location)}
		return exp, nil
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
		exp, err := p.parseSimpleExpression()
		if err != nil {
			return nil, err
		}
		if exp.Track == nil {
			p.log.LogError(errlog.ErrorNoTrackInRepeatExpression, exp.Location)
		}
		rep := &Expression{Repeat: &RepeatExpression{Count: int(count), TrackExpression: exp}}
		return rep, nil
	}
	panic("Oooops")
}

func (p *Parser) parseParameter() (interface{}, *errlog.Error) {
	// TODO
	return p.expectMulti(TokenIdentifier, TokenInteger)
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

func (p *Parser) parseLayout(t *Token) (l *Layout, err *errlog.Error) {
	_, err = p.expect(TokenOpenParanthesis)
	if err != nil {
		return
	}
	str, lstr := p.l.ScanRawText(')')
	_, err = p.expect(TokenCloseParanthesis)
	if err != nil {
		return
	}
	return &Layout{RawText: str, LocationToken: t.Location, LocationText: lstr}, nil
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

/*
func (p *Parser) expectError(tokenKind ...TokenKind) *errlog.Error {
	t := p.scan()
	var str = []string{t.Raw}
	for _, kind := range tokenKind {
		str = append(str, p.l.TokenKindToString(kind))
	}
	err := p.log.LogError(errlog.ErrorExpectedToken, t.Location, str...)
	return err
}
*/

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
