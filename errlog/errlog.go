package errlog

import "fmt"

type ErrorLog struct {
	errors   []*Error
	warnings []*Error
	files    []*SourceFile
}

func NewErrorLog() *ErrorLog {
	return &ErrorLog{}
}

func (log *ErrorLog) LogError(code ErrorCode, loc LocationRange, args ...string) *Error {
	if code == 0 {
		panic("Oooops")
	}
	err := NewError(code, loc, args...)
	log.errors = append(log.errors, err)
	return err
}

func (log *ErrorLog) LogWarning(code ErrorCode, loc LocationRange, args ...string) *Error {
	if code == 0 {
		panic("Oooops")
	}
	err := NewError(code, loc, args...)
	log.warnings = append(log.warnings, err)
	return err
}

func (log *ErrorLog) AddError(err *Error) {
	log.errors = append(log.errors, err)
}

func (log *ErrorLog) HasErrors() bool {
	return len(log.errors) > 0
}

func (log *ErrorLog) HasWarnings() bool {
	return len(log.warnings) > 0
}

func (log *ErrorLog) AddFile(f *SourceFile) int {
	log.files = append(log.files, f)
	return len(log.files) - 1
}

func (log *ErrorLog) Decode(loc Location) (*SourceFile, int, int) {
	file := log.files[uint64(loc)>>48]
	line := int((uint64(loc) & 0xffff00000000) >> 32)
	pos := int(uint64(loc) & 0xffffffff)
	return file, line, pos
}

func (log *ErrorLog) ErrorToString(err *Error) string {
	loc := err.Location()
	file, line, pos := log.Decode(loc.From)
	//	_, to := l.Resolve(loc.To)
	return fmt.Sprintf("%v %v:%v: %v", file.Name, line, pos, err.ToString(log))
}

func (log *ErrorLog) ToString() string {
	str := ""
	for _, err := range log.errors {
		str += log.ErrorToString(err) + "\n"
	}
	return str
}

func (log *ErrorLog) Print() {
	if log.HasErrors() || log.HasWarnings() {
		println("-----------------------------------------------------------")
		println("Errors:", len(log.errors), " Warnings:", len(log.warnings))
	}
	for _, err := range log.errors {
		println(log.ErrorToString(err))
	}
	for _, w := range log.warnings {
		println(log.ErrorToString(w))
	}
}
