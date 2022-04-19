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
	namedRailways    map[string]*namedRailway
}

type railwayColumn struct {
	first    *tracks.TrackConnection
	last     *tracks.TrackConnection
	queue    []*parser.Expression
	marked   bool
	location errlog.LocationRange
}

type namedRailway struct {
	ast       *parser.RailWay
	processed bool
	column    *railwayColumn
}

func NewInterpreter(errlog *errlog.ErrorLog) *Interpreter {
	return &Interpreter{errlog: errlog, namedRailways: make(map[string]*namedRailway)}
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
		} else if s.RailWay != nil && s.RailWay.Name != "" {
			b.namedRailways[s.RailWay.Name] = &namedRailway{ast: s.RailWay}
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
		} else if s.RailWay != nil && s.RailWay.Name == "" {
			b.processRailWay(s.RailWay)
		} else if s.RailWay != nil {
			// Do nothing by intention
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
	var columns []*railwayColumn
	var named *namedRailway
	if ast.Name != "" {
		named = b.namedRailways[ast.Name]
	}

	for _, row := range ast.Rows {
		for _, exp := range row.Expressions {
			if exp.Switch != nil {
				// println("Switch at", exp.Location.Position())
				index := getColumnAtPosition(columns, exp.Location.Position())
				if index == -1 {
					if len(columns) == 0 && named != nil {
						columns, index = startOrEndColumnAtPosition(columns, exp.Location)
						named.column = columns[index]
					} else {
						b.errlog.LogError(errlog.ErrorColumnStartedUnexpectedly, exp.Location)
						continue
					}
				}
				col := columns[index]
				col.marked = true
				var first *tracks.TrackConnection
				columns, first, col.last = b.processSwitchExpression(ast, exp, columns, index)
				if col.first == nil {
					col.first = first
				}
				b.processColumnQueue(col)
			} else if exp.TrackMark != nil || exp.Anchor != nil {
				// `"Mark"` or `@(...)`. Queue it up
				index := getColumnAtPosition(columns, exp.Location.Position())
				if index == -1 {
					if len(columns) == 0 && named != nil {
						columns, index = startOrEndColumnAtPosition(columns, exp.Location)
						named.column = columns[index]
					} else {
						b.errlog.LogError(errlog.ErrorColumnStartedUnexpectedly, exp.Location)
						continue
					}
				}
				// println("Mark/Anchor at", exp.Location.Position())
				columns[index].marked = true
				columns[index].queue = append(columns[index].queue, exp)
			} else if exp.TrackTermination != nil && exp.TrackTermination.EllipsisLeft {
				// `... Label`. Start a new columns and queue up the epxression
				index := getColumnAtPosition(columns, exp.Location.Position())
				if index != -1 {
					b.errlog.LogError(errlog.ErrorColumnStartedUnexpectedly, exp.Location)
					continue
				}
				columns, index = startOrEndColumnAtPosition(columns, exp.Location)
				columns[index].marked = true
				columns[index].queue = append(columns[index].queue, exp)
			} else if exp.TrackTermination != nil && exp.TrackTermination.EllipsisRight {
				// `Label ...`. Terminate the column
				// println("Label at", exp.Location.Position())
				index := getColumnAtPosition(columns, exp.Location.Position())
				if index == -1 {
					b.errlog.LogError(errlog.ErrorColumnStartedUnexpectedly, exp.Location)
					continue
				}
				columns[index].marked = true
				if columns[index].last == nil {
					b.errlog.LogError(errlog.ErrorColumnHasNoTracks, exp.Location)
					continue
				}
				b.processTrackTerminationWithConnection(exp, columns[index].last)
				columns, _ = startOrEndColumnAtPosition(columns, exp.Location)
			} else if exp.TrackTermination != nil {
				// `end`. Terminate the column
				// println("END at", exp.Location.Position())
				columns, _ = startOrEndColumnAtPosition(columns, exp.Location)
			} else if exp.Placeholder {
				// `|`
				index := getColumnAtPosition(columns, exp.Location.Position())
				if index == -1 {
					if len(columns) == 0 && named != nil {
						columns, index = startOrEndColumnAtPosition(columns, exp.Location)
						named.column = columns[index]
					} else {
						b.errlog.LogError(errlog.ErrorColumnStartedUnexpectedly, exp.Location)
						continue
					}
				}
				columns[index].marked = true
			} else {
				// println("TRACK at", exp.Location.Position())
				index := getColumnAtPosition(columns, exp.Location.Position())
				if index == -1 {
					if len(columns) == 0 && named != nil {
						columns, index = startOrEndColumnAtPosition(columns, exp.Location)
						named.column = columns[index]
					} else {
						b.errlog.LogError(errlog.ErrorColumnStartedUnexpectedly, exp.Location)
						continue
					}
				}
				columns[index].marked = true
				var first *tracks.TrackConnection
				first, columns[index].last = b.processExpression(ast, exp, columns[index].last)
				if columns[index].first == nil {
					columns[index].first = first
				}
				b.processColumnQueue(columns[index])
			}
		}
		b.checkForUnmarkedColumms(columns)
	}

	// Unnamed railways must properly terminate all columns
	if named != nil {
		if len(columns) != 1 {
			b.errlog.LogError(errlog.ErrorIllegalColumnCountInNamedRailway, ast.Location)
		}
		if len(columns) == 1 && columns[0] != named.column {
			b.errlog.LogError(errlog.ErrorIllegalStartEndColumnInNamedRailway, ast.Location)
		}
		named.processed = true
	} else {
		for _, c := range columns {
			b.errlog.LogError(errlog.ErrorColumnTerminatedUnexpectedly, c.location)
		}
	}
}

func (b *Interpreter) processColumnQueue(column *railwayColumn) {
	if len(column.queue) == 0 {
		return
	}
	if column.first == nil {
		panic("Oooops")
	}
	for _, exp := range column.queue {
		if exp.TrackMark != nil {
			b.processTrackMark(exp, column.first.Track)
		} else if exp.Anchor != nil {
			b.processAnchor(exp, column.first)
		} else if exp.TrackTermination != nil {
			b.processTrackTerminationWithConnection(exp, column.first)
		}
	}
	column.queue = nil
}

func (b *Interpreter) checkForUnmarkedColumms(columns []*railwayColumn) {
	for i := 0; i < len(columns); i++ {
		if !columns[i].marked {
			// println("CHECK", i, columns[i].location.Position())
			b.errlog.LogError(errlog.ErrorColumnTerminatedUnexpectedly, columns[i].location)
		} else {
			columns[i].marked = false
		}
	}
}

func (b *Interpreter) processSwitchExpression(rail *parser.RailWay, exp *parser.Expression, columns []*railwayColumn, index int) (columnsNew []*railwayColumn, first *tracks.TrackConnection, last *tracks.TrackConnection) {
	trackType := exp.Switch.Type
	if exp.Switch.JoinRight && !exp.Switch.JoinLeft && !exp.Switch.SplitLeft {
		trackType = tracks.JunctionRight + exp.Switch.Type
	} else if exp.Switch.JoinLeft && !exp.Switch.JoinRight && !exp.Switch.SplitRight {
		trackType = tracks.JunctionLeft + exp.Switch.Type
	} else if exp.Switch.SplitLeft && !exp.Switch.JoinRight && !exp.Switch.SplitRight {
		trackType = tracks.JunctionRight + exp.Switch.Type
	} else if exp.Switch.SplitRight && !exp.Switch.JoinLeft && !exp.Switch.SplitLeft {
		trackType = tracks.JunctionLeft + exp.Switch.Type
	} else if exp.Switch.SplitLeft && exp.Switch.JoinRight {
		trackType = tracks.JunctionCross + exp.Switch.Type
	} else if exp.Switch.JoinLeft && exp.Switch.SplitRight {
		trackType = tracks.JunctionCross + exp.Switch.Type
	} else {
		panic("Oooops")
	}

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

	var indexLeft, indexRight int
	if exp.Switch.JoinLeft {
		indexLeft = getColumnAtPosition(columns, exp.Switch.PositionLeft)
		if indexLeft == -1 {
			b.errlog.LogError(errlog.ErrorColumnStartedUnexpectedly, exp.Location)
			return columns, columns[index].first, columns[index].last
		}
		columns[indexLeft].marked = true
	} else if exp.Switch.SplitLeft {
		indexLeft = getColumnAtPosition(columns, exp.Switch.PositionLeft)
		if indexLeft != -1 {
			b.errlog.LogError(errlog.ErrorColumnTerminatedUnexpectedly, exp.Location)
			return columns, columns[index].first, columns[index].last
		}
		columns, indexLeft = startOrEndColumnAtPosition(columns, errlog.EncodeLocationRange(exp.Location.File(), exp.Location.Line(), exp.Switch.PositionLeft, exp.Location.Line(), exp.Switch.PositionLeft))
	}
	if exp.Switch.JoinRight {
		indexRight = getColumnAtPosition(columns, exp.Switch.PositionRight)
		if indexRight == -1 {
			b.errlog.LogError(errlog.ErrorColumnStartedUnexpectedly, exp.Location)
			return columns, columns[index].first, columns[index].last
		}
		columns[indexRight].marked = true
	} else if exp.Switch.SplitRight {
		indexRight = getColumnAtPosition(columns, exp.Switch.PositionRight)
		if indexRight != -1 {
			b.errlog.LogError(errlog.ErrorColumnTerminatedUnexpectedly, exp.Location)
			return columns, columns[index].first, columns[index].last
		}
		columns, indexRight = startOrEndColumnAtPosition(columns, errlog.EncodeLocationRange(exp.Location.File(), exp.Location.Line(), exp.Switch.PositionRight, exp.Location.Line(), exp.Switch.PositionRight))
	}

	if exp.Switch.JoinRight && !exp.Switch.JoinLeft && !exp.Switch.SplitLeft {
		// `W10------/`
		if newTrack.Geometry.IncomingConnectionCount != 1 || newTrack.Geometry.OutgoingConnectionCount != 2 {
			panic("Track has wrong connection count")
		}
		if columns[index].last != nil {
			columns[index].last.Connect(newTrack.Connection(1))
		}
		if columns[indexRight].last != nil {
			columns[indexRight].last.Connect(newTrack.Connection(2))
		} else {
			columns[indexRight].last = newTrack.Connection(2)
			columns[indexRight].first = newTrack.Connection(2)
		}
		columns = endColumn(columns, indexRight)
		return columns, newTrack.Connection(1), newTrack.Connection(0)
	} else if exp.Switch.JoinLeft && !exp.Switch.JoinRight && !exp.Switch.SplitRight {
		// `\------W10`
		if newTrack.Geometry.IncomingConnectionCount != 1 || newTrack.Geometry.OutgoingConnectionCount != 2 {
			panic("Track has wrong connection count")
		}
		if columns[index].last != nil {
			columns[index].last.Connect(newTrack.Connection(2))
		}
		if columns[indexLeft].last != nil {
			columns[indexLeft].last.Connect(newTrack.Connection(1))
		} else {
			columns[indexLeft].last = newTrack.Connection(1)
			columns[indexLeft].first = newTrack.Connection(1)
		}
		columns = endColumn(columns, indexLeft)
		return columns, newTrack.Connection(2), newTrack.Connection(0)
	} else if exp.Switch.SplitLeft && exp.Switch.JoinRight {
		// `/-------W10------/`
		if newTrack.Geometry.IncomingConnectionCount != 2 || newTrack.Geometry.OutgoingConnectionCount != 2 {
			panic("Track has wrong connection count")
		}
		if columns[index].last != nil {
			columns[index].last.Connect(newTrack.Connection(0))
		}
		if columns[indexRight].last != nil {
			columns[indexRight].last.Connect(newTrack.Connection(1))
		} else {
			columns[indexRight].last = newTrack.Connection(1)
			columns[indexRight].first = newTrack.Connection(1)
		}
		columns[indexLeft].last = newTrack.Connection(3)
		columns[indexLeft].first = newTrack.Connection(3)
		columns = endColumn(columns, indexRight)
		return columns, newTrack.Connection(0), newTrack.Connection(2)
	} else if exp.Switch.JoinLeft && exp.Switch.SplitRight {
		// `\-------W10------\`
		if newTrack.Geometry.IncomingConnectionCount != 2 || newTrack.Geometry.OutgoingConnectionCount != 2 {
			panic("Track has wrong connection count")
		}
		if columns[index].last != nil {
			columns[index].last.Connect(newTrack.Connection(1))
		}
		if columns[indexLeft].last != nil {
			columns[indexLeft].last.Connect(newTrack.Connection(0))
		} else {
			columns[indexLeft].last = newTrack.Connection(0)
			columns[indexLeft].first = newTrack.Connection(0)
		}
		columns[indexRight].last = newTrack.Connection(2)
		columns[indexRight].first = newTrack.Connection(2)
		columns = endColumn(columns, indexLeft)
		return columns, newTrack.Connection(1), newTrack.Connection(3)
	} else if exp.Switch.SplitLeft && !exp.Switch.JoinRight && !exp.Switch.SplitRight {
		// `/-------W10`
		if newTrack.Geometry.IncomingConnectionCount != 1 || newTrack.Geometry.OutgoingConnectionCount != 2 {
			panic("Track has wrong connection count")
		}
		if columns[index].last != nil {
			columns[index].last.Connect(newTrack.Connection(0))
		}
		columns[indexLeft].last = newTrack.Connection(2)
		columns[indexLeft].first = newTrack.Connection(2)
		return columns, newTrack.Connection(0), newTrack.Connection(1)
	} else if exp.Switch.SplitRight && !exp.Switch.JoinLeft && !exp.Switch.SplitLeft {
		// `W10------\`
		if newTrack.Geometry.IncomingConnectionCount != 1 || newTrack.Geometry.OutgoingConnectionCount != 2 {
			panic("Track has wrong connection count")
		}
		if columns[index].last != nil {
			columns[index].last.Connect(newTrack.Connection(0))
		}
		columns[indexRight].last = newTrack.Connection(1)
		columns[indexRight].first = newTrack.Connection(1)
		return columns, newTrack.Connection(0), newTrack.Connection(2)
	}
	panic("Oooops")
}

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

func (b *Interpreter) processTrackTerminationWithConnection(exp *parser.Expression, con *tracks.TrackConnection) {
	if exp.TrackTermination == nil || (!exp.TrackTermination.EllipsisLeft && !exp.TrackTermination.EllipsisRight) {
		panic("Not a track termination with connection")
	}
	if m := b.trackSystem.GetMark(exp.TrackTermination.Name); m != nil && m.Connection != nil {
		if m.Connection.IsConnected() {
			b.errlog.LogError(errlog.ErrorTrackConnectedTwice, exp.Location)
		} else {
			m.Connection.Connect(con)
		}
		return
	}
	if !con.AddMark(exp.TrackTermination.Name) {
		b.errlog.LogError(errlog.ErrorTrackMarkDefinedTwice, exp.Location)
	}
}

func (b *Interpreter) processTrackMark(exp *parser.Expression, track *tracks.Track) {
	if exp.TrackMark == nil {
		panic("Not a track")
	}
	if !track.AddMark(exp.TrackMark.Position, exp.TrackMark.Name) {
		b.errlog.LogError(errlog.ErrorTrackMarkDefinedTwice, exp.Location)
	}
}

func (b *Interpreter) processRepeatExpression(rail *parser.RailWay, exp *parser.Expression, prevCon *tracks.TrackConnection) (first *tracks.TrackConnection, last *tracks.TrackConnection) {
	if exp.Repeat == nil {
		panic("Oooops")
	}
	last = prevCon
	for i := 0; i < exp.Repeat.Count; i++ {
		f, l := b.processExpression(rail, exp.Repeat.TrackExpression, last)
		if first == nil {
			first = f
		}
		last = l
	}
	return
}

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

/********************************************************************
 *
 * Helper functions for arrays of railwayColumn.
 *
 ********************************************************************/

/*
func markOrCreateColumnAtPosition(columns []*railwayColumn, position int) ([]*railwayColumn, int) {
	for i, c := range columns {
		if c.position == position {
			c.marked = true
			return columns, i
		}
	}
	c := &railwayColumn{position: position}
	columns = append(columns, c)
	return columns, len(columns) - 1
}
*/

func startOrEndColumnAtPosition(columns []*railwayColumn, location errlog.LocationRange) ([]*railwayColumn, int) {
	index := getColumnAtPosition(columns, location.Position())
	if index == -1 {
		c := &railwayColumn{location: location, marked: true}
		return append(columns, c), len(columns)
	}
	return endColumn(columns, index), index
}

func endColumn(columns []*railwayColumn, index int) []*railwayColumn {
	return append(columns[:index], columns[index+1:]...)
}

func getColumnAtPosition(columns []*railwayColumn, position int) int {
	for i, c := range columns {
		if c.location.Position() == position {
			return i
		}
	}
	return -1
}
