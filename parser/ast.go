package parser

import (
	"encoding/json"

	"github.com/weistn/ferrovia/errlog"
)

type File struct {
	Statements []*Statement    `json:"statements"`
	Location   errlog.Location `json:"-"`
}

type Statement struct {
	Layer       *Layer       `json:"layer"`
	RailWay     *RailWay     `json:"way"`
	GroundPlate *GroundPlate `json:"ground"`
	Layout      *Layout      `json:"layout"`
}

type GroundPlate struct {
	Top      float64              `json:"top"`
	Left     float64              `json:"left"`
	Width    float64              `json:"width"`
	Height   float64              `json:"height"`
	Polygon  []GroundPoint        `json:"polygon,omitempty"`
	Location errlog.LocationRange `json:"-"`
}

type GroundPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Layout struct {
	RawText  string
	LocationToken errlog.LocationRange `json:"-"`
	LocationText errlog.LocationRange `json:"-"`
}

type Layer struct {
	Name     string               `json:"name"`
	Location errlog.LocationRange `json:"-"`
}

type RailWay struct {
	Name     string
	Rows     []*ExpressionRow     `json:"rows"`
	Layer    string               `json:"layer"`
	Location errlog.LocationRange `json:"-"`
}

type ExpressionRow struct {
	Expressions []*Expression `json:"exp"`
}

type Expression struct {
	Placeholder      bool
	TrackTermination *TrackTerminationExpression `json:"term"`
	Switch           *SwitchExpression           `json:"switch"`
	Repeat           *RepeatExpression           `json:"repeat"`
	Track            *TrackExpression            `json:"track"`
	TrackMark        *TrackMarkExpression        `json:"mark"`
	Anchor           *AnchorExpression           `json:"anchor"`
	Location         errlog.LocationRange        `json:"-"`
}

type SwitchExpression struct {
	TrackExpression `json:"track"`
	PositionLeft    int  `json:"posleft"`
	JoinLeft        bool `json:"joinleft"`
	SplitLeft       bool `json:"splitleft"`
	PositionRight   int  `json:"posright"`
	JoinRight       bool `json:"joinright"`
	SplitRight      bool `json:"splitright"`
}

type RepeatExpression struct {
	Count           int         `json:"n"`
	TrackExpression *Expression `json:"exp"`
}

type TrackExpression struct {
	Type string `json:"t"`
	// A Token or a more complex expression
	Parameters []interface{} `json:"params"`
}

type TrackMarkExpression struct {
	Name     string  `json:"name"`
	Position float32 `json:"pos"`
}

// In case of "end", EllipsisLeft and EllipsisRight are nil.
type TrackTerminationExpression struct {
	EllipsisLeft  bool
	Name          string
	EllipsisRight bool
}

type AnchorExpression struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Z     float64 `json:"z"`
	Angle float64 `json:"a"`
}

func Load(data []byte) (*File, error) {
	main := &File{}
	err := json.Unmarshal(data, main)
	return main, err
}
