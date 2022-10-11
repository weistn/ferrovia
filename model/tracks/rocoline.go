package tracks

import (
	"math"
)

var rocoG025 *TrackGeometry
var rocoG05 *TrackGeometry
var rocoG1 *TrackGeometry
var rocoG4 *TrackGeometry
var rocoDG1 *TrackGeometry
var rocoR5 *TrackGeometry
var rocoR6 *TrackGeometry
var rocoR9 *TrackGeometry
var rocoR10 *TrackGeometry
var rocoW15R *TrackGeometry
var rocoW15L *TrackGeometry
var rocoBWR5 *TrackGeometry
var rocoBWL5 *TrackGeometry
var rocoBWR9 *TrackGeometry
var rocoBWL9 *TrackGeometry
var rocoDKW15 *TrackGeometry
var rocoK15 *TrackGeometry
var rocoW10R *TrackGeometry
var rocoW10L *TrackGeometry
var rocoDKW10 *TrackGeometry
var rocoD2 *TrackGeometry
var rocoD8 *TrackGeometry

func registerCustomCurve(name string, radius float64, angle float64) {
	geo := &TrackGeometry{
		Name: name,
		Paths: []ITrackGeometryPath{
			&TrackGeometryArc{TrackAngle: angle, Radius: radius, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
		},
		ConnectionPoints: []TrackGeometryPoint{
			{
				Position: [2]float64{0, 0},
				Angle:    0,
			},
			{
				Position: [2]float64{-radius * (1. - math.Cos(math.Pi/180*angle)), radius * math.Sin(math.Pi/180*angle)},
				Angle:    180 + angle,
			},
		},
		IncomingConnectionCount: 1,
		OutgoingConnectionCount: 1,
		TurnoutOptions:          []TurnoutOption{{From: 0, To: 1}},
	}

	fright := func(l *TrackLayer, id int) *Track {
		InitRoco()
		t := NewTrack(l, id, geo, false)
		return t
	}
	fleft := func(l *TrackLayer, id int) *Track {
		InitRoco()
		t := NewTrack(l, id, geo, true)
		return t
	}

	RegisterTrackFactory("R"+name, fright)
	RegisterTrackFactory("L"+name, fleft)
}

func InitRoco() {
	if rocoG1 != nil {
		return
	}
	rocoD2 = &TrackGeometry{
		Name: "D2",
		Paths: []ITrackGeometryPath{
			&TrackGeometryLine{Size: 2.5, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
		},
		ConnectionPoints: []TrackGeometryPoint{
			{
				Position: [2]float64{0, 0},
				Angle:    0,
			},
			{
				Position: [2]float64{0, 2.5},
				Angle:    180,
			},
		},
		IncomingConnectionCount: 1,
		OutgoingConnectionCount: 1,
		TurnoutOptions:          []TurnoutOption{{From: 0, To: 1}},
	}
	rocoD8 = &TrackGeometry{
		Name: "D8",
		Paths: []ITrackGeometryPath{
			&TrackGeometryLine{Size: 8, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
		},
		ConnectionPoints: []TrackGeometryPoint{
			{
				Position: [2]float64{0, 0},
				Angle:    0,
			},
			{
				Position: [2]float64{0, 8},
				Angle:    180,
			},
		},
		IncomingConnectionCount: 1,
		OutgoingConnectionCount: 1,
		TurnoutOptions:          []TurnoutOption{{From: 0, To: 1}},
	}
	rocoG025 = &TrackGeometry{
		Name: "G025",
		Paths: []ITrackGeometryPath{
			&TrackGeometryLine{Size: 57.5, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
		},
		ConnectionPoints: []TrackGeometryPoint{
			{
				Position: [2]float64{0, 0},
				Angle:    0,
			},
			{
				Position: [2]float64{0, 57.5},
				Angle:    180,
			},
		},
		IncomingConnectionCount: 1,
		OutgoingConnectionCount: 1,
		TurnoutOptions:          []TurnoutOption{{From: 0, To: 1}},
	}
	rocoG05 = &TrackGeometry{
		Name: "G05",
		Paths: []ITrackGeometryPath{
			&TrackGeometryLine{Size: 115, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
		},
		ConnectionPoints: []TrackGeometryPoint{
			{
				Position: [2]float64{0, 0},
				Angle:    0,
			},
			{
				Position: [2]float64{0, 115},
				Angle:    180,
			},
		},
		IncomingConnectionCount: 1,
		OutgoingConnectionCount: 1,
		TurnoutOptions:          []TurnoutOption{{From: 0, To: 1}},
	}
	rocoG1 = &TrackGeometry{
		Name: "G1",
		Paths: []ITrackGeometryPath{
			&TrackGeometryLine{Size: 230, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
		},
		ConnectionPoints: []TrackGeometryPoint{
			{
				Position: [2]float64{0, 0},
				Angle:    0,
			},
			{
				Position: [2]float64{0, 230},
				Angle:    180,
			},
		},
		IncomingConnectionCount: 1,
		OutgoingConnectionCount: 1,
		TurnoutOptions:          []TurnoutOption{{From: 0, To: 1}},
	}
	rocoG4 = &TrackGeometry{
		Name: "G4",
		Paths: []ITrackGeometryPath{
			&TrackGeometryLine{Size: 230 * 4, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
		},
		ConnectionPoints: []TrackGeometryPoint{
			{
				Position: [2]float64{0, 0},
				Angle:    0,
			},
			{
				Position: [2]float64{0, 230 * 4},
				Angle:    180,
			},
		},
		IncomingConnectionCount: 1,
		OutgoingConnectionCount: 1,
		TurnoutOptions:          []TurnoutOption{{From: 0, To: 1}},
	}
	rocoDG1 = &TrackGeometry{
		Name: "DG1",
		Paths: []ITrackGeometryPath{
			&TrackGeometryLine{Size: 119, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
		},
		ConnectionPoints: []TrackGeometryPoint{
			{
				Position: [2]float64{0, 0},
				Angle:    0,
			},
			{
				Position: [2]float64{0, 119},
				Angle:    180,
			},
		},
		IncomingConnectionCount: 1,
		OutgoingConnectionCount: 1,
		TurnoutOptions:          []TurnoutOption{{From: 0, To: 1}},
	}
	// Track turning right
	rocoR5 = &TrackGeometry{
		Name: "R5",
		Paths: []ITrackGeometryPath{
			&TrackGeometryArc{TrackAngle: 30, Radius: 542.8, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
		},
		ConnectionPoints: []TrackGeometryPoint{
			{
				Position: [2]float64{0, 0},
				Angle:    0,
			},
			{
				Position: [2]float64{-542.8 * (1. - math.Cos(math.Pi/6.0)), 542.8 * math.Sin(math.Pi/6.0)},
				Angle:    180 + 30,
			},
		},
		IncomingConnectionCount: 1,
		OutgoingConnectionCount: 1,
		TurnoutOptions:          []TurnoutOption{{From: 0, To: 1}},
	}
	// Track turning right
	rocoR6 = &TrackGeometry{
		Name: "R6",
		Paths: []ITrackGeometryPath{
			&TrackGeometryArc{TrackAngle: 30, Radius: 604.4, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
		},
		ConnectionPoints: []TrackGeometryPoint{
			{
				Position: [2]float64{0, 0},
				Angle:    0,
			},
			{
				Position: [2]float64{-604.4 * (1. - math.Cos(math.Pi/6.0)), 604.4 * math.Sin(math.Pi/6.0)},
				Angle:    180 + 30,
			},
		},
		IncomingConnectionCount: 1,
		OutgoingConnectionCount: 1,
		TurnoutOptions:          []TurnoutOption{{From: 0, To: 1}},
	}
	// Track turning right
	rocoR9 = &TrackGeometry{
		Name: "R9",
		Paths: []ITrackGeometryPath{
			&TrackGeometryArc{TrackAngle: 15, Radius: 826.4, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
		},
		ConnectionPoints: []TrackGeometryPoint{
			{
				Position: [2]float64{0, 0},
				Angle:    0,
			},
			{
				Position: [2]float64{-826.4 * (1. - math.Cos(math.Pi/12.0)), 826.4 * math.Sin(math.Pi/12.0)},
				Angle:    180 + 15,
			},
		},
		IncomingConnectionCount: 1,
		OutgoingConnectionCount: 1,
		TurnoutOptions:          []TurnoutOption{{From: 0, To: 1}},
	}
	// Track turning right
	rocoR10 = &TrackGeometry{
		Name: "R10",
		Paths: []ITrackGeometryPath{
			&TrackGeometryArc{TrackAngle: 15, Radius: 888, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
		},
		ConnectionPoints: []TrackGeometryPoint{
			{
				Position: [2]float64{0, 0},
				Angle:    0,
			},
			{
				Position: [2]float64{-888 * (1. - math.Cos(math.Pi/12.0)), 888 * math.Sin(math.Pi/12.0)},
				Angle:    180 + 15,
			},
		},
		IncomingConnectionCount: 1,
		OutgoingConnectionCount: 1,
		TurnoutOptions:          []TurnoutOption{{From: 0, To: 1}},
	}
	// Track turning right
	rocoW15R = &TrackGeometry{
		Name: "WR15",
		Paths: []ITrackGeometryPath{
			&TrackGeometryLine{Size: 230, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
			&TrackGeometryArc{TrackAngle: 15, Radius: 888, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
		},
		ConnectionPoints: []TrackGeometryPoint{
			{
				Position: [2]float64{0, 0},
				Angle:    0,
			},
			{
				Position: [2]float64{0, 230},
				Angle:    180,
			},
			{
				Position: [2]float64{-888 * (1. - math.Cos(math.Pi/12.0)), 888 * math.Sin(math.Pi/12.0)},
				Angle:    180 + 15,
			},
		},
		IncomingConnectionCount: 1,
		OutgoingConnectionCount: 2,
		TurnoutOptions:          []TurnoutOption{{From: 0, To: 1}, {From: 0, To: 2}},
	}
	// Track turning right
	rocoW15L = &TrackGeometry{
		Name: "WL15",
		Paths: []ITrackGeometryPath{
			&TrackGeometryLine{Size: 230, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
			&TrackGeometryArc{TrackAngle: 15, Radius: 888, Anchor: TrackGeometryPoint{
				Position: [2]float64{-888 * (1. - math.Cos(math.Pi/12.0)), -888 * math.Sin(math.Pi/12.0)},
				Angle:    180 - 15,
			}},
		},
		ConnectionPoints: []TrackGeometryPoint{
			{
				Position: [2]float64{0, 0},
				Angle:    0,
			},
			{
				Position: [2]float64{888 * (1. - math.Cos(math.Pi/12.0)), 888 * math.Sin(math.Pi/12.0)},
				Angle:    180 - 15,
			},
			{
				Position: [2]float64{0, 230},
				Angle:    180,
			},
		},
		IncomingConnectionCount: 1,
		OutgoingConnectionCount: 2,
		TurnoutOptions:          []TurnoutOption{{From: 0, To: 2}, {From: 0, To: 1}},
	}
	rocoBWR5 = &TrackGeometry{
		Name: "BWR5",
		Paths: []ITrackGeometryPath{
			&TrackGeometryArc{TrackAngle: 30, Radius: 542.8, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
			&TrackGeometryLine{Size: 61, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
			&TrackGeometryArc{TrackAngle: 30, Radius: 542.8, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, -61},
				Angle:    0,
			}},
		},
		ConnectionPoints: []TrackGeometryPoint{
			{
				Position: [2]float64{0, 0},
				Angle:    0,
			},
			{
				Position: [2]float64{-542.8 * (1. - math.Cos(math.Pi/6.0)), 542.8*math.Sin(math.Pi/6.0) + 61},
				Angle:    180 + 30,
			},
			{
				Position: [2]float64{-542.8 * (1. - math.Cos(math.Pi/6.0)), 542.8 * math.Sin(math.Pi/6.0)},
				Angle:    180 + 30,
			},
		},
		IncomingConnectionCount: 1,
		OutgoingConnectionCount: 2,
		TurnoutOptions:          []TurnoutOption{{From: 0, To: 2}, {From: 0, To: 1}},
	}
	rocoBWL5 = &TrackGeometry{
		Name: "BWL5",
		Paths: []ITrackGeometryPath{
			&TrackGeometryArc{TrackAngle: 30, Radius: 542.8, Anchor: TrackGeometryPoint{
				Position: [2]float64{-542.8 * (1. - math.Cos(math.Pi/6.0)), -542.8*math.Sin(math.Pi/6.0) - 61},
				Angle:    180 - 30,
			}},
			&TrackGeometryLine{Size: 61, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
			&TrackGeometryArc{TrackAngle: 30, Radius: 542.8, Anchor: TrackGeometryPoint{
				Position: [2]float64{-542.8 * (1. - math.Cos(math.Pi/6.0)), -542.8 * math.Sin(math.Pi/6.0)},
				Angle:    180 - 30,
			}},
		},
		ConnectionPoints: []TrackGeometryPoint{
			{
				Position: [2]float64{0, 0},
				Angle:    0,
			},
			{
				Position: [2]float64{542.8 * (1. - math.Cos(math.Pi/6.0)), 542.8 * math.Sin(math.Pi/6.0)},
				Angle:    180 - 30,
			},
			{
				Position: [2]float64{542.8 * (1. - math.Cos(math.Pi/6.0)), 542.8*math.Sin(math.Pi/6.0) + 61},
				Angle:    180 - 30,
			},
		},
		IncomingConnectionCount: 1,
		OutgoingConnectionCount: 2,
		TurnoutOptions:          []TurnoutOption{{From: 0, To: 2}, {From: 0, To: 1}},
	}
	rocoBWR9 = &TrackGeometry{
		Name: "BWR9",
		Paths: []ITrackGeometryPath{
			&TrackGeometryArc{TrackAngle: 30, Radius: 826.4, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
			&TrackGeometryLine{Size: 61, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
			&TrackGeometryArc{TrackAngle: 30, Radius: 826.4, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, -61},
				Angle:    0,
			}},
		},
		ConnectionPoints: []TrackGeometryPoint{
			{
				Position: [2]float64{0, 0},
				Angle:    0,
			},
			{
				Position: [2]float64{-826.4 * (1. - math.Cos(math.Pi/6.0)), 826.4*math.Sin(math.Pi/6.0) + 61},
				Angle:    180 + 30,
			},
			{
				Position: [2]float64{-826.4 * (1. - math.Cos(math.Pi/6.0)), 826.4 * math.Sin(math.Pi/6.0)},
				Angle:    180 + 30,
			},
		},
		IncomingConnectionCount: 1,
		OutgoingConnectionCount: 2,
		TurnoutOptions:          []TurnoutOption{{From: 0, To: 2}, {From: 0, To: 1}},
	}
	rocoBWL9 = &TrackGeometry{
		Name: "BWL9",
		Paths: []ITrackGeometryPath{
			&TrackGeometryArc{TrackAngle: 30, Radius: 826.4, Anchor: TrackGeometryPoint{
				Position: [2]float64{-826.4 * (1. - math.Cos(math.Pi/6.0)), -826.4*math.Sin(math.Pi/6.0) - 61},
				Angle:    180 - 30,
			}},
			&TrackGeometryLine{Size: 61, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
			&TrackGeometryArc{TrackAngle: 30, Radius: 826.4, Anchor: TrackGeometryPoint{
				Position: [2]float64{-826.4 * (1. - math.Cos(math.Pi/6.0)), -826.4 * math.Sin(math.Pi/6.0)},
				Angle:    180 - 30,
			}},
		},
		ConnectionPoints: []TrackGeometryPoint{
			{
				Position: [2]float64{0, 0},
				Angle:    0,
			},
			{
				Position: [2]float64{826.4 * (1. - math.Cos(math.Pi/6.0)), 826.4 * math.Sin(math.Pi/6.0)},
				Angle:    180 - 30,
			},
			{
				Position: [2]float64{826.4 * (1. - math.Cos(math.Pi/6.0)), 826.4*math.Sin(math.Pi/6.0) + 61},
				Angle:    180 - 30,
			},
		},
		IncomingConnectionCount: 1,
		OutgoingConnectionCount: 2,
		TurnoutOptions:          []TurnoutOption{{From: 0, To: 2}, {From: 0, To: 1}},
	}
	rocoDKW15 = &TrackGeometry{
		Name: "DKW15",
		Paths: []ITrackGeometryPath{
			&TrackGeometryLine{Size: 230, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
			&TrackGeometryLine{Size: 230, Anchor: TrackGeometryPoint{
				Position: [2]float64{115 * math.Sin(math.Pi/12), -115 * (1 - math.Cos(math.Pi/12))},
				Angle:    345,
			}},
		},
		ConnectionPoints: []TrackGeometryPoint{
			{
				Position: [2]float64{-115 * math.Sin(math.Pi/12), 115 * (1 - math.Cos(math.Pi/12))},
				Angle:    345,
			},
			{
				Position: [2]float64{0, 0},
				Angle:    0,
			},
			{
				Position: [2]float64{115 * math.Sin(math.Pi/12), 230 - 115*(1-math.Cos(math.Pi/12))},
				Angle:    165,
			},
			{
				Position: [2]float64{0, 230},
				Angle:    180,
			},
		},
		IncomingConnectionCount: 2,
		OutgoingConnectionCount: 2,
		TurnoutOptions:          []TurnoutOption{{From: 0, To: 2}, {From: 0, To: 3}, {From: 1, To: 2}, {From: 1, To: 3}},
	}
	rocoK15 = &TrackGeometry{
		Name: "K15",
		Paths: []ITrackGeometryPath{
			&TrackGeometryLine{Size: 230, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
			&TrackGeometryLine{Size: 230, Anchor: TrackGeometryPoint{
				Position: [2]float64{115 * math.Sin(math.Pi/12), -115 * (1 - math.Cos(math.Pi/12))},
				Angle:    345,
			}},
		},
		ConnectionPoints: []TrackGeometryPoint{
			{
				Position: [2]float64{-115 * math.Sin(math.Pi/12), 115 * (1 - math.Cos(math.Pi/12))},
				Angle:    345,
			},
			{
				Position: [2]float64{0, 0},
				Angle:    0,
			},
			{
				Position: [2]float64{115 * math.Sin(math.Pi/12), 230 - 115*(1-math.Cos(math.Pi/12))},
				Angle:    165,
			},
			{
				Position: [2]float64{0, 230},
				Angle:    180,
			},
		},
		IncomingConnectionCount: 2,
		OutgoingConnectionCount: 2,
		TurnoutOptions:          []TurnoutOption{{From: 0, To: 2}, {From: 1, To: 3}},
	}
	// Track turning right
	rocoW10R = &TrackGeometry{
		Name: "WR10",
		Paths: []ITrackGeometryPath{
			&TrackGeometryLine{Size: 345, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
			&TrackGeometryArc{TrackAngle: 10, Radius: 1946, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
		},
		ConnectionPoints: []TrackGeometryPoint{
			{
				Position: [2]float64{0, 0},
				Angle:    0,
			},
			{
				Position: [2]float64{0, 345},
				Angle:    180,
			},
			{
				Position: [2]float64{-1946 * (1. - math.Cos(math.Pi/18.0)), 1946 * math.Sin(math.Pi/18.0)},
				Angle:    180 + 10,
			},
		},
		IncomingConnectionCount: 1,
		OutgoingConnectionCount: 2,
		TurnoutOptions:          []TurnoutOption{{From: 0, To: 1}, {From: 0, To: 2}},
	}
	// Track turning right
	rocoW10L = &TrackGeometry{
		Name: "WL10",
		Paths: []ITrackGeometryPath{
			&TrackGeometryLine{Size: 345, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
			&TrackGeometryArc{TrackAngle: 10, Radius: 1946, Anchor: TrackGeometryPoint{
				Position: [2]float64{-1946 * (1. - math.Cos(math.Pi/18.0)), -1946 * math.Sin(math.Pi/18.0)},
				Angle:    180 - 10,
			}},
		},
		ConnectionPoints: []TrackGeometryPoint{
			{
				Position: [2]float64{0, 0},
				Angle:    0,
			},
			{
				Position: [2]float64{1946 * (1. - math.Cos(math.Pi/18.0)), 1946 * math.Sin(math.Pi/18.0)},
				Angle:    180 - 10,
			},
			{
				Position: [2]float64{0, 345},
				Angle:    180,
			},
		},
		IncomingConnectionCount: 1,
		OutgoingConnectionCount: 2,
		TurnoutOptions:          []TurnoutOption{{From: 0, To: 1}, {From: 0, To: 2}},
	}
	rocoDKW10 = &TrackGeometry{
		Name: "DKW10",
		Paths: []ITrackGeometryPath{
			&TrackGeometryLine{Size: 345, Anchor: TrackGeometryPoint{
				Position: [2]float64{0, 0},
				Angle:    0,
			}},
			&TrackGeometryLine{Size: 345, Anchor: TrackGeometryPoint{
				Position: [2]float64{172.5 * math.Sin(math.Pi/18), -172.5 * (1 - math.Cos(math.Pi/18))},
				Angle:    350,
			}},
		},
		ConnectionPoints: []TrackGeometryPoint{
			{
				Position: [2]float64{-172.5 * math.Sin(math.Pi/18), 172.5 * (1 - math.Cos(math.Pi/18))},
				Angle:    350,
			},
			{
				Position: [2]float64{0, 0},
				Angle:    0,
			},
			{
				Position: [2]float64{172.5 * math.Sin(math.Pi/18), 345 - 172.5*(1-math.Cos(math.Pi/18))},
				Angle:    170,
			},
			{
				Position: [2]float64{0, 345},
				Angle:    180,
			},
		},
		IncomingConnectionCount: 2,
		OutgoingConnectionCount: 2,
		TurnoutOptions:          []TurnoutOption{{From: 0, To: 2}, {From: 0, To: 3}, {From: 1, To: 2}, {From: 1, To: 3}},
	}

	RegisterTrackFactory("R5", NewR5Right)
	RegisterTrackFactory("L5", NewR5Left)
	RegisterTrackFactory("R6", NewR6Right)
	RegisterTrackFactory("L6", NewR6Left)
	RegisterTrackFactory("R9", NewR9Right)
	RegisterTrackFactory("L9", NewR9Left)
	RegisterTrackFactory("R10", NewR10Right)
	RegisterTrackFactory("L10", NewR10Left)
	RegisterTrackFactory("WR15", NewW15Right)
	RegisterTrackFactory("WL15", NewW15Left)
	RegisterTrackFactory("BWR5", NewBWR5)
	RegisterTrackFactory("BWL5", NewBWL5)
	RegisterTrackFactory("BWR9", NewBWR9)
	RegisterTrackFactory("BWL9", NewBWL9)
	RegisterTrackFactory("DKW15", NewDKW15)
	RegisterTrackFactory("K15", NewK15)
	RegisterTrackFactory("G025", NewG025)
	RegisterTrackFactory("G05", NewG05)
	RegisterTrackFactory("G1", NewG1)
	RegisterTrackFactory("G4", NewG4)
	RegisterTrackFactory("DG1", NewDG1)
	RegisterTrackFactory("WR10", NewW10Right)
	RegisterTrackFactory("WL10", NewW10Left)
	RegisterTrackFactory("DKW10", NewDKW10)
	RegisterTrackFactory("D2", NewD2)
	RegisterTrackFactory("D8", NewD8)
	registerCustomCurve("20", 1962, 5)

	registerCustomCurve("C5", 542.8, 5)
	registerCustomCurve("C6", 542.8+61.6, 5)
	registerCustomCurve("C7", 542.8+2*61.6, 5)
	registerCustomCurve("C8", 888-2*61.6, 5)
	registerCustomCurve("C9", 888-61.6, 5)
	registerCustomCurve("C10", 888, 5)
	registerCustomCurve("C11", 888+61.6, 5)
	registerCustomCurve("C12", 888+2*61.6, 5)
	registerCustomCurve("C13", 888+3*61.6, 5)
	registerCustomCurve("C14", 888+4*61.6, 5)
	registerCustomCurve("C15", 888+5*61.6, 5)
	registerCustomCurve("C16", 888+6*61.6, 5)
	registerCustomCurve("C17", 888+7*61.6, 5)
	registerCustomCurve("C18", 888+8*61.6, 5)
	registerCustomCurve("C19", 888+9*61.6, 5)
	registerCustomCurve("C20", 888+10*61.6, 5)
	registerCustomCurve("C21", 888+11*61.6, 5)
	registerCustomCurve("C22", 888+12*61.6, 5)
	registerCustomCurve("C23", 888+13*61.6, 5)
	registerCustomCurve("C24", 888+14*61.6, 5)
	registerCustomCurve("C25", 888+15*61.6, 5)
	registerCustomCurve("C26", 888+16*61.6, 5)
	registerCustomCurve("C27", 888+17*61.6, 5)
	registerCustomCurve("C28", 888+18*61.6, 5)
	registerCustomCurve("C29", 888+19*61.6, 5)
	registerCustomCurve("C30", 888+20*61.6, 5)
	registerCustomCurve("C31", 888+21*61.6, 5)
	registerCustomCurve("C32", 888+22*61.6, 5)
	registerCustomCurve("C33", 888+23*61.6, 5)
	registerCustomCurve("C34", 888+24*61.6, 5)
	registerCustomCurve("C35", 888+25*61.6, 5)
	registerCustomCurve("C36", 888+26*61.6, 5)
	registerCustomCurve("C37", 888+27*61.6, 5)
	registerCustomCurve("C38", 888+28*61.6, 5)
	registerCustomCurve("C39", 888+29*61.6, 5)
	registerCustomCurve("C40", 888+30*61.6, 5)
}

func NewR5Left(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoR5, true)
	return t
}

func NewR5Right(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoR5, false)
	return t
}

func NewR6Left(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoR6, true)
	return t
}

func NewR6Right(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoR6, false)
	return t
}

func NewR9Left(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoR9, true)
	return t
}

func NewR9Right(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoR9, false)
	return t
}

func NewR10Left(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoR10, true)
	return t
}

func NewR10Right(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoR10, false)
	return t
}

func NewG025(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoG025, false)
	return t
}

func NewG05(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoG05, false)
	return t
}

func NewG1(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoG1, false)
	return t
}

func NewG4(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoG4, false)
	return t
}

func NewDG1(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoDG1, false)
	return t
}

func NewW15Right(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoW15R, false)
	return t
}

func NewW15Left(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoW15L, false)
	return t
}

func NewBWR5(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoBWR5, false)
	return t
}

func NewBWL5(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoBWL5, false)
	return t
}

func NewBWR9(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoBWR9, false)
	return t
}

func NewBWL9(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoBWL9, false)
	return t
}

func NewDKW15(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoDKW15, false)
	return t
}

func NewK15(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoK15, false)
	return t
}

func NewW10Right(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoW10R, false)
	return t
}

func NewW10Left(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoW10L, false)
	return t
}

func NewDKW10(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoDKW10, false)
	return t
}

func NewD2(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoD2, false)
	return t
}

func NewD8(l *TrackLayer, id int) *Track {
	InitRoco()
	t := NewTrack(l, id, rocoD8, false)
	return t
}
