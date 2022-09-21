package switchtower

import (
	"github.com/weistn/ferrovia/model/structure"
)

type LayoutDescription struct {
	RowCount    int            `json:"rows"`
	ColumnCount int            `json:"columns"`
	Tracks      []*LayoutTrack `json:"tracks"`
}

type LayoutTrack struct {
	X    int                        `json:"c"`
	Y    int                        `json:"r"`
	Kind structure.ASCIICellType `json:"kind"`
}

func Render(layouts []*structure.ASCIIStructure) *LayoutDescription {
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
				if c.Anchor != nil {
					// Do nothing
				} else if c.Type == structure.TrackDoubleSlash {
					d.Tracks = append(d.Tracks, &LayoutTrack{X: x, Y: y, Kind: structure.TrackDiagonalLower})
					d.Tracks = append(d.Tracks, &LayoutTrack{X: x, Y: y, Kind: structure.TrackDiagonalUpper})
				} else if c.Type == structure.TrackDoubleBackslash {
					d.Tracks = append(d.Tracks, &LayoutTrack{X: x, Y: y, Kind: structure.TrackDiagonalBackLower})
					d.Tracks = append(d.Tracks, &LayoutTrack{X: x, Y: y, Kind: structure.TrackDiagonalBackUpper})
				} else if c.Type != structure.UnprocessedCell && c.Type < 100 {
					d.Tracks = append(d.Tracks, &LayoutTrack{X: x, Y: y, Kind: c.Type})
				}
			}
		}
	}
	return d
}
