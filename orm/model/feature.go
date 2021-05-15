package model

import "gorm.io/datatypes"

type GeoFeature struct {
	Type        string         `json:"type"`
	Coordinates datatypes.JSON `json:"coordinates"`
	Properties  datatypes.JSON `json:"properties"`
}
