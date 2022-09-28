package errlog

import "fmt"

// ErrorCode ...
type ErrorCode int

const (
	// Parser Errors
	ErrorIllegalNumber ErrorCode = 1 + iota
	ErrorIllegalRune
	ErrorIllegalString
	ErrorIllegalCharacter
	ErrorIllegalUnit
	ErrorUnexpectedEOF
	ErrorExpectedToken
	ErrorUnknownDirective
	ErrorMalformedLayout

	// Runtime errors
	ErrorTracksWithoutPosition
	ErrorUnknownTrackType
	ErrorTrackConnectedTwice
	ErrorTrackMarkDefinedTwice
	ErrorTrackPositionedTwice
	ErrorDuplicateLayer
	ErrorUnknownLayer
	ErrorNoTrackInRepeatExpression
	ErrorNamedRailwayUsedTwice
	ErrorTypeMismtach
	ErrorArgumentCountMismatch
)

type Error struct {
	code     ErrorCode
	location LocationRange
	args     []string
	//	locations []LocationRange
}

func NewError(code ErrorCode, loc LocationRange, args ...string) *Error {
	return &Error{code: code, location: loc, args: args}
}

// Error ...
func (e *Error) Error() string {
	return e.ToString(nil)
}

// ToString ...
func (e *Error) ToString(log *ErrorLog) string {
	switch e.code {
	case ErrorIllegalNumber:
		return "Illegal number"
	case ErrorIllegalRune:
		return "Illegal rune"
	case ErrorIllegalString:
		return "Illegal string"
	case ErrorIllegalCharacter:
		return "Illegal character"
	case ErrorIllegalUnit:
		return "Illegal unit " + e.args[0]
	case ErrorTracksWithoutPosition:
		return "Tracks have no defined position"
	case ErrorUnknownTrackType:
		return "Unknown track type " + e.args[0]
	case ErrorTrackConnectedTwice:
		return "The track has been connected twice"
	case ErrorTrackMarkDefinedTwice:
		return "The mark " + e.args[0] + " has been defined twice on the same track"
	case ErrorTrackPositionedTwice:
		return "More than one position has been defined for the track"
	case ErrorUnknownDirective:
		return "Unknown directive " + e.args[0]
	case ErrorExpectedToken:
		str := "`" + e.args[1] + "`"
		if e.args[1] == "\n" || e.args[1] == "\r\n" {
			str = "`end of line`"
		}
		for i := 2; i < len(e.args); i++ {
			if e.args[i] == "\n" || e.args[i] == "\r\n" {
				str += " or " + "`end of line`"
			} else {
				str += " or " + "`" + e.args[i] + "`"
			}
		}
		if e.args[0] == "\n" || e.args[0] == "\r\n" {
			return "Expected " + str + " but got " + "end of line"
		}
		return "Expected " + str + " but got " + "`" + e.args[0] + "`"
	case ErrorDuplicateLayer:
		return "Two layers of the same name `" + e.args[0] + "` defined"
	case ErrorUnknownLayer:
		return "Unknown layer `" + e.args[0] + "`"
	case ErrorNoTrackInRepeatExpression:
		return "A repeat expression requires a track type to repeat"
	case ErrorNamedRailwayUsedTwice:
		return fmt.Sprintf("The railway %v has been used twice", e.args[0])
	case ErrorUnexpectedEOF:
		return "Unexpected end of file"
	case ErrorMalformedLayout:
		return "Malformed layout: " + e.args[0]
	case ErrorArgumentCountMismatch:
		return fmt.Sprintf("Argument count mismatch. Expected %v parameters", e.args[0])
	}
	println(e.code)
	panic("Should not happen")
}

// Location ...
func (e *Error) Location() LocationRange {
	return e.location
}
