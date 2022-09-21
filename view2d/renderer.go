package view2d

import (
	"math"

	"github.com/weistn/ferrovia/parser"
	"github.com/weistn/ferrovia/model/tracks"
)

type Canvas struct {
	Name string `json:"name"`
	// Maps layer names to sequences of tracks
	Layers []*Layer              `json:"layers"`
	Width  float64               `json:"width"`
	Height float64               `json:"height"`
	Ground []*parser.GroundPlate `json:"ground"`
}

type Layer struct {
	Name   string   `json:"name"`
	Tracks []*Track `json:"tracks"`
}

type Track struct {
	Lines      []*Line      `json:"l,omitempty"`
	Arcs       []*Arc       `json:"a,omitempty"`
	Delimiters []*Delimiter `json:"d,omitempty"`
}

type Line struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Angle  float64 `json:"a"`
	Length float64 `json:"l"`
}

type Delimiter struct {
	X1 float64 `json:"x1"`
	Y1 float64 `json:"y1"`
	X2 float64 `json:"x2"`
	Y2 float64 `json:"y2"`
}

type Arc struct {
	CenterX    float64 `json:"cx"`
	CenterY    float64 `json:"cy"`
	Radius     float64 `json:"r"`
	StartAngle float64 `json:"sa"`
	TrackAngle float64 `json:"ta"`
}

func Render(ts *tracks.TrackSystem) *Canvas {
	tracks.NewEpoch()
	c := &Canvas{}
	c.Name = ts.Name
	for _, l := range ts.Layers {
		cl := &Layer{Name: l.Name}
		c.Layers = append(c.Layers, cl)
		for _, track := range l.Tracks {
			renderTrack(track, cl, c)
		}
	}
	for _, ground := range ts.Ground {
		c.Ground = append(c.Ground, ground)
		if len(ground.Polygon) != 0 {
			for _, p := range ground.Polygon {
				c.Width = math.Max(c.Width, ground.Left+p.X)
				c.Height = math.Max(c.Height, ground.Top+p.Y)
			}
		} else {
			c.Width = math.Max(c.Width, ground.Left+ground.Width)
			c.Height = math.Max(c.Height, ground.Top+ground.Height)
		}
	}
	return c
}

func renderTrack(track *tracks.Track, cl *Layer, c *Canvas) {
	t := &Track{}
	for _, path := range track.Geometry.Paths {
		switch p := path.(type) {
		case *tracks.TrackGeometryLine:
			from := track.Location.Center.Add2(p.Anchor.Position.Rotate(track.Location.Rotation))
			startAngle := track.Location.Rotation + p.Anchor.Angle
			length := p.Size
			t.Lines = append(t.Lines, &Line{X: from[0], Y: from[1], Angle: normalizeAngle(startAngle), Length: length})

			to := from.Add2(tracks.Vec2{0, -length}.Rotate(startAngle))
			c.Width = math.Max(math.Max(from[0]+100, c.Width), to[0]+100)
			c.Height = math.Max(math.Max(from[1]+100, c.Height), to[1]+100)
		case *tracks.TrackGeometryArc:
			from := track.Location.Center.Add2(p.Anchor.Position.Rotate(track.Location.Rotation))
			angle := track.Location.Rotation + p.Anchor.Angle
			trackAngle := p.TrackAngle
			radius := p.Radius
			sin := math.Sin(angle * math.Pi / 180)
			cos := math.Cos(angle * math.Pi / 180)
			x := cos
			y := sin
			centerX := from[0] + x*radius
			centerY := from[1] + y*radius
			t.Arcs = append(t.Arcs, &Arc{CenterX: centerX, CenterY: centerY, Radius: radius, StartAngle: normalizeAngle(angle - 90), TrackAngle: trackAngle})

			maxSize := 2 * math.Pi * radius * trackAngle / 360
			c.Width = math.Max(math.Max(from[0], c.Width), from[0]+maxSize)
			c.Height = math.Max(math.Max(from[1], c.Height), from[1]+maxSize)
		default:
			panic("Not implemented")
		}
	}
	track.Tag()
	for i := 0; i < track.ConnectionCount(); i++ {
		c := track.Connection(i)
		if c.Opposite != nil && c.Opposite.Track.IsTagged() {
			continue
		}
		pos := track.Location.Center.Add2(track.Geometry.ConnectionPoints[i].Position.Invert().Rotate(track.Location.Rotation))
		angle := track.Location.Rotation + track.Geometry.ConnectionPoints[i].Angle
		t.Delimiters = append(t.Delimiters, renderTrackDelimiter(pos, angle, 30))
	}
	cl.Tracks = append(cl.Tracks, t)
}

func renderTrackDelimiter(pos tracks.Vec3, angle float64, size float64) *Delimiter {
	sin := math.Sin(angle * math.Pi / 180)
	cos := math.Cos(angle * math.Pi / 180)
	x := cos * size / 2
	y := sin * size / 2
	return &Delimiter{X1: pos[0] + x, Y1: pos[1] + y, X2: pos[0] - x, Y2: pos[1] - y}
}

func normalizeAngle(a float64) float64 {
	if a < 0 {
		for a < 0 {
			a += 360
		}
	} else {
		for a >= 360 {
			a -= 360
		}
	}
	return a
}
