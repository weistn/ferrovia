package errlog

// ErrorCode ...
type ErrorCode int

const (
	ErrorIllegalNumber ErrorCode = 1 + iota
	ErrorIllegalRune
	ErrorIllegalString
	ErrorIllegalCharacter
	ErrorIllegalUnit
	ErrorIllegalProperty
	ErrorTracksWithoutPosition
	ErrorUnknownTrackType
	ErrorTrackConnectedTwice
	ErrorTrackMarkDefinedTwice
	ErrorTrackPositionedTwice
	ErrorMismatchJunctionCount
	ErrorIllegalJunction
	ErrorDuplicateLayer
	ErrorUnknownLayer
	ErrorUnexpectedEOF
	ErrorExpectedToken
)

type Error struct {
	code      ErrorCode
	location  LocationRange
	args      []string
	locations []LocationRange
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
	case ErrorIllegalProperty:
		return "Illegal property " + e.args[0]
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
	case ErrorMismatchJunctionCount:
		return "The number of connected tracks does not match the number of provided connections"
	case ErrorIllegalJunction:
		return "The junction is not allowed in this place"
	case ErrorUnexpectedEOF:
		return "Unexpected end of file"
	}
	println(e.code)
	panic("Should not happen")
}

// Location ...
func (e *Error) Location() LocationRange {
	return e.location
}
