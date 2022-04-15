package tracks

// Tracks of the same kind all share the same track geometry.
type TrackGeometry struct {
	Name string
	// The paths describe the lines and arcs that make up a track.
	Paths []ITrackGeometryPath
	// Position of a connection point relative to the center of the track
	ConnectionPoints []TrackGeometryPoint
	// Lists all possible ways to drive from an incoming connection point to an outgoing connection point.
	// A normal track has exactly two connection points and therefore one TurnoutOption,
	// i.e. the option of driving along the track.
	// Of course tracks can be used in two directions, but this is ignore here
	// since we only talk about geometry and orientation.
	TurnoutOptions []TurnoutOption
	// A connection point is incoming, if a drain driving along the orientation of the track
	// enters the track via the connection point.
	// The ConnectionPoints array holds first the incoming and then the outgoing connection points
	// sorted clock-wise.
	IncomingConnectionCount int
	OutgoingConnectionCount int
}

// A turnout can allow the train to drive
// from one of its connection points to another.
type TurnoutOption struct {
	// Index of a connection point.
	From int
	// Index of a connection point.
	To int
}

// A point relative to the center and orientation of a track.
type TrackGeometryPoint struct {
	// Track entry angle in degrees relative to the
	// orientation of the track
	Angle float64
	// Position vector in mm pointing from the connection point to
	// the center of the track.
	Position Vec2
}

// A path describes a piece of track that can be graphically represented.
type ITrackGeometryPath interface {
}

// Implements ITrackGeometryPath.
type TrackGeometryLine struct {
	// Size in mm
	Size   float64
	Anchor TrackGeometryPoint
}

// Implements ITrackGeometryPath.
type TrackGeometryArc struct {
	// Angle in degree
	TrackAngle float64
	// Radius in mm.
	Radius float64
	Anchor TrackGeometryPoint
}

// Describes the location and orientation of a track.
type TrackLocation struct {
	// Vector pointing from the origin to the center of the track.
	Center Vec3
	// Angle in degree. Zero means the track is aligned with the x-axis
	// and is heading to the right side.
	Rotation float64
	// Incline in percent of the track length
	Incline float64
}

func NewTrackLocation(con *TrackConnection, pos Vec3, angle float64) *TrackLocation {
	p := &con.Track.Geometry.ConnectionPoints[con.Track.ConnectionIndex(con)]
	r := angle - p.Angle
	c := p.Position.Rotate(r)
	center := [3]float64{pos[0] + c[0], pos[1] + c[1], pos[2]}
	rotation := normalizeAngle(r)
	return &TrackLocation{Center: center, Rotation: rotation}
}

// Given a track location and a track geometry, the function returns the
// position and angle of a connection point of the track.
func (l *TrackLocation) Connection(index int, geometry *TrackGeometry) (pos Vec3, angle float64) {
	c := &geometry.ConnectionPoints[index]
	angle = normalizeAngle(l.Rotation + c.Angle + 180)
	pos = l.Center.Add2(c.Position.Invert().Rotate(l.Rotation))
	return
}
