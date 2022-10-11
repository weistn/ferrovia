package interpreter

import (
	"github.com/weistn/ferrovia/errlog"
	"github.com/weistn/ferrovia/model/tracks"
	"github.com/weistn/ferrovia/parser"
)

type pendingAnchor struct {
	x     float64
	y     float64
	z     float64
	angle float64
}

// Implements IContext
type TracksContext struct {
	layer *tracks.TrackLayer
	// track *tracks.Track
	first      *tracks.TrackConnection
	last       *tracks.TrackConnection
	anchor     *pendingAnchor
	atFunc     FuncValue
	layerFunc  FuncValue
	trackFuncs map[string]*FuncValue
}

func NewTracksContext(layer *tracks.TrackLayer) *TracksContext {
	ctx := &TracksContext{layer: layer}
	ctx.trackFuncs = make(map[string]*FuncValue)
	ctx.layerFunc = FuncValue{
		Name: "layer",
		Func: func(b *Interpreter, c []IContext, loc errlog.LocationRange, args ...parser.IExpression) (*ExprValue, *errlog.Error) {
			if len(args) != 1 {
				return nil, b.errlog.LogError(errlog.ErrorArgumentCountMismatch, loc, "1")
			}
			name, err := b.evalToString(c, args[0])
			if err != nil {
				return nil, err
			}
			l, ok := b.model.Tracks.Layers[name]
			if !ok {
				return nil, b.errlog.LogError(errlog.ErrorUnknownLayer, loc, name)
			}
			ctx.layer = l
			return nil, err
		},
	}
	ctx.atFunc = FuncValue{
		Name: "@",
		Func: func(b *Interpreter, c []IContext, loc errlog.LocationRange, args ...parser.IExpression) (*ExprValue, *errlog.Error) {
			if len(args) != 4 {
				return nil, b.errlog.LogError(errlog.ErrorArgumentCountMismatch, loc, "1")
			}
			x, err := b.evalToFloat(c, args[0])
			if err != nil {
				return nil, err
			}
			y, err := b.evalToFloat(c, args[1])
			if err != nil {
				return nil, err
			}
			z, err := b.evalToFloat(c, args[2])
			if err != nil {
				return nil, err
			}
			angle, err := b.evalToFloat(c, args[3])
			if err != nil {
				return nil, err
			}
			if ctx.last == nil {
				// Apply to the next track
				ctx.anchor = &pendingAnchor{x: x, y: y, z: z, angle: angle}
			} else {
				// Apply to previous track
				angle += 180
				if angle > 360 {
					angle -= 360
				}
				l := tracks.NewTrackLocation(ctx.last, tracks.Vec3{x, y, z}, angle)
				if !ctx.last.Track.SetLocation(l) {
					return nil, b.errlog.LogError(errlog.ErrorTrackPositionedTwice, loc)
				}
				b.tracksWithAnchor = append(b.tracksWithAnchor, ctx.last.Track)
			}
			return nil, nil
		},
	}
	return ctx
}

/*
// Returns (nil, false, nil) if the call could not be made, because the function is unknown.
func (c *TrackContext) Call(b *Interpreter, loc errlog.LocationRange, name string, args ...*ExprValue) (*ExprValue, bool, *errlog.Error) {
	if name == "layer" {
		if len(args) != 1 {
			return nil, true, errlog.NewError(errlog.ErrorArgumentCountMismatch, loc, "1")
		}
		name, err := b.ToString(args[0], loc)
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
		x, err := b.ToFloat(args[0], loc)
		if err != nil {
			return nil, true, err
		}
		y, err := b.ToFloat(args[1], loc)
		if err != nil {
			return nil, true, err
		}
		z, err := b.ToFloat(args[2], loc)
		if err != nil {
			return nil, true, err
		}
		angle, err := b.ToFloat(args[3], loc)
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
			// A track of this name does not exist.
			return nil, false, nil
		}
		newTrack.SourceLocation = loc
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
*/

func (c *TracksContext) Lookup(b *Interpreter, loc errlog.LocationRange, name string) (*ExprValue, *errlog.Error) {
	switch name {
	case "layer":
		return &ExprValue{Type: funcType, FuncValue: &c.layerFunc}, nil
	case "@":
		return &ExprValue{Type: funcType, FuncValue: &c.atFunc}, nil
	default:
		if f, ok := c.trackFuncs[name]; ok {
			return &ExprValue{Type: funcType, FuncValue: f}, nil
		}
		f := &FuncValue{}
		f.Name = name
		f.Func = func(b *Interpreter, ctx []IContext, loc errlog.LocationRange, args ...parser.IExpression) (*ExprValue, *errlog.Error) {
			if len(args) != 0 {
				return nil, b.errlog.LogError(errlog.ErrorArgumentCountMismatch, loc, "1")
			}
			newTrack := c.layer.NewTrack(name)
			if newTrack == nil {
				// A track of this name does not exist.
				return nil, b.errlog.LogError(errlog.ErrorUnknownTrackType, loc, name)
			}
			newTrack.SourceLocation = loc
			con := newTrack.FirstConnection()
			if c.anchor != nil {
				l := tracks.NewTrackLocation(con, tracks.Vec3{c.anchor.x, c.anchor.y, c.anchor.z}, c.anchor.angle)
				if !con.Track.SetLocation(l) {
					return nil, b.errlog.LogError(errlog.ErrorTrackPositionedTwice, loc)
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
			return nil, nil
		}
		c.trackFuncs[name] = f
		return &ExprValue{Type: funcType, FuncValue: f}, nil
	}
}
