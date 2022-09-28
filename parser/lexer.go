package parser

import (
	"bytes"

	"github.com/weistn/ferrovia/errlog"
)

// Lexer ...
type Lexer struct {
	t    *tokenizer
	file int
	str  string
	pos  int
	log  *errlog.ErrorLog
}

// NewLexer ...
func NewLexer(file int, str string, log *errlog.ErrorLog) *Lexer {
	t := newTokenizer()
	t.addTokenDefinition("/*", TokenBlockComment)
	t.addTokenDefinition("//", TokenLineComment)
	t.addTokenDefinition("(", TokenOpenParanthesis)
	t.addTokenDefinition(")", TokenCloseParanthesis)
	t.addTokenDefinition("{", TokenOpenBraces)
	t.addTokenDefinition("}", TokenCloseBraces)
	t.addTokenDefinition("*", TokenAsterisk)
	t.addTokenDefinition(",", TokenComma)
	t.addTokenDefinition("@", TokenAt)
	t.addTokenDefinition("mm", TokenUnit)
	t.addTokenDefinition("cm", TokenUnit)
	t.addTokenDefinition("m", TokenUnit)
	t.addTokenDefinition("deg", TokenUnit)
	t.addTokenDefinition("\r\n", TokenNewline)
	t.addTokenDefinition("\n", TokenNewline)

	// Not used
	t.addTokenDefinition(":", TokenColon)
	t.addTokenDefinition("-", TokenDash)
	t.addTokenDefinition("|", TokenPipe)
	t.addTokenDefinition("...", TokenEllipsis)
	t.addTokenDefinition("/", TokenSlash)
	t.addTokenDefinition("&", TokenAmpersand)
	t.addTokenDefinition("[", TokenOpenBracket)
	t.addTokenDefinition("]", TokenCloseBracket)
	t.addTokenDefinition(".", TokenDot)
	t.addTokenDefinition(";", TokenSemicolon)
	t.addTokenDefinition("~", TokenTilde)
	t.addTokenDefinition("#", TokenHash)
	t.addTokenDefinition("!", TokenBang)
	t.addTokenDefinition("^", TokenCaret)
	t.addTokenDefinition("+", TokenPlus)
	t.addTokenDefinition("%", TokenPercent)
	t.addTokenDefinition("||", TokenLogicalOr)
	t.addTokenDefinition("&&", TokenLogicalAnd)
	t.addTokenDefinition("&^", TokenBitClear)
	t.addTokenDefinition("=", TokenAssign)
	t.addTokenDefinition("+=", TokenAssignPlus)
	t.addTokenDefinition("-=", TokenAssignMinus)
	t.addTokenDefinition("*=", TokenAssignAsterisk)
	t.addTokenDefinition("/=", TokenAssignDivision)
	t.addTokenDefinition("<<=", TokenAssignShiftLeft)
	t.addTokenDefinition(">>=", TokenAssignShiftRight)
	t.addTokenDefinition("|=", TokenAssignBinaryOr)
	t.addTokenDefinition("&=", TokenAssignBinaryAnd)
	t.addTokenDefinition("%=", TokenAssignPercent)
	t.addTokenDefinition("^=", TokenAssignCaret)
	t.addTokenDefinition("&^=", TokenAssignAndCaret)
	t.addTokenDefinition(":=", TokenWalrus)
	t.addTokenDefinition("<<", TokenShiftLeft)
	t.addTokenDefinition(">>", TokenShiftRight)
	t.addTokenDefinition("==", TokenEqual)
	t.addTokenDefinition("!=", TokenNotEqual)
	t.addTokenDefinition("<=", TokenLessOrEqual)
	t.addTokenDefinition(">=", TokenGreaterOrEqual)
	t.addTokenDefinition("<", TokenLess)
	t.addTokenDefinition(">", TokenGreater)
	t.addTokenDefinition("true", TokenTrue)
	t.addTokenDefinition("false", TokenFalse)
	t.addTokenDefinition("null", TokenNull)
	t.polish()
	l := &Lexer{t: t, str: str, file: file, log: log}
	return l
}

// TokenKindToString ...
func (l *Lexer) TokenKindToString(kind TokenKind) string {
	switch kind {
	case TokenIdentifier:
		return "identifier"
	case TokenString:
		return "string literal"
	case TokenRune:
		return "rune literal"
	case TokenInteger:
		return "integer number"
	case TokenFloat:
		return "floating point number"
	case TokenHex:
		return "hex number"
	case TokenOctal:
		return "octal number"
	}
	return l.t.tokenKindToString(kind)
}

// Scan ...
func (l *Lexer) Scan() *Token {
	start := l.pos
	var token *Token
	for {
		l.skipWhitespace()
		if l.pos == len(l.str) {
			token = &Token{Kind: TokenEOF, Location: encodeRange(l.file, l.str, l.pos, l.pos)}
			break
		}
		token, l.pos = l.t.scan(l.file, l.str, l.pos)
		if token.Kind == TokenError {
			l.log.LogError(token.ErrorCode, token.Location)
			continue
		}
		if token.Kind == TokenLineComment {
			l.skipLineComment()
		} else if token.Kind == TokenBlockComment {
			l.skipBlockComment()
		} else {
			break
		}
	}
	token.Raw = l.str[start:l.pos]
	return token
}

func (l *Lexer) skipLineComment() {
	for ; l.pos < len(l.str); l.pos++ {
		ch := l.str[l.pos]
		if ch == '\n' || (ch == '\r' && l.pos+1 < len(l.str) && l.str[l.pos+1] == '\n') {
			break
		}
	}
}

// This function is used to scan ASCII-Art text, until it encounters a new line that starts with the character 'term'.
// The function expands tabulators, assuming a tab-width of 4 characters.
func (l *Lexer) ScanRawText(term byte) (string, errlog.LocationRange) {
	l.skipWhitespace()
	if l.pos == len(l.str) {
		return "", encodeRange(l.file, l.str, l.pos, l.pos)

	}
	if l.str[l.pos] == '\n' {
		l.pos++
	}
	newline := true
	start := l.pos
	end := start
	linepos := 0
	tab := "    "
	var buf bytes.Buffer
	// var buf bytes.Buffer
	for ; l.pos < len(l.str); l.pos++ {
		ch := l.str[l.pos]
		if ch == '\n' {
			newline = true
			linepos = 0
			end = l.pos
			buf.WriteByte(ch)
		} else if ch == '\r' {
			continue
		} else if newline && ch == term {
			break
		} else if ch == ' ' {
			linepos++
			buf.WriteByte(ch)
		} else if ch == '\t' {
			off := linepos % 4
			linepos += 4 - off
			buf.WriteString(tab[:4-off])
		} else {
			linepos++
			buf.WriteByte(ch)
			newline = false
		}
	}
	return buf.String(), encodeRange(l.file, l.str, start, end)
}

func (l *Lexer) skipBlockComment() {
	for ; l.pos < len(l.str); l.pos++ {
		ch := l.str[l.pos]
		if ch == '*' && l.pos+1 < len(l.str) && l.str[l.pos+1] == '/' {
			l.pos += 2
			break
		}
	}
}

func (l *Lexer) skipWhitespace() {
	for ; l.pos < len(l.str); l.pos++ {
		ch := l.str[l.pos]
		if ch != ' ' && ch != '\t' {
			break
		}
	}
}
