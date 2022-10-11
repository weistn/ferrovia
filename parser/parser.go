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
		t, err := p.expectMulti(TokenEOF, TokenNewline, TokenIdentifier)
		if err != nil {
			p.log.AddError(err)
			return
		}
		if t.Kind == TokenNewline {
			continue
		}
		if t.Kind == TokenIdentifier {
			if t.StringValue == "tracks" {
				tracks, err := p.parseTracks(t)
				if err != nil {
					p.log.AddError(err)
					return
				}
				f.Statements = append(f.Statements, tracks)
			} else if t.StringValue == "layer" {
				l, err := p.parseLayer(t)
				if err != nil {
					p.log.AddError(err)
					return
				}
				f.Statements = append(f.Statements, l)
			} else if t.StringValue == "ground" {
				ground, err := p.parseGround(t)
				if err != nil {
					p.log.AddError(err)
					return
				}
				f.Statements = append(f.Statements, ground)
			} else if t.StringValue == "switchboard" {
				sb, err := p.parseSwitchboard(t)
				if err != nil {
					return
				}
				f.Statements = append(f.Statements, sb)
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

func (p *Parser) parseBody() ([]IExpression, *errlog.Error) {
	var expressions []IExpression
	// Parse body
	for {
		if _, ok := p.optional(TokenNewline); ok {
			continue
		}
		if _, ok := p.optional(TokenCloseBraces); ok {
			break
		}
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		//
		// ContextExpression?
		//
		if _, ok := p.optional(TokenOpenBraces); ok {
			statements, err := p.parseBody()
			if err != nil {
				return nil, err
			}
			expr = &ContextExpression{Object: expr, Statements: statements}
		}

		expressions = append(expressions, expr)
	}
	return expressions, nil
}

func (p *Parser) parseLayer(t *Token) (*Layer, *errlog.Error) {
	l := &Layer{Location: t.Location}
	var err *errlog.Error

	// Parse name
	l.Name, err = p.expect(TokenIdentifier)
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(TokenOpenBraces); err != nil {
		return nil, err
	}
	if _, err := p.expect(TokenNewline); err != nil {
		return nil, err
	}

	// Parse body
	l.Expressions, err = p.parseBody()
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (p *Parser) parseTracks(t *Token) (*Tracks, *errlog.Error) {
	tracks := &Tracks{Location: t.Location}

	// Parse optional name and optional parameter list
	if t, ok := p.optional(TokenIdentifier); ok {
		tracks.Name = t
		if _, ok := p.optional(TokenOpenParanthesis); ok {
			var params []*Parameter
			for {
				if _, ok := p.optional(TokenCloseParanthesis); ok {
					break
				}
				if len(params) != 0 {
					if _, err := p.expect(TokenComma); err != nil {
						return nil, err
					}
				}
				t, err := p.expect(TokenIdentifier)
				if err != nil {
					return nil, err
				}
				params = append(params, &Parameter{Name: t})
			}
			tracks.Parameters = params
		}
	}

	if _, err := p.expect(TokenOpenBraces); err != nil {
		return nil, err
	}
	if _, err := p.expect(TokenNewline); err != nil {
		return nil, err
	}

	// Parse body
	var err *errlog.Error
	tracks.Expressions, err = p.parseBody()
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

func (p *Parser) parseExpression() (IExpression, *errlog.Error) {
	expr, err := p.parseDotExpression()
	if err != nil {
		return nil, err
	}

	//
	// BinaryExpression ?
	//
	if t, ok := p.optional(TokenAsterisk); ok {
		expr2, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		return &BinaryExpression{Left: expr, Op: t, Right: expr2}, nil
	}
	return expr, nil
}

func (p *Parser) parseDotExpression() (IExpression, *errlog.Error) {
	expr, err := p.parseCallExpression()
	if err != nil {
		return nil, err
	}

	//
	// DotExpression?
	//
	for {
		if _, ok := p.optional(TokenDot); !ok {
			return expr, nil
		}
		ident, err := p.expect(TokenIdentifier)
		if err != nil {
			return nil, err
		}
		dot := &DotExpression{Context: expr, Identifier: ident}
		if _, ok := p.optional(TokenOpenParanthesis); ok {
			var args []IExpression
			for {
				if _, ok := p.optional(TokenCloseParanthesis); ok {
					break
				}
				if len(args) != 0 {
					if _, err := p.expect(TokenComma); err != nil {
						return nil, err
					}
				}
				arg, err := p.parseExpression()
				if err != nil {
					return nil, err
				}
				args = append(args, arg)
			}
			dot.Arguments = args
		}
		expr = dot
	}
}

func (p *Parser) parseCallExpression() (IExpression, *errlog.Error) {
	expr, err := p.parseSimpleExpression()
	if err != nil {
		return nil, err
	}

	//
	// CallExpression ?
	//
	if _, ok := p.optional(TokenOpenParanthesis); ok {
		var args []IExpression
		for {
			if _, ok := p.optional(TokenCloseParanthesis); ok {
				break
			}
			if len(args) != 0 {
				if _, err := p.expect(TokenComma); err != nil {
					return nil, err
				}
			}
			arg, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
		}
		return &CallExpression{Func: expr, Arguments: args}, nil
	}
	return expr, nil
}

func (p *Parser) parseSimpleExpression() (IExpression, *errlog.Error) {
	t, err := p.expectMulti(TokenIdentifier, TokenAt, TokenString, TokenInteger, TokenFloat, TokenOpenParanthesis, TokenOpenBracket)
	if err != nil {
		return nil, err
	}

	// For strings and @, dimensions are not allowed
	if t.Kind == TokenString {
		return &ConstantExpression{Value: t}, nil
	} else if t.Kind == TokenAt {
		return &IdentifierExpression{Identifier: t}, nil
	} else if t.Kind == TokenOpenBracket {
		var values []IExpression
		for {
			if _, ok := p.optional(TokenCloseBracket); ok {
				break
			}
			if len(values) != 0 {
				if _, err := p.expect(TokenComma); err != nil {
					return nil, err
				}
			}
			value, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			values = append(values, value)
		}
		return &VectorExpression{Values: values, Location: t.Location}, nil
	}

	var expr IExpression
	if t.Kind == TokenIdentifier || t.Kind == TokenAt {
		expr = &IdentifierExpression{Identifier: t}
	} else if t.Kind == TokenInteger {
		if !t.IntegerValue.IsInt64() {
			return nil, errlog.NewError(errlog.ErrorIllegalNumber, t.Location)
		}
		expr = &ConstantExpression{Value: t}
	} else if t.Kind == TokenFloat {
		expr = &ConstantExpression{Value: t}
	} else if t.Kind == TokenOpenParanthesis {
		expr, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(TokenCloseParanthesis); err != nil {
			return nil, err
		}
	} else {
		panic("Oooops")
	}

	// Check for dimension
	if t, ok := p.optional(TokenUnit); ok {
		return &DimensionExpression{Value: expr, Dimension: t}, nil
	}
	return expr, nil
}

func (p *Parser) parseGround(t *Token) (*GroundPlate, *errlog.Error) {
	if _, err := p.expect(TokenOpenBraces); err != nil {
		return nil, err
	}
	if _, err := p.expect(TokenNewline); err != nil {
		return nil, err
	}
	ground := &GroundPlate{Location: t.Location}

	// Parse body
	var err *errlog.Error
	ground.Expressions, err = p.parseBody()
	if err != nil {
		return nil, err
	}

	return ground, nil
}

func (p *Parser) parseSwitchboard(t *Token) (sb *Switchboard, err *errlog.Error) {
	_, err = p.expect(TokenOpenBraces)
	if err != nil {
		return
	}
	str, lstr := p.l.ScanRawText('}')
	_, err = p.expect(TokenCloseBraces)
	if err != nil {
		return
	}
	return &Switchboard{RawText: str, LocationToken: t.Location, LocationText: lstr}, nil
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

/*
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
*/

/*
func (p *Parser) peek(tokenKind TokenKind) bool {
	t := p.scan()
	p.savedToken = t
	return t.Kind == tokenKind
}
*/

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
