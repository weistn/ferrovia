package interpreter

import (
	"github.com/weistn/ferrovia/errlog"
	"github.com/weistn/ferrovia/parser"
	"github.com/weistn/ferrovia/tracks"
)

type Interpreter struct {
	errlog           *errlog.ErrorLog
	ast              *parser.File
	trackSystem      *tracks.TrackSystem
	tracksWithAnchor []*tracks.Track
}

type junctionInfo struct {
	expr   *parser.JunctionExpression
	first  bool
	second bool
}

func NewInterpreter(errlog *errlog.ErrorLog) *Interpreter {
	return &Interpreter{errlog: errlog}
}

func (b *Interpreter) Process(ast *parser.File) *tracks.TrackSystem {
	b.trackSystem = tracks.NewTrackSystem()
	b.ast = ast
	for _, s := range ast.Statements {
		if s == nil {
			break
		}
		if s.GroundPlate != nil {
			b.trackSystem.Ground = append(b.trackSystem.Ground, s.GroundPlate)
		} else if s.Layer != nil {
			b.processLayer(s.Layer)
		}
	}

	for _, s := range ast.Statements {
		if s == nil {
			break
		}
		if s.GroundPlate != nil {
			// Do nothing by intention
		} else if s.Layer != nil {
			// Do nothing by intention
		} else if s.RailWay != nil {
			b.processRailWay(s.RailWay)
		} else {
			panic("Oooops")
		}
	}

	// Determine the location of all tracks which are directly or indirectly anchored
	for _, t := range b.tracksWithAnchor {
		b.computeLocationFromAnchor(t)
	}

	// All tracks should be located by now
	for _, l := range b.trackSystem.Layers {
		for _, t := range l.Tracks {
			if t.Location == nil {
				b.errlog.LogError(errlog.ErrorTracksWithoutPosition, t.SourceLocation)
				break
			}
		}
	}

	return b.trackSystem
}

func (b *Interpreter) processLayer(ast *parser.Layer) {
	if _, ok := b.trackSystem.Layers[ast.Name]; ok {
		b.errlog.LogError(errlog.ErrorDuplicateLayer, ast.Location, ast.Name)
		return
	}
	b.trackSystem.NewLayer(ast.Name)
}

func (b *Interpreter) processRailWay(ast *parser.RailWay) {
	b.processExpressions(ast, ast.Expressions, nil)
}

// `start` may be nil. In this case the resulting tracks form a block of tracks that has two unconnected ends.
// Otherwiese the new tracks will be connected to `start`.
func (b *Interpreter) processExpressions(rail *parser.RailWay, expressions []*parser.Expression, start *tracks.TrackConnection) (first *tracks.TrackConnection, last *tracks.TrackConnection) {
	var queue []*parser.Expression
	current := start
	for i, t := range expressions {
		if t == nil {
			break
		}
		var c *tracks.TrackConnection
		if t.ConnectionMark != nil {
			if current == nil {
				queue = expressions[:i+1]
			} else {
				b.processConnectionMark(t.ConnectionMark, current)
			}
		} else if t.TrackMark != nil {
			if current == nil {
				queue = expressions[:i+1]
			} else {
				b.processTrackMark(t.TrackMark, current.Track)
			}
		} else if t.Track != nil {
			c1, c2, err := b.processTrack(rail, t.Track, current)
			if err == nil {
				c = c1
				last = c2
			}
		} else if t.Repeat != nil {
			c, last = b.processRepeatExpression(rail, t.Repeat, current)
		} else if t.Anchor != nil {
			if current == nil {
				queue = expressions[:i+1]
			} else {
				b.processAnchor(t.Anchor, current)
			}
		} else {
			panic("Oooops")
		}
		// This is the first track?
		if c != nil && first == nil {
			first = c
			if len(queue) != 0 {
				b.processExpressions(rail, queue, c)
			}
		}
		if last != nil {
			current = last
		}
	}
	return
}

func (b *Interpreter) processTrack(rail *parser.RailWay, ast *parser.TrackExpression, prevCon *tracks.TrackConnection) (first *tracks.TrackConnection, last *tracks.TrackConnection, err *errlog.Error) {
	l, ok := b.trackSystem.Layers[rail.Layer]
	if !ok {
		err = b.errlog.LogError(errlog.ErrorUnknownLayer, ast.Location, rail.Layer)
		return nil, nil, err
	}
	newTrack := l.NewTrack(ast.Type)
	if newTrack == nil {
		err = b.errlog.LogError(errlog.ErrorUnknownTrackType, ast.Location, ast.Type)
		return nil, nil, err
	}

	// Junctions? Then it must be a switch
	if len(ast.JunctionsOnLeft) != 0 || len(ast.JunctionsOnRight) != 0 {
		// A list of junctions, clock-wise around the switch, starting with tracks coming in from top left.
		var junctions []junctionInfo
		var incoming, outgoing int
		for _, j := range ast.JunctionsOnLeft {
			if j.Arrow == parser.TokenArrowLeftIn {
				junctions = append(junctions, junctionInfo{expr: j})
				incoming++
			}
		}
		junctions = append(junctions, junctionInfo{first: true})
		incoming++
		for _, j := range ast.JunctionsOnRight {
			if j.Arrow == parser.TokenArrowRightIn {
				junctions = append(junctions, junctionInfo{expr: j})
				incoming++
			}
		}
		for _, j := range ast.JunctionsOnRight {
			if j.Arrow == parser.TokenArrowRightOut {
				junctions = append(junctions, junctionInfo{expr: j})
				outgoing++
			}
		}
		junctions = append(junctions, junctionInfo{second: true})
		outgoing++
		for _, j := range ast.JunctionsOnLeft {
			if j.Arrow == parser.TokenArrowLeftOut {
				junctions = append(junctions, junctionInfo{expr: j})
				outgoing++
			}
		}

		var offset int
		if outgoing < incoming {
			offset = newTrack.Geometry.IncomingConnectionCount
			tmp := incoming
			incoming = outgoing
			outgoing = tmp
		}
		if newTrack.Geometry.IncomingConnectionCount != incoming || newTrack.Geometry.OutgoingConnectionCount != outgoing {
			// println(newTrack.Geometry.IncomingConnectionCount, incoming, newTrack.Geometry.OutgoingConnectionCount, outgoing)
			err = b.errlog.LogError(errlog.ErrorMismatchJunctionCount, ast.Location, ast.Type)
			return nil, nil, err
		}
		var connectionCount = newTrack.ConnectionCount()
		for i, j := range junctions {
			if j.expr != nil {
				b.processExpressions(rail, j.expr.Expressions, newTrack.Connection((i+offset)%connectionCount))
			} else if j.first {
				first = newTrack.Connection((i + offset) % connectionCount)
				if prevCon != nil {
					prevCon.Connect(first)
				}
			} else {
				last = newTrack.Connection((i + offset) % connectionCount)
			}
		}
	} else {
		if newTrack.Geometry.IncomingConnectionCount != 1 || newTrack.Geometry.OutgoingConnectionCount != 1 {
			err = b.errlog.LogError(errlog.ErrorMismatchJunctionCount, ast.Location, ast.Type)
			return nil, nil, err
		}
		// Not a switch. Simply connect both ends
		first = newTrack.FirstConnection()
		if prevCon != nil {
			prevCon.Connect(first)
		}
		last = newTrack.SecondConnection()
	}

	return
}

func (b *Interpreter) processConnectionMark(ast *parser.ConnectionMarkExpression, con *tracks.TrackConnection) {
	if m := b.trackSystem.GetMark(ast.Name); m != nil && m.Connection != nil {
		if m.Connection.IsConnected() {
			b.errlog.LogError(errlog.ErrorTrackConnectedTwice, ast.Location)
		} else {
			m.Connection.Connect(con)
		}
		return
	}
	if !con.AddMark(ast.Name) {
		b.errlog.LogError(errlog.ErrorTrackMarkDefinedTwice, ast.Location)
	}
}

func (b *Interpreter) processTrackMark(ast *parser.TrackMarkExpression, track *tracks.Track) {
	if !track.AddMark(ast.Position, ast.Name) {
		b.errlog.LogError(errlog.ErrorTrackMarkDefinedTwice, ast.Location)
	}
}

func (b *Interpreter) processRepeatExpression(rail *parser.RailWay, repeat *parser.RepeatExpression, prevCon *tracks.TrackConnection) (first *tracks.TrackConnection, last *tracks.TrackConnection) {
	last = prevCon
	for i := 0; i < repeat.Count; i++ {
		f, l := b.processExpressions(rail, repeat.TrackExpressions, last)
		if first == nil {
			first = f
		}
		last = l
	}
	return
}

func (b *Interpreter) processAnchor(ast *parser.AnchorExpression, con *tracks.TrackConnection) {
	l := tracks.NewTrackLocation(con, tracks.Vec3{ast.X, ast.Y, ast.Z}, ast.Angle)
	if !con.Track.SetLocation(l) {
		b.errlog.LogError(errlog.ErrorTrackPositionedTwice, ast.Location)
	}
	b.tracksWithAnchor = append(b.tracksWithAnchor, con.Track)
}

func (b *Interpreter) computeLocationFromAnchor(track *tracks.Track) {
	if track.Location == nil {
		panic("Track has no anchor")
	}
	tracks.NewEpoch()
	track.Tag()
	// println(track.Geometry.Name, track.Id, track.Location.Center[0], track.Location.Center[1], track.Location.Rotation)
	b.computeLocationOfConnectedTracks(track)
}

func (b *Interpreter) computeLocationOfConnectedTracks(track *tracks.Track) {
	for i := 0; i < track.ConnectionCount(); i++ {
		c := track.Connection(i)
		c2 := c.Opposite
		if c2 == nil || c2.Track.IsTagged() {
			continue
		}
		if c2.Track.Location != nil {
			b.errlog.LogError(errlog.ErrorTrackPositionedTwice, c2.Track.SourceLocation)
		}
		cpos, cangle := track.Location.Connection(i, track.Geometry)
		// println("    con at ", cpos[0], cpos[1], cangle, i)
		l := tracks.NewTrackLocation(c2, cpos, cangle)
		// println(c2.Track.Geometry.Name, c2.Track.Id, l.Center[0], l.Center[1], l.Rotation)
		c2.Track.SetLocation(l)
		c2.Track.Tag()
		b.computeLocationOfConnectedTracks(c2.Track)
	}
}
