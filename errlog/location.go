package errlog

type Location int64

type LocationRange struct {
	From Location
	To   Location
}

// SourceFile represents a file in a LocationMap.
type SourceFile struct {
	Name string
}

func (loc Location) Position() int {
	pos := int(uint64(loc) & 0xffffffff)
	return pos
}

func (loc Location) Line() int {
	line := int((uint64(loc) & 0xffff00000000) >> 32)
	return line
}

func (r LocationRange) File() int {
	return int(uint64(r.From) >> 48)
}

func (r LocationRange) Position() int {
	pos := int(uint64(r.From) & 0xffffffff)
	return pos
}

func (r LocationRange) Line() int {
	line := int((uint64(r.From) & 0xffff00000000) >> 32)
	return line
}

func (loc Location) File() int {
	return int(uint64(loc) >> 48)
}

func NewSourceFile(name string) *SourceFile {
	return &SourceFile{Name: name}
}

func EncodeLocation(file int, line int, pos int) Location {
	return Location((uint64(file) << 48) | (uint64(line<<32) | (uint64(pos))))
}

func EncodeLocationRange(file int, fromLine int, fromPos int, toLine int, toPos int) LocationRange {
	from := Location((uint64(file) << 48) | (uint64(fromLine<<32) | (uint64(fromPos))))
	to := Location((uint64(file) << 48) | (uint64(toLine<<32) | (uint64(toPos))))
	return LocationRange{From: from, To: to}
}

// Join ...
func (l LocationRange) Join(l2 LocationRange) LocationRange {
	if l.IsNull() {
		return l2
	}
	if l2.IsNull() {
		return l
	}
	return LocationRange{From: l.From, To: l2.To}
}

// IsNull ...
func (l LocationRange) IsNull() bool {
	return l.From == 0
}

// IsEqualLine ...
func IsEqualLine(l1, l2 LocationRange) bool {
	f1 := uint64(l1.From) >> 48
	line1 := int((uint64(l1.From) & 0xffff00000000) >> 32)
	f2 := uint64(l2.From) >> 48
	line2 := int((uint64(l2.From) & 0xffff00000000) >> 32)

	return f1 == f2 && line1 == line2
}
