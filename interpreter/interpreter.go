package interpreter

import (
	"strings"

	"github.com/weistn/ferrovia/errlog"
	"github.com/weistn/ferrovia/model"
	"github.com/weistn/ferrovia/model/tracks"
	"github.com/weistn/ferrovia/parser"
)

type Interpreter struct {
	errlog           *errlog.ErrorLog
	ast              *parser.File
	model            *model.Model
	tracksWithAnchor []*tracks.Track
	identifiers      map[string]interface{}
}

func NewInterpreter(errlog *errlog.ErrorLog) *Interpreter {
	m := &model.Model{}
	m.Tracks = tracks.NewTrackSystem()
	return &Interpreter{errlog: errlog, model: m, identifiers: make(map[string]interface{})}
}

func (b *Interpreter) ProcessStatics(ast *parser.File) *model.Model {
	b.ast = ast

	// Determine all identifiers
	for _, s := range ast.Statements {
		if s == nil {
			break
		}
		switch t := s.(type) {
		case *parser.GroundPlate:
			// Do nothing by intention
		case *parser.Layer:
			if t.Name.StringValue != "" {
				b.identifiers[t.Name.StringValue] = s
			}
		case *parser.Switchboard:
			// Do nothing by intention
		case *parser.Tracks:
			if t.Name.StringValue != "" {
				b.identifiers[t.Name.StringValue] = s
			}
		default:
			panic("Ooooops")
		}
	}

	// Compute all ground plates and layers
	for _, s := range ast.Statements {
		if s == nil {
			break
		}
		switch t := s.(type) {
		case *parser.GroundPlate:

		case *parser.Layer:
			b.processLayer(t)
		case *parser.Switchboard:
			b.processSwitchboard(t)
		case *parser.Tracks:
			// Do nothing by intention
		default:
			panic("Ooooops")
		}
	}

	// Compute all tracks
	for _, s := range ast.Statements {
		if s == nil {
			break
		}
		switch t := s.(type) {
		case *parser.GroundPlate:
			// Do nothing by intention
		case *parser.Layer:
			// Do nothing by intention
		case *parser.Switchboard:
			// Do nothing by intention
		case *parser.Tracks:
			b.processTracks(t)
		default:
			panic("Ooooops")
		}
	}

	// Determine the location of all tracks which are directly or indirectly anchored
	for _, t := range b.tracksWithAnchor {
		b.computeLocationFromAnchor(t)
	}

	// All tracks should be located by now
	for _, l := range b.model.Tracks.Layers {
		for _, t := range l.Tracks {
			if t.Location == nil {
				b.errlog.LogError(errlog.ErrorTracksWithoutPosition, t.SourceLocation)
				break
			}
		}
	}

	return b.model
}

func (b *Interpreter) processSwitchboard(ast *parser.Switchboard) {
	lines := strings.Split(ast.RawText, "\n")
	sb := processASCIIStructure(lines, ast.LocationText, b.errlog)
	if sb != nil {
		b.model.Switchboards = append(b.model.Switchboards, sb)
	}
}

func (b *Interpreter) processLayer(ast *parser.Layer) {
	if _, ok := b.model.Tracks.Layers[ast.Name.StringValue]; ok {
		b.errlog.LogError(errlog.ErrorDuplicateLayer, ast.Location, ast.Name.StringValue)
		return
	}
	// TODO: execute expressions
	b.model.Tracks.NewLayer(ast.Name.StringValue)
}

func (b *Interpreter) processTracks(ast *parser.Tracks) {
	// TODO: execute expressions
}

/*
func (b *Interpreter) processSwitchExpression(rail *parser.RailWay, exp *parser.Expression, columns []*railwayColumn, index int) (columnsNew []*railwayColumn, first *tracks.TrackConnection, last *tracks.TrackConnection) {

	l, ok := b.trackSystem.Layers[rail.Layer]
	if !ok {
		b.errlog.LogError(errlog.ErrorUnknownLayer, exp.Location, rail.Layer)
		return columns, columns[index].first, columns[index].last
	}
	newTrack := l.NewTrack(trackType)
	if newTrack == nil {
		b.errlog.LogError(errlog.ErrorUnknownTrackType, exp.Location, exp.Switch.Type)
		return columns, columns[index].first, columns[index].last
	}
	newTrack.SourceLocation = exp.Location

}
*/

/*
// `start` may be nil. In this case the resulting tracks form a block of tracks that has two unconnected ends.
// Otherwiese the new tracks will be connected to `start`.
func (b *Interpreter) processExpression(rail *parser.RailWay, exp *parser.Expression, start *tracks.TrackConnection) (first *tracks.TrackConnection, last *tracks.TrackConnection) {
	var c *tracks.TrackConnection
	if exp.Track != nil {
		c1, c2, err := b.processTrack(rail, exp, start)
		if err == nil {
			c = c1
			last = c2
		}
	} else if exp.Repeat != nil {
		c, last = b.processRepeatExpression(rail, exp, start)
	} else {
		panic("Oooops")
	}
	// This is the first track?
	if first == nil {
		first = c
	}
	return
}
*/

/*
func (b *Interpreter) processTrack(rail *parser.RailWay, ast *parser.Expression, prevCon *tracks.TrackConnection) (first *tracks.TrackConnection, last *tracks.TrackConnection, err *errlog.Error) {
	if ast.Track == nil {
		panic("Not a track")
	}
	l, ok := b.trackSystem.Layers[rail.Layer]
	if !ok {
		err = b.errlog.LogError(errlog.ErrorUnknownLayer, ast.Location, rail.Layer)
		return nil, nil, err
	}

	if named, ok := b.namedRailways[ast.Track.Type]; ok {
		if named.processed {
			err = b.errlog.LogError(errlog.ErrorNamedRailwayUsedTwice, ast.Location, ast.Track.Type)
			return nil, nil, err
		}
		b.processRailWay(named.ast)
		if named.column.first != nil && prevCon != nil {
			prevCon.Connect(named.column.first)
		}
		return named.column.first, named.column.last, nil
	}

	newTrack := l.NewTrack(ast.Track.Type)
	if newTrack == nil {
		err = b.errlog.LogError(errlog.ErrorUnknownTrackType, ast.Location, ast.Track.Type)
		return nil, nil, err
	}
	newTrack.SourceLocation = ast.Location
	if newTrack.Geometry.IncomingConnectionCount != 1 || newTrack.Geometry.OutgoingConnectionCount != 1 {
		panic("Track has more than one incoming or outgoing connections")
	}
	// Not a switch. Simply connect both ends
	first = newTrack.FirstConnection()
	if prevCon != nil {
		prevCon.Connect(newTrack.FirstConnection())
	}
	last = newTrack.SecondConnection()
	return
}
*/

/*
func (b *Interpreter) processTrackMark(exp *parser.Expression, track *tracks.Track) {
	if exp.TrackMark == nil {
		panic("Not a track")
	}
	if !track.AddMark(exp.TrackMark.Position, exp.TrackMark.Name) {
		b.errlog.LogError(errlog.ErrorTrackMarkDefinedTwice, exp.Location)
	}
}
*/

/*
func (b *Interpreter) processAnchor(exp *parser.Expression, con *tracks.TrackConnection) {
	if exp.Anchor == nil {
		panic("Oooops")
	}
	ast := exp.Anchor
	l := tracks.NewTrackLocation(con, tracks.Vec3{ast.X, ast.Y, ast.Z}, ast.Angle)
	if !con.Track.SetLocation(l) {
		b.errlog.LogError(errlog.ErrorTrackPositionedTwice, exp.Location)
	}
	b.tracksWithAnchor = append(b.tracksWithAnchor, con.Track)
}
*/

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
