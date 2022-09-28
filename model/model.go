package model

import (
	"github.com/weistn/ferrovia/errlog"
	"github.com/weistn/ferrovia/model/switchboard"
	"github.com/weistn/ferrovia/model/tracks"
)

type Model struct {
	Name         string
	GroundPlates []*GroundPlate
	Switchboards []*switchboard.ASCIISwitchboard
	Tracks       *tracks.TrackSystem
	Location     errlog.LocationRange
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
	X        float64              `json:"x"`
	Y        float64              `json:"y"`
	Location errlog.LocationRange `json:"-"`
}
