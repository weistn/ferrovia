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

type Layer struct {
	Name     string               `json:"name"`
	Location errlog.LocationRange `json:"-"`
}

type RailWay struct {
	Expressions []*Expression        `json:"exp"`
	Layer       string               `json:"layer"`
	Location    errlog.LocationRange `json:"-"`
}

type Expression struct {
	Repeat         *RepeatExpression         `json:"repeat"`
	Track          *TrackExpression          `json:"track"`
	TrackMark      *TrackMarkExpression      `json:"mark"`
	ConnectionMark *ConnectionMarkExpression `json:"con"`
	Anchor         *AnchorExpression         `json:"anchor"`
}

type RepeatExpression struct {
	Count            int           `json:"n"`
	TrackExpressions []*Expression `json:"exp"`
}

type TrackExpression struct {
	Type             string                `json:"t"`
	JunctionsOnLeft  []*JunctionExpression `json:"jleft,omitempty"`
	JunctionsOnRight []*JunctionExpression `json:"jright,omitempty"`
	Location         errlog.LocationRange  `json:"-"`
}

type JunctionExpression struct {
	Arrow       TokenKind     `json:"arrow"`
	Expressions []*Expression `json:"exp"`
}

type ConnectionMarkExpression struct {
	Name     string               `json:"name"`
	Location errlog.LocationRange `json:"-"`
}

type TrackMarkExpression struct {
	Name     string               `json:"name"`
	Position float32              `json:"pos"`
	Location errlog.LocationRange `json:"-"`
}

type AnchorExpression struct {
	X        float64              `json:"x"`
	Y        float64              `json:"y"`
	Z        float64              `json:"z"`
	Angle    float64              `json:"a"`
	Location errlog.LocationRange `json:"-"`
}

func Load(data []byte) (*File, error) {
	main := &File{}
	err := json.Unmarshal(data, main)
	return main, err
}
