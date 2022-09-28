package switchboard

import (
	sb "github.com/weistn/ferrovia/model/switchboard"
)

type TrackDiagram struct {
	RowCount    int      `json:"rows"`
	ColumnCount int      `json:"columns"`
	Tracks      []*Track `json:"tracks"`
}

type Track struct {
	X    int              `json:"c"`
	Y    int              `json:"r"`
	Kind sb.ASCIICellType `json:"kind"`
}

func Render(layouts []*sb.ASCIISwitchboard) *TrackDiagram {
	d := &TrackDiagram{}
	for _, layout := range layouts {
		if layout.ColumnCount > d.ColumnCount {
			d.ColumnCount = layout.ColumnCount
		}
		if layout.LineCount > d.RowCount {
			d.RowCount = layout.LineCount
		}
		for y := 0; y < layout.LineCount; y++ {
			for x := 0; x < layout.ColumnCount; x++ {
				c := layout.Cell(x, y)
				if c.Anchor != nil {
					// Do nothing
				} else if c.Type == sb.TrackDoubleSlash {
					d.Tracks = append(d.Tracks, &Track{X: x, Y: y, Kind: sb.TrackDiagonalLower})
					d.Tracks = append(d.Tracks, &Track{X: x, Y: y, Kind: sb.TrackDiagonalUpper})
				} else if c.Type == sb.TrackDoubleBackslash {
					d.Tracks = append(d.Tracks, &Track{X: x, Y: y, Kind: sb.TrackDiagonalBackLower})
					d.Tracks = append(d.Tracks, &Track{X: x, Y: y, Kind: sb.TrackDiagonalBackUpper})
				} else if c.Type != sb.UnprocessedCell && c.Type < 100 {
					d.Tracks = append(d.Tracks, &Track{X: x, Y: y, Kind: c.Type})
				}
			}
		}
	}
	return d
}
