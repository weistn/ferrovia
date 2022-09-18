package interpreter

import (
	"unicode"

	"github.com/weistn/ferrovia/errlog"
)

// type LayoutCellDirection int
type scanDirection int
type LayoutCellConnection int
type LayoutCellType int

const (
	ConnectTop    LayoutCellConnection = 1
	ConnectRight  LayoutCellConnection = 2
	ConnectBottom LayoutCellConnection = 4
	ConnectLeft   LayoutCellConnection = 8
)
const (
	scanLeft scanDirection = iota
	scanRight
	scanUpwards
	scanDownwards
	/*
		UnknownDirection LayoutCellDirection = iota
		Vertical
		Diag45
		Diag135
		Diag225
		Diag315
	*/
)

const (
	UnprocessedCell LayoutCellType = iota
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
	TrackHorizontal LayoutCellType = 20 + iota
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
	switchHorizontalSlash LayoutCellType = 100 + iota
	switchVerticalSlash
	switchHorizontalBackslash
	switchVerticalBackslash
	TrackDoubleBackslash
	TrackDoubleSlash
)

type Layout struct {
	LineCount   int
	ColumnCount int
	Cells       []LayoutCell
	lines       []string
	Location    errlog.LocationRange
	errlog      *errlog.ErrorLog
}

type LayoutCell struct {
	Type        LayoutCellType
	Connections LayoutCellConnection
	Rune        rune
	X           int
	Y           int
}

func NewLayout(lines []string, loc errlog.LocationRange) *Layout {
	l := &Layout{LineCount: len(lines), lines: lines, Location: loc}
	for _, line := range lines {
		c := 0
		for range line {
			c++
		}
		if c > l.ColumnCount {
			l.ColumnCount = c
		}
	}
	l.Cells = make([]LayoutCell, l.LineCount*l.ColumnCount)
	for y, line := range lines {
		for x, r := range line {
			l.Cells[y*l.ColumnCount+x] = LayoutCell{Rune: r, X: x, Y: y}
		}
	}
	return l
}

func (l *Layout) Cell(x int, y int) *LayoutCell {
	if x < 0 || y < 0 || x >= l.ColumnCount || y >= l.LineCount {
		return nil
	}
	return &l.Cells[y*l.ColumnCount+x]
}

func (l *Layout) CellBelow(x int, y int) *LayoutCell {
	if y+1 >= l.LineCount {
		return nil
	}
	return &l.Cells[(y+1)*l.ColumnCount+x]
}

func (l *Layout) CellAbove(x int, y int) *LayoutCell {
	if y == 0 {
		return nil
	}
	return &l.Cells[(y-1)*l.ColumnCount+x]
}

func (l *Layout) CellRightOf(x int, y int) *LayoutCell {
	if x+1 >= l.ColumnCount {
		return nil
	}
	return &l.Cells[y*l.ColumnCount+x+1]
}

func (l *Layout) CellLeftOf(x int, y int) *LayoutCell {
	if x == 0 {
		return nil
	}
	return &l.Cells[y*l.ColumnCount+x-1]
}

func (l *Layout) Process(log *errlog.ErrorLog) {
	l.errlog = log
	var cellLines [][]*LayoutCell
	for y := 0; y < l.LineCount; y++ {
		var cellLine []*LayoutCell
		for x := 0; x < l.ColumnCount; x++ {
			cellLine = append(cellLine, l.Cell(x, y))
		}
		cellLines = append(cellLines, cellLine)
	}

	var cellColumns [][]*LayoutCell
	for x := 0; x < l.ColumnCount; x++ {
		var cellColumn []*LayoutCell
		for y := 0; y < l.LineCount; y++ {
			cellColumn = append(cellColumn, l.Cell(x, y))
		}
		cellColumns = append(cellColumns, cellColumn)
	}

	for y := 0; y < l.LineCount; y++ {
		for x := 0; x < l.ColumnCount; x++ {
			cell := l.Cell(x, y)
			switch cell.Rune {
			case '-':
				cell.Type = TrackHorizontal
				cell.Connections = ConnectLeft | ConnectRight
				l.processCells(cellLines[y], x, scanLeft)
				l.processCells(cellLines[y], x, scanRight)
			case '|':
				cell.Type = TrackVertical
				cell.Connections = ConnectTop | ConnectBottom
				l.processCells(cellColumns[x], y, scanUpwards)
				l.processCells(cellColumns[x], y, scanDownwards)
			case '/':
				cell.Connections = ConnectLeft | ConnectRight | ConnectTop | ConnectBottom
			case '\\':
				cell.Connections = ConnectLeft | ConnectRight | ConnectTop | ConnectBottom
			case '@', '>', '<', '^':
				// Do nothing by intention
			case ',':
				cell.Type = TrackDiagonalLower
				cell.Connections = ConnectBottom | ConnectRight
			case '.':
				cell.Type = TrackDiagonalBackLower
				cell.Connections = ConnectBottom | ConnectLeft
			case '`':
				cell.Type = TrackDiagonalBackUpper
				cell.Connections = ConnectTop | ConnectRight
			case '\'':
				cell.Type = TrackDiagonalUpper
				cell.Connections = ConnectTop | ConnectLeft
			case '%':
				cell.Type = TrackDoubleSlash
				cell.Connections = ConnectTop | ConnectRight | ConnectBottom | ConnectRight
			case '&':
				cell.Type = TrackDoubleBackslash
				cell.Connections = ConnectTop | ConnectRight | ConnectBottom | ConnectRight
			case 0:
				// Do nothing by intention
			default:
				if unicode.IsSpace(cell.Rune) || unicode.IsLetter(cell.Rune) || unicode.IsDigit(cell.Rune) {
					// Do nothing by intention
				} else {
					l.addError(errlog.ErrorIllegalRune, x, y)
				}
			}
		}
	}

	for y := 0; y < l.LineCount; y++ {
		for x := 0; x < l.ColumnCount; x++ {
			cell := l.Cell(x, y)
			switch cell.Type {
			case switchHorizontalBackslash:
				println("------------ back")
				if cell_above := l.CellAbove(x, y); cell_above != nil && cell_above.ConnectsToBottom() {
					println("above back")
					cell.Connections |= ConnectTop
					cell.Type = SwitchHorizontalDiagonalBackUpper
				}
				if cell_below := l.CellBelow(x, y); cell_below != nil && cell_below.ConnectsToTop() {
					println("below back")
					cell.Connections |= ConnectBottom
					cell.Type = SwitchHorizontalDiagonalLower
				}
				if cell.ConnectsToTop() && cell.ConnectsToBottom() {
					l.addError(errlog.ErrorMalformedLayout, x, y, "Switch has too many connections")
				} else if !cell.ConnectsToTop() && !cell.ConnectsToBottom() {
					l.addError(errlog.ErrorMalformedLayout, x, y, "Switch has too few connections")
				}
			case switchHorizontalSlash:
				println("------------ slash")
				if cell_above := l.CellAbove(x, y); cell_above != nil && cell_above.ConnectsToBottom() {
					println("above")
					cell.Connections = ConnectTop | ConnectLeft | ConnectRight
					cell.Type = SwitchHorizontalDiagonalUpper
				}
				if cell_below := l.CellBelow(x, y); cell_below != nil && cell_below.ConnectsToTop() {
					println("below")
					cell.Connections = ConnectBottom | ConnectLeft | ConnectRight
					cell.Type = SwitchHorizontalDiagonalLower
				}
				if cell.ConnectsToTop() && cell.ConnectsToBottom() {
					l.addError(errlog.ErrorMalformedLayout, x, y, "Switch has too many connections")
				} else if !cell.ConnectsToTop() && !cell.ConnectsToBottom() {
					l.addError(errlog.ErrorMalformedLayout, x, y, "Switch has too few connections")
				}
			case switchVerticalBackslash:
			case switchVerticalSlash:
			}
		}
	}
}

func (l *Layout) processCells(cells []*LayoutCell, pos int, dir scanDirection) {
	inc := 1
	if dir == scanLeft || dir == scanUpwards {
		inc = -1
	}
	// Inspect the following cell
	i := pos + inc
	if i < 0 || i >= len(cells) {
		return
	}
	cell := cells[i]
	switch cell.Rune {
	case '<':
		return
	case '>':
		return
	case '^':
		return
	case '-':
		return
	case '|':
		return
	case '.', ',', '`', '\'', '&', '%':
		return
	case '/':
		if dir == scanRight && i+1 < len(cells) && (cells[i+1].Rune == '/' || cells[i+1].Rune == '-' || cells[i+1].Rune == '\'' || cells[i+1].Rune == '.') {
			cell.Type = switchHorizontalSlash
			cell.Connections |= ConnectLeft | ConnectRight | ConnectTop | ConnectBottom
			l.processCells(cells, i, dir)
		} else if dir == scanLeft && i > 0 && (cells[i-1].Rune == '/' || cells[i-1].Rune == '-' || cells[i-1].Rune == ',' || cells[i+1].Rune == '`') {
			cell.Type = switchHorizontalSlash
			cell.Connections |= ConnectLeft | ConnectRight | ConnectTop | ConnectBottom
			l.processCells(cells, i, dir)
		} else if dir == scanDownwards && i+1 < len(cells) && (cells[i+1].Rune == '/' || cells[i+1].Rune == '|' || cells[i+1].Rune == '\'' || cells[i+1].Rune == '`') {
			cell.Type = switchVerticalSlash
			cell.Connections |= ConnectTop | ConnectBottom | ConnectLeft | ConnectRight
			l.processCells(cells, i, dir)
		} else if dir == scanUpwards && i > 0 && (cells[i-1].Rune == '/' || cells[i-1].Rune == '|' || cells[i-1].Rune == '.' || cells[i+1].Rune == ',') {
			cell.Type = switchVerticalSlash
			cell.Connections |= ConnectTop | ConnectBottom | ConnectLeft | ConnectRight
			l.processCells(cells, i, dir)
		}
		return
	case '\\':
		if dir == scanRight || dir == scanLeft {
			println("!!!!!!!!!!!!!!!! \\")
			cell.Type = switchHorizontalBackslash
		} else {
			cell.Type = switchVerticalBackslash
		}
		l.processCells(cells, i, dir)
		return
	case '@':
		if cell.Type != UnprocessedCell {
			l.addError(errlog.ErrorMalformedLayout, cell.X, cell.Y, "'@' has more than one track connection")
		}
		if dir == scanLeft {
			cell.Type = TrackHorizontalStopLeft
			cell.Connections = ConnectRight
		} else if dir == scanRight {
			cell.Type = TrackHorizontalStopRight
			cell.Connections = ConnectLeft
		} else if dir == scanUpwards {
			cell.Type = TrackVerticalStopTop
			cell.Connections = ConnectBottom
		} else {
			cell.Type = TrackVerticalStopBottom
			cell.Connections = ConnectTop
		}
		return
	default:
		if cell.Type != UnprocessedCell {
			l.addError(errlog.ErrorMalformedLayout, cell.X, cell.Y, " has more than one track connection")
		}
		if cell.Rune == 'B' {
			if (dir == scanDownwards || dir == scanRight) && i+4 < len(cells) && cells[i+1].Rune == 'B' && cells[i+2].Rune == 'B' && cells[i+3].Rune == 'B' && !unicode.IsLetter(cells[i+4].Rune) && !unicode.IsDigit(cells[i+4].Rune) {
				if dir == scanDownwards {
					makeBlock(cells[i:i+4], TrackVerticalBlock)
				} else {
					makeBlock(cells[i:i+4], TrackHorizontalBlock)
				}
			} else if (dir == scanUpwards || dir == scanLeft) && i >= 4 && cells[i-1].Rune == 'B' && cells[i-2].Rune == 'B' && cells[i-3].Rune == 'B' && !unicode.IsLetter(cells[i-4].Rune) && !unicode.IsDigit(cells[i-4].Rune) {
				if dir == scanUpwards {
					makeBlock(cells[i-3:i+1], TrackVerticalBlock)
				} else {
					makeBlock(cells[i-3:i+1], TrackHorizontalBlock)
				}
			}
		}

		if unicode.IsSpace(cell.Rune) || unicode.IsLetter(cell.Rune) || unicode.IsDigit(cell.Rune) {
			start := i
			space := unicode.IsSpace(cell.Rune)
			if space {
				start += inc
			}
			var j int
			for j = start; j >= 0 && j < len(cells); j++ {
				if unicode.IsLetter(cells[j].Rune) || unicode.IsDigit(cells[j].Rune) {
					space = false
					continue
				}
				if unicode.IsSpace(cells[j].Rune) && !space {
					space = true
					continue
				}
				break
			}
			if space {
				j -= inc
			}
			if start == j {
				// Do nothing
			} else if dir == scanDownwards {
				makeLabel(cells[start:j+1], TrackVerticalLabel)
			} else if dir == scanRight {
				makeLabel(cells[start:j+1], TrackHorizontalLabel)
			} else if dir == scanUpwards {
				makeLabel(cells[start:j+1], TrackVerticalLabel)
			} else {
				makeLabel(cells[j:start+1], TrackHorizontalLabel)
			}
		}
	}
}

func makeLabel(cells []*LayoutCell, kind LayoutCellType) {

}

func makeBlock(cells []*LayoutCell, kind LayoutCellType) {
	for _, c := range cells {
		if kind == TrackHorizontalBlock {
			c.Connections = ConnectLeft | ConnectRight
		} else {
			c.Connections = ConnectTop | ConnectBottom
		}
		c.Type = kind
	}
}

func (l *Layout) addError(code errlog.ErrorCode, x int, y int, args ...string) *errlog.Error {
	loc := errlog.EncodeLocationRange(l.Location.File(), l.Location.Line()+y, l.Location.Position()+x, l.Location.Line()+y, l.Location.Position()+x)
	return l.errlog.LogError(code, loc, args...)
}

func (c *LayoutCell) ConnectsToTop() bool {
	return c.Connections&ConnectTop == ConnectTop
}

func (c *LayoutCell) ConnectsToBottom() bool {
	return c.Connections&ConnectBottom == ConnectBottom
}
