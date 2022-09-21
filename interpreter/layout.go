package interpreter

import (
	"unicode"

	"github.com/weistn/ferrovia/errlog"
	. "github.com/weistn/ferrovia/model/structure"
)

type scanDirection int

const (
	scanLeft scanDirection = iota
	scanRight
	scanUpwards
	scanDownwards
)

func processASCIIStructure(lines []string, loc errlog.LocationRange, log *errlog.ErrorLog) *ASCIIStructure {
	l := &ASCIIStructure{LineCount: len(lines), Lines: lines, Location: loc}
	for _, line := range lines {
		c := 0
		for range line {
			c++
		}
		if c > l.ColumnCount {
			l.ColumnCount = c
		}
	}
	l.Cells = make([]ASCIICell, l.LineCount*l.ColumnCount)
	for y, line := range lines {
		for x, r := range line {
			l.Cells[y*l.ColumnCount+x] = ASCIICell{Rune: r, X: x, Y: y}
		}
	}

	// Create lines of cells
	var cellLines [][]*ASCIICell
	for y := 0; y < l.LineCount; y++ {
		var cellLine []*ASCIICell
		for x := 0; x < l.ColumnCount; x++ {
			cellLine = append(cellLine, l.Cell(x, y))
		}
		cellLines = append(cellLines, cellLine)
	}

	// Create columns of cells
	var cellColumns [][]*ASCIICell
	for x := 0; x < l.ColumnCount; x++ {
		var cellColumn []*ASCIICell
		for y := 0; y < l.LineCount; y++ {
			cellColumn = append(cellColumn, l.Cell(x, y))
		}
		cellColumns = append(cellColumns, cellColumn)
	}

	// Determine the meaning of symbols such as / \ @, which have to be interpreted by their surrounding cells.
	// Use symbols of known meaning like - | as a starting point and process neighbouring cells from there.
	for y := 0; y < l.LineCount; y++ {
		for x := 0; x < l.ColumnCount; x++ {
			cell := l.Cell(x, y)
			switch cell.Rune {
			case '-':
				cell.Type = TrackHorizontal
				cell.Connections = ConnectLeft | ConnectRight
				processCells(l, cellLines[y], x, scanLeft, log)
				processCells(l, cellLines[y], x, scanRight, log)
			case '|':
				cell.Type = TrackVertical
				cell.Connections = ConnectTop | ConnectBottom
				processCells(l, cellColumns[x], y, scanUpwards, log)
				processCells(l, cellColumns[x], y, scanDownwards, log)
			case '/':
				cell.Connections = ConnectLeft | ConnectRight | ConnectTop | ConnectBottom
			case '\\':
				cell.Connections = ConnectLeft | ConnectRight | ConnectTop | ConnectBottom
			case '@', '>', '<', '^':
				// Do nothing by intention
			case ',':
				cell.Type = TrackDiagonalLower
				cell.Connections = ConnectBottom | ConnectRight
				processCells(l, cellLines[y], x, scanRight, log)
				processCells(l, cellColumns[x], y, scanDownwards, log)
			case '.':
				cell.Type = TrackDiagonalBackLower
				cell.Connections = ConnectBottom | ConnectLeft
				processCells(l, cellLines[y], x, scanLeft, log)
				processCells(l, cellColumns[x], y, scanDownwards, log)
			case '`':
				cell.Type = TrackDiagonalBackUpper
				cell.Connections = ConnectTop | ConnectRight
				processCells(l, cellLines[y], x, scanRight, log)
				processCells(l, cellColumns[x], y, scanUpwards, log)
			case '\'':
				cell.Type = TrackDiagonalUpper
				cell.Connections = ConnectTop | ConnectLeft
				processCells(l, cellLines[y], x, scanLeft, log)
				processCells(l, cellColumns[x], y, scanUpwards, log)
			case '%':
				cell.Type = TrackDoubleSlash
				cell.Connections = ConnectTop | ConnectRight | ConnectBottom | ConnectRight
				processCells(l, cellLines[y], x, scanLeft, log)
				processCells(l, cellLines[y], x, scanRight, log)
				processCells(l, cellColumns[x], y, scanUpwards, log)
				processCells(l, cellColumns[x], y, scanDownwards, log)
			case '&':
				cell.Type = TrackDoubleBackslash
				cell.Connections = ConnectTop | ConnectRight | ConnectBottom | ConnectRight
				processCells(l, cellLines[y], x, scanLeft, log)
				processCells(l, cellLines[y], x, scanRight, log)
				processCells(l, cellColumns[x], y, scanUpwards, log)
				processCells(l, cellColumns[x], y, scanDownwards, log)
			case 0:
				// Do nothing by intention
			default:
				if unicode.IsSpace(cell.Rune) || unicode.IsLetter(cell.Rune) || unicode.IsDigit(cell.Rune) {
					// Do nothing by intention
				} else {
					addASCIIStructureError(log, errlog.ErrorIllegalRune, loc, x, y)
				}
			}
		}
	}

	for y := 0; y < l.LineCount; y++ {
		for x := 0; x < l.ColumnCount; x++ {
			cell := l.Cell(x, y)
			switch cell.Type {
			case SwitchHorizontalBackslash:
				if cell_above := l.CellAbove(x, y); cell_above != nil && cell_above.ConnectsToBottom() {
					cell.Connections = ConnectTop | ConnectLeft | ConnectRight
					cell.Type = SwitchHorizontalDiagonalBackUpper
				} else if cell_below := l.CellBelow(x, y); cell_below != nil && cell_below.ConnectsToTop() {
					cell.Connections |= ConnectBottom | ConnectLeft | ConnectRight
					cell.Type = SwitchHorizontalDiagonalBackLower
				} else {
					addASCIIStructureError(log, errlog.ErrorMalformedLayout, loc, x, y, "Switch has too few connections")
				}
			case SwitchHorizontalSlash:
				if cell_above := l.CellAbove(x, y); cell_above != nil && cell_above.ConnectsToBottom() {
					cell.Connections = ConnectTop | ConnectLeft | ConnectRight
					cell.Type = SwitchHorizontalDiagonalUpper
				} else if cell_below := l.CellBelow(x, y); cell_below != nil && cell_below.ConnectsToTop() {
					cell.Connections = ConnectBottom | ConnectLeft | ConnectRight
					cell.Type = SwitchHorizontalDiagonalLower
				} else {
					addASCIIStructureError(log, errlog.ErrorMalformedLayout, loc, x, y, "Switch has too few connections")
				}
			case SwitchVerticalBackslash:
				if cell_left := l.CellLeftOf(x, y); cell_left != nil && cell_left.ConnectsToRight() {
					cell.Connections = ConnectTop | ConnectBottom | ConnectRight
					cell.Type = SwitchVerticalDiagonalBackLower
				} else if cell_right := l.CellRightOf(x, y); cell_right != nil && cell_right.ConnectsToLeft() {
					cell.Connections = ConnectBottom | ConnectLeft | ConnectTop
					cell.Type = SwitchVerticalDiagonalBackUpper
				} else {
					addASCIIStructureError(log, errlog.ErrorMalformedLayout, loc, x, y, "Switch has too few connections")
				}
			case SwitchVerticalSlash:
				if cell_left := l.CellLeftOf(x, y); cell_left != nil && cell_left.ConnectsToRight() {
					cell.Connections = ConnectTop | ConnectBottom | ConnectRight
					cell.Type = SwitchVerticalDiagonalUpper
				} else if cell_right := l.CellRightOf(x, y); cell_right != nil && cell_right.ConnectsToLeft() {
					cell.Connections = ConnectBottom | ConnectLeft | ConnectTop
					cell.Type = SwitchVerticalDiagonalLower
				} else {
					addASCIIStructureError(log, errlog.ErrorMalformedLayout, loc, x, y, "Switch has too few connections")
				}
			case SwitchHorizontalSlash | SwitchVerticalSlash:
				addASCIIStructureError(log, errlog.ErrorMalformedLayout, loc, x, y, "Switch has too many connections")
			case SwitchHorizontalBackslash | SwitchVerticalBackslash:
				addASCIIStructureError(log, errlog.ErrorMalformedLayout, loc, x, y, "Switch has too many connections")
			}
		}
	}
	return l
}

func processCells(l *ASCIIStructure, cells []*ASCIICell, pos int, dir scanDirection, log *errlog.ErrorLog) {
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
		if dir == scanRight && i+1 < len(cells) && (cells[i+1].Rune == '/' || cells[i+1].Rune == '\\' || cells[i+1].Rune == '-' || cells[i+1].Rune == '\'' || cells[i+1].Rune == '.' || cells[i+1].Rune == '&' || cells[i+1].Rune == '%') {
			cell.Type |= SwitchHorizontalSlash
			processCells(l, cells, i, dir, log)
		} else if dir == scanLeft && i > 0 && (cells[i-1].Rune == '/' || cells[i-1].Rune == '\\' || cells[i-1].Rune == '-' || cells[i-1].Rune == ',' || cells[i-1].Rune == '`' || cells[i-1].Rune == '&' || cells[i-1].Rune == '%') {
			cell.Type |= SwitchHorizontalSlash
			processCells(l, cells, i, dir, log)
		} else if dir == scanDownwards && i+1 < len(cells) && (cells[i+1].Rune == '/' || cells[i+1].Rune == '\\' || cells[i+1].Rune == '|' || cells[i+1].Rune == '\'' || cells[i+1].Rune == '`' || cells[i+1].Rune == '&' || cells[i+1].Rune == '%') {
			cell.Type |= SwitchVerticalSlash
			processCells(l, cells, i, dir, log)
		} else if dir == scanUpwards && i > 0 && (cells[i-1].Rune == '/' || cells[i-1].Rune == '\\' || cells[i-1].Rune == '|' || cells[i-1].Rune == '.' || cells[i-1].Rune == ',' || cells[i-1].Rune == '&' || cells[i-1].Rune == '%') {
			cell.Type |= SwitchVerticalSlash
			processCells(l, cells, i, dir, log)
		}
		return
	case '\\':
		if dir == scanRight && i+1 < len(cells) && (cells[i+1].Rune == '\\' || cells[i+1].Rune == '/' || cells[i+1].Rune == '-' || cells[i+1].Rune == '\'' || cells[i+1].Rune == '.' || cells[i+1].Rune == '&' || cells[i+1].Rune == '%') {
			cell.Type |= SwitchHorizontalBackslash
			processCells(l, cells, i, dir, log)
		} else if dir == scanLeft && i > 0 && (cells[i-1].Rune == '\\' || cells[i-1].Rune == '/' || cells[i-1].Rune == '-' || cells[i-1].Rune == ',' || cells[i-1].Rune == '`' || cells[i-1].Rune == '&' || cells[i-1].Rune == '%') {
			cell.Type |= SwitchHorizontalBackslash
			processCells(l, cells, i, dir, log)
		} else if dir == scanDownwards && i+1 < len(cells) && (cells[i+1].Rune == '\\' || cells[i+1].Rune == '/' || cells[i+1].Rune == '|' || cells[i+1].Rune == '\'' || cells[i+1].Rune == '`' || cells[i+1].Rune == '&' || cells[i+1].Rune == '%') {
			cell.Type |= SwitchVerticalBackslash
			processCells(l, cells, i, dir, log)
		} else if dir == scanUpwards && i > 0 && (cells[i-1].Rune == '\\' || cells[i-1].Rune == '/' || cells[i-1].Rune == '|' || cells[i-1].Rune == '.' || cells[i-1].Rune == ',' || cells[i-1].Rune == '&' || cells[i-1].Rune == '%') {
			cell.Type |= SwitchVerticalBackslash
			processCells(l, cells, i, dir, log)
		}
		return
	case '@':
		if cell.Type != UnprocessedCell {
			addASCIIStructureError(log, errlog.ErrorMalformedLayout, l.Location, cell.X, cell.Y, "'@' has more than one track connection")
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
			return
			// l.addASCIIStructureError(errlog.ErrorMalformedLayout, cell.X, cell.Y, " has more than one track connection")
		}
		if cell.Rune == 'B' {
			if (dir == scanDownwards || dir == scanRight) && i+4 < len(cells) && cells[i+1].Rune == 'B' && cells[i+2].Rune == 'B' && cells[i+3].Rune == 'B' && !unicode.IsLetter(cells[i+4].Rune) && !unicode.IsDigit(cells[i+4].Rune) {
				if dir == scanDownwards {
					MakeBlock(cells[i:i+4], TrackVerticalBlock)
				} else {
					MakeBlock(cells[i:i+4], TrackHorizontalBlock)
				}
			} else if (dir == scanUpwards || dir == scanLeft) && i >= 4 && cells[i-1].Rune == 'B' && cells[i-2].Rune == 'B' && cells[i-3].Rune == 'B' && !unicode.IsLetter(cells[i-4].Rune) && !unicode.IsDigit(cells[i-4].Rune) {
				if dir == scanUpwards {
					MakeBlock(cells[i-3:i+1], TrackVerticalBlock)
				} else {
					MakeBlock(cells[i-3:i+1], TrackHorizontalBlock)
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
				MakeLabel(cells[start:j+1], TrackVerticalLabel)
			} else if dir == scanRight {
				MakeLabel(cells[start:j+1], TrackHorizontalLabel)
			} else if dir == scanUpwards {
				MakeLabel(cells[start:j+1], TrackVerticalLabel)
			} else {
				MakeLabel(cells[j:start+1], TrackHorizontalLabel)
			}
		}
	}
}

func addASCIIStructureError(log *errlog.ErrorLog, code errlog.ErrorCode, loc errlog.LocationRange, x int, y int, args ...string) *errlog.Error {
	newloc := errlog.EncodeLocationRange(loc.File(), loc.Line()+y, loc.Position()+x, loc.Line()+y, loc.Position()+x)
	return log.LogError(code, newloc, args...)
}
