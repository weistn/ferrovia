package tracks

import (
	"github.com/weistn/ferrovia/errlog"
)

// A TrackSystem connsist of tracks, which are usually connected with each other (but not necessarily).
type TrackSystem struct {
	Layers map[string]*TrackLayer
	marks  map[string]*TrackMark
}

type TrackLayer struct {
	TrackSystem *TrackSystem
	Name        string
	Color       string
	Tracks      []*Track
}

type Track struct {
	// Immutable
	Layer *TrackLayer
	// Immutable
	Id int
	// Immutable
	Geometry *TrackGeometry
	// Read-only
	Location       *TrackLocation
	connections    []*TrackConnection
	connectReverse bool
	tag            int
	// Read-Only
	Marks []*TrackMark
	// The index of the currently selected option.
	SelectedTurnoutOption int
	SourceLocation        errlog.LocationRange
}

// Each track has multiple connection points, each represented by TrackConnection.
// Physically joining two tracks means joining two TrackConnectins.
type TrackConnection struct {
	Track *Track
	// If null, the track ends here.
	Opposite *TrackConnection
	Marks    []*TrackMark
}

// A TrackMark marks a position on a Track, usually the beginning or end of a track.
// Tracks marks are only allowed on straight tracks, i.e. not on turnout tracks.
// TrackMarks can be used to position signaling or to denote blocks.
type TrackMark struct {
	// A value in the range [0,1].
	// A value of 0 means the positon of the first connection point.
	// A value of 1 means the Position of the second connection point.
	// Values inbetween indicate a proportional Position on the track.
	position float32
	// Optional unique name of the TrackMark or the empty string
	name  string
	track *Track
	// May be nil.
	Connection *TrackConnection
}

type TrackFactoryFunc func(l *TrackLayer, id int) *Track

const (
	JunctionLeft   = "L-"
	JunctionRight  = "R-"
	JunctionDouble = "D-"
	JunctionCross  = "K-"
)

// A map of all registered track factories.
var TrackFactories map[string]TrackFactoryFunc = make(map[string]TrackFactoryFunc)

// Global variable used by NewEpoch.
var epoch int

// Global variable used by NewTrack.
var trackId int

// Registers a function that can create a track of a certain kind
func RegisterTrackFactory(kind string, fn TrackFactoryFunc) {
	TrackFactories[kind] = fn
}

func NewTrackSystem() *TrackSystem {
	ts := &TrackSystem{marks: make(map[string]*TrackMark), Layers: make(map[string]*TrackLayer)}
	ts.AddLayer(&TrackLayer{Name: ""})
	return ts
}

func (ts *TrackSystem) AddLayer(l *TrackLayer) {
	if _, ok := ts.Layers[l.Name]; ok {
		panic("Duplicate layer name")
	}
	ts.Layers[l.Name] = l
}

func (ts *TrackSystem) GetMark(name string) *TrackMark {
	m, ok := ts.marks[name]
	if ok {
		return m
	}
	return nil
}

// Creates a track of the given kind.
// Returns nil if no corresponding factory has been registered.
func (l *TrackLayer) NewTrack(kind string) *Track {
	fn, ok := TrackFactories[kind]
	if !ok {
		return nil
	}
	return fn(l, 0)
}

func NewEpoch() int {
	epoch++
	return epoch
}

func (m *TrackMark) Track() *Track {
	return m.track
}

func (m *TrackMark) Name() string {
	return m.name
}

func (m *TrackMark) Position() float32 {
	return m.position
}

func (c *TrackConnection) Connect(c2 *TrackConnection) {
	// println("CONNECT", c.Track.Geometry.Name, c.Track.ConnectionIndex(c), c2.Track.Geometry.Name, c2.Track.ConnectionIndex(c2))
	if c.Opposite != nil && c.Opposite != c2 {
		panic("Already connected")
	}
	if c2.Opposite != nil && c2.Opposite != c {
		panic("Already connected")
	}
	c.Opposite = c2
	c2.Opposite = c
}

func (c *TrackConnection) IsConnected() bool {
	return c.Opposite != nil
}

// Returns false if a mark of the same name already exists
func (c *TrackConnection) AddMark(name string) bool {
	if _, ok := c.Track.Layer.TrackSystem.marks[name]; ok {
		return false
	}
	mark := &TrackMark{name: name, Connection: c, track: c.Track}
	c.Marks = append(c.Marks, mark)
	if name != "" {
		c.Track.Layer.TrackSystem.marks[name] = mark
	}
	return true
}

// Should only be used by track factories.
// Otherwise use the TrackSystem to create new tracks.
func NewTrack(layer *TrackLayer, id int, geometry *TrackGeometry, reverse bool) *Track {
	if id == 0 {
		trackId++
		id = trackId
	}
	t := &Track{Layer: layer, Id: id, Geometry: geometry, connections: make([]*TrackConnection, len(geometry.ConnectionPoints)), SelectedTurnoutOption: 0, connectReverse: reverse}
	// Initialize the connections of the track
	for i := range t.connections {
		t.connections[i] = &TrackConnection{Track: t}
	}
	// Append the track to a system of tracks
	layer.Tracks = append(layer.Tracks, t)
	return t
}

// Returns false if a mark of the same name already exists
func (t *Track) AddMark(position float32, name string) bool {
	if _, ok := t.Layer.TrackSystem.marks[name]; ok {
		return false
	}
	mark := &TrackMark{position: position, name: name, track: t}
	t.Marks = append(t.Marks, mark)
	if name != "" {
		t.Layer.TrackSystem.marks[name] = mark
	}
	return true
}

// Returns false if a location has already been set
func (t *Track) SetLocation(l *TrackLocation) bool {
	if t.Location != nil {
		return false
	}
	t.Location = l
	return true
}

// At each point in time, a track reaches from one connection to another.
// Turnouts can change the reachable connection by switching.
// Thus, a track can have 2, 3, 4 or more connections, but only two are reachable at any time.
func (t *Track) ActiveConnections() (*TrackConnection, *TrackConnection) {
	return t.connections[t.Geometry.TurnoutOptions[t.SelectedTurnoutOption].From], t.connections[t.Geometry.TurnoutOptions[t.SelectedTurnoutOption].To]
}

// Each track has one selected TrackOption, that describes that trains can drive from
// one connection to another. FirstConnection returns the first such connection.
// First and second connection are reversed if the track has been reversed.
func (t *Track) FirstConnection() *TrackConnection {
	option := t.Geometry.TurnoutOptions[t.SelectedTurnoutOption]
	if t.connectReverse {
		return t.connections[option.To]
	}
	return t.connections[option.From]
}

// Each track has one selected TrackOption, that describes that trains can drive from
// one connection to another. FirstConnection returns the first such connection.
// First and second connection are reversed if the track has been reversed.
func (t *Track) SecondConnection() *TrackConnection {
	option := t.Geometry.TurnoutOptions[t.SelectedTurnoutOption]
	if t.connectReverse {
		return t.connections[option.From]
	}
	return t.connections[option.To]
}

// Used during track construction to retrieve the next track connection that
// is located on a turnout not currently selected.
// Do not use the function directly. Use Connect instead.
func (t *Track) TurnoutConnection(number int) *TrackConnection {
	option := t.Geometry.TurnoutOptions[t.SelectedTurnoutOption]
	count := 0
	for i := 0; i < len(t.connections); i++ {
		if i == option.From || i == option.To {
			continue
		}
		if count == number {
			return t.connections[i]
		}
		count++
	}
	return nil
}

// TODO: Remove?
// Used during track construction.
func (t *Track) Connect(track *Track) *Track {
	con1 := t.SecondConnection()
	con2 := track.FirstConnection()
	if con1.IsConnected() || con2.IsConnected() {
		panic("tracks are already connected")
	}
	con1.Connect(con2)
	return track
}

// Used when algorithms traverse the tracks to detect loops.
// To avoid confusion, call NewEpoch() before tagging as this
// invalidates all tags made in the previous epoch.
func (t *Track) Tag() {
	t.tag = epoch
}

// Used when algorithms traverse the tracks to detect loops.
// Returns true if the tracks has been tagged in the current epoch.
func (t *Track) IsTagged() bool {
	return t.tag == epoch
}

func (t *Track) ConnectionIndex(c *TrackConnection) int {
	for i, c2 := range t.connections {
		if c2 == c {
			return i
		}
	}
	panic("Connection not part of this track")
}

func (t *Track) Connection(index int) *TrackConnection {
	if index < 0 || index >= len(t.connections) {
		panic("Oooops")
	}
	return t.connections[index]
}

func (t *Track) ConnectionCount() int {
	return len(t.connections)
}
