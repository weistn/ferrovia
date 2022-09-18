package viewlayout

import "github.com/weistn/ferrovia/interpreter"

type LayoutDescription struct {
	RowCount    int            `json:"rows"`
	ColumnCount int            `json:"columns"`
	Tracks      []*LayoutTrack `json:"tracks"`
}

type LayoutTrack struct {
	X    int                        `json:"c"`
	Y    int                        `json:"r"`
	Kind interpreter.LayoutCellType `json:"kind"`
}

func Render(layouts []*interpreter.Layout) *LayoutDescription {
	d := &LayoutDescription{}
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
				if c.Type == interpreter.TrackDoubleSlash {
					d.Tracks = append(d.Tracks, &LayoutTrack{X: x, Y: y, Kind: interpreter.TrackDiagonalLower})
					d.Tracks = append(d.Tracks, &LayoutTrack{X: x, Y: y, Kind: interpreter.TrackDiagonalUpper})
				} else if c.Type == interpreter.TrackDoubleBackslash {
					d.Tracks = append(d.Tracks, &LayoutTrack{X: x, Y: y, Kind: interpreter.TrackDiagonalBackLower})
					d.Tracks = append(d.Tracks, &LayoutTrack{X: x, Y: y, Kind: interpreter.TrackDiagonalBackUpper})
				} else if c.Type != interpreter.UnprocessedCell && c.Type < 100 {
					d.Tracks = append(d.Tracks, &LayoutTrack{X: x, Y: y, Kind: c.Type})
				}
			}
		}
	}
	return d
}
