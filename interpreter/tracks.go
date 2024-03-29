package interpreter

import (
	"github.com/weistn/ferrovia/errlog"
	"github.com/weistn/ferrovia/model/tracks"
	"github.com/weistn/ferrovia/parser"
)

type pendingAnchor struct {
	x        float64
	y        float64
	z        float64
	angle    float64
	location errlog.LocationRange
}

// Implements IContext
type TracksContext struct {
	// The currently selected layer
	layer *tracks.TrackLayer
	// A list of *Track or *pendingAnchor instances.
	// The list is processed upon Close().
	elements []interface{}
	// Populate after Close()
	first *tracks.TrackConnection
	// Populate after Close()
	last      *tracks.TrackConnection
	atFunc    FuncValue
	layerFunc FuncValue
	// A cache
	trackFuncs map[string]*FuncValue
	location   errlog.LocationRange
}

// Implements IContext
type TurnoutContext struct {
	track      *tracks.Track
	left       *TracksContext
	right      *TracksContext
	middle     *TracksContext
	backleft   *TracksContext
	backright  *TracksContext
	backmiddle *TracksContext
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
			return nil, nil
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
			// ctx.elements = append(ctx.elements, &pendingAnchor{x: x, y: y, z: z, angle: angle, location: loc})
			// return nil, nil
			return &ExprValue{Type: contextType, Context: &ValueContext{Value: &pendingAnchor{x: x, y: y, z: z, angle: angle, location: loc}}}, nil
		},
	}
	return ctx
}

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
		// Is it a track type?
		_, ok := tracks.TrackFactories[name]
		if !ok {
			return nil, nil
		}
		f := &FuncValue{}
		f.Name = name
		f.Func = func(b *Interpreter, ctx []IContext, loc errlog.LocationRange, args ...parser.IExpression) (*ExprValue, *errlog.Error) {
			if len(args) != 0 {
				return nil, b.errlog.LogError(errlog.ErrorArgumentCountMismatch, loc, "0")
			}
			newTrack := c.layer.NewTrack(name)
			if newTrack == nil {
				// A track of this name does not exist.
				return nil, b.errlog.LogError(errlog.ErrorUnknownTrackType, loc, name)
			}
			newTrack.SourceLocation = loc
			// In case of a turnout, create a TurnoutContext
			if newTrack.Geometry.IncomingConnectionCount+newTrack.Geometry.OutgoingConnectionCount > 2 {
				return &ExprValue{Type: contextType, Context: NewTurnoutContext(newTrack)}, nil
			}
			return &ExprValue{Type: contextType, Context: &ValueContext{Value: newTrack}}, nil
		}
		c.trackFuncs[name] = f
		return &ExprValue{Type: funcType, FuncValue: f}, nil
	}
}

func (c *TracksContext) Process(b *Interpreter, loc errlog.LocationRange, value *ExprValue) *errlog.Error {
	if value.Type == contextType {
		switch t := value.Context.(type) {
		case *TracksContext:
			c.elements = append(c.elements, t)
			return nil
		case *TurnoutContext:
			c.elements = append(c.elements, t.track)
			return nil
		case *ValueContext:
			switch v := t.Value.(type) {
			case *tracks.Track:
				c.elements = append(c.elements, v)
				return nil
			case *pendingAnchor:
				c.elements = append(c.elements, v)
				return nil
			}
		}
	}
	return b.errlog.LogError(errlog.ErrorIllegalInThisContext, loc)
}

func (c *TracksContext) Close(b *Interpreter) *errlog.Error {
	//
	// Connect all tracks
	//
	var anchor *pendingAnchor
	for _, el := range c.elements {
		switch e := el.(type) {
		case *tracks.Track:
			con := e.FirstConnection()
			if anchor != nil {
				l := tracks.NewTrackLocation(con, tracks.Vec3{anchor.x, anchor.y, anchor.z}, anchor.angle)
				if !con.Track.SetLocation(l) {
					return b.errlog.LogError(errlog.ErrorTrackPositionedTwice, e.SourceLocation)
				}
				b.tracksWithAnchor = append(b.tracksWithAnchor, con.Track)
				anchor = nil
			}
			if c.last != nil {
				c.last.Connect(con)
			} else {
				c.first = con
			}
			c.last = e.SecondConnection()
		case *TracksContext:
			if e.first != nil {
				if c.last != nil {
					c.last.Connect(e.first)
				} else {
					c.first = e.first
				}
				c.last = e.last
			}
		case *pendingAnchor:
			if c.last == nil {
				// Apply to the next track
				anchor = e
			} else {
				// Apply to previous track
				angle := e.angle
				angle += 180
				if angle > 360 {
					angle -= 360
				}
				l := tracks.NewTrackLocation(c.last, tracks.Vec3{e.x, e.y, e.z}, angle)
				if !c.last.Track.SetLocation(l) {
					return b.errlog.LogError(errlog.ErrorTrackPositionedTwice, e.location)
				}
				b.tracksWithAnchor = append(b.tracksWithAnchor, c.last.Track)
			}
		default:
			panic("Ooooops")
		}
	}
	return nil
}

func NewTurnoutContext(track *tracks.Track) *TurnoutContext {
	return &TurnoutContext{track: track}
}

func (c *TurnoutContext) Lookup(b *Interpreter, loc errlog.LocationRange, name string) (*ExprValue, *errlog.Error) {
	switch name {
	case "left":
		c.left = NewTracksContext(c.track.Layer)
		c.left.location = loc
		return &ExprValue{Type: contextType, Context: c.left}, nil
	case "right":
		c.right = NewTracksContext(c.track.Layer)
		c.right.location = loc
		return &ExprValue{Type: contextType, Context: c.right}, nil
	case "middle":
		c.middle = NewTracksContext(c.track.Layer)
		c.middle.location = loc
		return &ExprValue{Type: contextType, Context: c.middle}, nil
	case "backleft":
		c.backleft = NewTracksContext(c.track.Layer)
		c.backleft.location = loc
		return &ExprValue{Type: contextType, Context: c.backleft}, nil
	case "backright":
		c.backright = NewTracksContext(c.track.Layer)
		c.backright.location = loc
		return &ExprValue{Type: contextType, Context: c.backright}, nil
	case "backmiddle":
		c.backmiddle = NewTracksContext(c.track.Layer)
		c.backmiddle.location = loc
		return &ExprValue{Type: contextType, Context: c.backmiddle}, nil
	}
	return nil, nil
}

func (c *TurnoutContext) connect(c1, c2 *tracks.TrackConnection) {
	if c2 == nil {
		return
	}
	c1.Connect(c2)
}

func (c *TurnoutContext) Process(b *Interpreter, loc errlog.LocationRange, value *ExprValue) *errlog.Error {
	// TODO: Throw error on unacceptable values
	return nil
}

func (c *TurnoutContext) Close(b *Interpreter) *errlog.Error {
	if c.track.Geometry.IncomingConnectionCount == 1 && c.track.Geometry.OutgoingConnectionCount == 2 {
		if c.middle != nil || c.backmiddle != nil {
			panic("TOOD: The 'middle' connection is not available on this track")
		}
		// A normal turnout
		if c.backleft != nil || c.backright != nil {
			// Reverse track
			if c.left != nil || c.right != nil {
				panic("TODO: Track is not a crossing")
			}
			if c.backleft != nil && c.backright != nil {
				panic("TODO: No free connection left on turnout")
			}
			c.track.Reverse()
			if c.backright != nil {
				c.connect(c.track.Connection(1), c.backright.last)
				c.track.SelectedTurnoutOption = 1
			} else {
				c.connect(c.track.Connection(2), c.backleft.last)
				c.track.SelectedTurnoutOption = 0
			}
		} else {
			if c.left != nil && c.right != nil {
				panic("TODO: No free connection left on turnout")
			}
			if c.left != nil {
				c.connect(c.track.Connection(1), c.left.first)
				c.track.SelectedTurnoutOption = 1
			} else if c.right != nil {
				c.connect(c.track.Connection(2), c.right.first)
				c.track.SelectedTurnoutOption = 0
			}
		}
	} else if c.track.Geometry.IncomingConnectionCount == 1 && c.track.Geometry.OutgoingConnectionCount == 3 {
		// A three-way turnout
		panic("TODO")
	} else if c.track.Geometry.IncomingConnectionCount == 2 && c.track.Geometry.OutgoingConnectionCount == 2 && len(c.track.Geometry.TurnoutOptions) == 2 {
		// A non-switching crossing
		if c.middle != nil || c.backmiddle != nil {
			panic("TOOD: The 'middle' connection is not available on this track")
		}
		if c.backleft != nil && c.backright != nil {
			panic("TODO: No free connection left on turnout")
		}
		if c.left != nil && c.right != nil {
			panic("TODO: No free connection left on turnout")
		}
		// A switching crossing
		if c.backright == nil && c.left == nil {
			c.track.SelectedTurnoutOption = 0
		} else if c.backleft == nil && c.right == nil {
			c.track.SelectedTurnoutOption = 1
		} else {
			panic("TODO: No free connection left on turnout")
		}
		if c.left != nil {
			c.connect(c.track.Connection(2), c.left.first)
		} else if c.right != nil {
			c.connect(c.track.Connection(3), c.right.first)
		}
		if c.backright != nil {
			c.connect(c.track.Connection(0), c.backright.last)
		} else if c.backleft != nil {
			c.connect(c.track.Connection(1), c.backleft.last)
		}
	} else if c.track.Geometry.IncomingConnectionCount == 2 && c.track.Geometry.OutgoingConnectionCount == 2 {
		if c.middle != nil || c.backmiddle != nil {
			panic("TOOD: The 'middle' connection is not available on this track")
		}
		// A switching crossing
		if c.backright == nil && c.left == nil {
			c.track.SelectedTurnoutOption = 0
		} else if c.backright == nil && c.right == nil {
			c.track.SelectedTurnoutOption = 1
		} else if c.backleft == nil && c.left == nil {
			c.track.SelectedTurnoutOption = 2
		} else if c.backleft == nil && c.right == nil {
			c.track.SelectedTurnoutOption = 3
		}
		if c.backleft != nil && c.backright != nil {
			panic("TODO: No free connection left on turnout")
		}
		if c.left != nil && c.right != nil {
			panic("TODO: No free connection left on turnout")
		}
		if c.left != nil {
			c.connect(c.track.Connection(2), c.left.first)
		} else if c.right != nil {
			c.connect(c.track.Connection(3), c.right.first)
		}
		if c.backright != nil {
			c.connect(c.track.Connection(0), c.backright.last)
		} else if c.backleft != nil {
			c.connect(c.track.Connection(1), c.backleft.last)
		}
	} else {
		panic("Ooops")
	}
	return nil
}
