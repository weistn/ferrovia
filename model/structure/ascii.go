package structure

import (
	"github.com/weistn/ferrovia/errlog"
)

type ASCIICellConnection int
type ASCIICellType int

const (
	ConnectTop    ASCIICellConnection = 1
	ConnectRight  ASCIICellConnection = 2
	ConnectBottom ASCIICellConnection = 4
	ConnectLeft   ASCIICellConnection = 8
)

const (
	UnprocessedCell ASCIICellType = iota
	SwitchVerticalDiagonalUpper
	SwitchVerticalDiagonalBackUpper
	SwitchVerticalDiagonalBackLower
	SwitchVerticalDiagonalLower
	SwitchHorizontalDiagonalUpper
	SwitchHorizontalDiagonalBackUpper
	SwitchHorizontalDiagonalBackLower
	SwitchHorizontalDiagonalLower
	SwitchVerticalCrossDiagonal
	SwitchVerticalCrossDiagonalBack
	SwitchHorizontalCrossDiagonal
	SwitchHorizontalCrossDiagonalBack
)

const (
	TrackHorizontal ASCIICellType = 20 + iota
	TrackVertical
	TrackDiagonalUpper
	TrackDiagonalLower
	TrackDiagonalBackUpper
	TrackDiagonalBackLower
	TrackHorizontalStopRight
	TrackHorizontalStopLeft
	TrackVerticalStopTop
	TrackVerticalStopBottom
	TrackHorizontalBlock
	TrackVerticalBlock
	TrackHorizontalLabel
	TrackVerticalLabel
)

const (
	// Used internally only by the interpreter
	SwitchHorizontalSlash ASCIICellType = 256
	// Used internally only by the interpreter
	SwitchVerticalSlash ASCIICellType = 512
	// Used internally only by the interpreter
	SwitchHorizontalBackslash ASCIICellType = 1024
	// Used internally only by the interpreter
	SwitchVerticalBackslash ASCIICellType = 2048
)

const (
	// Used internally only, not understood by the HTML UI
	TrackDoubleBackslash ASCIICellType = iota + 3000
	// Used internally only, not understood by the HTML UI
	TrackDoubleSlash
)
/*
 * ASCIIStructure represents an ASCII representation of a layout structure.
 * Instances of this type are created while interpreting *.via files.
 */
type ASCIIStructure struct {
	LineCount   int
	ColumnCount int
	Lines       []string
	Cells       []ASCIICell
	Location    errlog.LocationRange
}

type ASCIICell struct {
	Type        ASCIICellType
	Connections ASCIICellConnection
	Rune        rune
	X           int
	Y           int
	Anchor      *ASCIICell
}

func (l *ASCIIStructure) Cell(x int, y int) *ASCIICell {
	if x < 0 || y < 0 || x >= l.ColumnCount || y >= l.LineCount {
		return nil
	}
	return &l.Cells[y*l.ColumnCount+x]
}

func (l *ASCIIStructure) CellBelow(x int, y int) *ASCIICell {
	if y+1 >= l.LineCount {
		return nil
	}
	return &l.Cells[(y+1)*l.ColumnCount+x]
}

func (l *ASCIIStructure) CellAbove(x int, y int) *ASCIICell {
	if y == 0 {
		return nil
	}
	return &l.Cells[(y-1)*l.ColumnCount+x]
}

func (l *ASCIIStructure) CellRightOf(x int, y int) *ASCIICell {
	if x+1 >= l.ColumnCount {
		return nil
	}
	return &l.Cells[y*l.ColumnCount+x+1]
}

func (l *ASCIIStructure) CellLeftOf(x int, y int) *ASCIICell {
	if x == 0 {
		return nil
	}
	return &l.Cells[y*l.ColumnCount+x-1]
}

func (c *ASCIICell) ConnectsToTop() bool {
	return c.Connections&ConnectTop == ConnectTop
}

func (c *ASCIICell) ConnectsToBottom() bool {
	return c.Connections&ConnectBottom == ConnectBottom
}

func (c *ASCIICell) ConnectsToLeft() bool {
	return c.Connections&ConnectLeft == ConnectLeft
}

func (c *ASCIICell) ConnectsToRight() bool {
	return c.Connections&ConnectRight == ConnectRight
}

func MakeLabel(cells []*ASCIICell, kind ASCIICellType) {

}

func MakeBlock(cells []*ASCIICell, kind ASCIICellType) {
	for i, c := range cells {
		if kind == TrackHorizontalBlock {
			c.Connections = ConnectLeft | ConnectRight
		} else {
			c.Connections = ConnectTop | ConnectBottom
		}
		c.Type = kind
		if i > 0 {
			c.Anchor = cells[0]
		}
	}
}
