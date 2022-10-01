package interpreter

import (
	"github.com/weistn/ferrovia/errlog"
	"github.com/weistn/ferrovia/model/tracks"
)

type pendingAnchor struct {
	x     float64
	y     float64
	z     float64
	angle float64
}

// Implements IContext
type TrackContext struct {
	layer *tracks.TrackLayer
	// track *tracks.Track
	first  *tracks.TrackConnection
	last   *tracks.TrackConnection
	anchor *pendingAnchor
}

// Returns (nil, false, nil) if the call could not be made, because the function is unknown.
func (c *TrackContext) Call(b *Interpreter, loc errlog.LocationRange, name string, args ...*ExprValue) (*ExprValue, bool, *errlog.Error) {
	if name == "layer" {
		if len(args) != 1 {
			return nil, true, errlog.NewError(errlog.ErrorArgumentCountMismatch, loc, "1")
		}
		name, err := args[0].ToString(loc)
		if err != nil {
			return nil, true, err
		}
		l, ok := b.model.Tracks.Layers[name]
		if !ok {
			return nil, true, errlog.NewError(errlog.ErrorUnknownLayer, loc, name)
		}
		c.layer = l
		return nil, true, nil
	} else if name == "@" {
		if len(args) != 4 {
			return nil, true, errlog.NewError(errlog.ErrorArgumentCountMismatch, loc, "1")
		}
		x, err := args[0].ToFloat(loc)
		if err != nil {
			return nil, true, err
		}
		y, err := args[1].ToFloat(loc)
		if err != nil {
			return nil, true, err
		}
		z, err := args[2].ToFloat(loc)
		if err != nil {
			return nil, true, err
		}
		angle, err := args[3].ToFloat(loc)
		if err != nil {
			return nil, true, err
		}
		if c.last == nil {
			// Apply to the next track
			c.anchor = &pendingAnchor{x: x, y: y, z: z, angle: angle}
		} else {
			// Apply to previous track
			angle += 180
			if angle > 360 {
				angle -= 360
			}
			l := tracks.NewTrackLocation(c.last, tracks.Vec3{x, y, z}, angle)
			if !c.last.Track.SetLocation(l) {
				return nil, true, errlog.NewError(errlog.ErrorTrackPositionedTwice, loc)
			}
			b.tracksWithAnchor = append(b.tracksWithAnchor, c.last.Track)
		}
		return nil, true, nil
	} else {
		newTrack := c.layer.NewTrack(name)
		if newTrack == nil {
			return nil, false, nil
		}
		newTrack.SourceLocation = loc
		if newTrack.Geometry.IncomingConnectionCount != 1 || newTrack.Geometry.OutgoingConnectionCount != 1 {
			panic("Track has more than one incoming or outgoing connections")
		}
		// Not a switch. Simply connect both ends
		con := newTrack.FirstConnection()
		if c.anchor != nil {
			l := tracks.NewTrackLocation(con, tracks.Vec3{c.anchor.x, c.anchor.y, c.anchor.z}, c.anchor.angle)
			if !con.Track.SetLocation(l) {
				return nil, true, errlog.NewError(errlog.ErrorTrackPositionedTwice, loc)
			}
			b.tracksWithAnchor = append(b.tracksWithAnchor, con.Track)
			c.anchor = nil
		}
		if c.last != nil {
			c.last.Connect(con)
		} else {
			c.first = con
		}
		c.last = newTrack.SecondConnection()
		return nil, true, nil
	}
}
