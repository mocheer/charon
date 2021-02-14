package tables

import "gorm.io/datatypes"

// DmapLayer 应用配置
type DmapLayer struct {
	ID         int            `json:"id" gorm:"primary_key"`
	ParentID   int            `json:"parent_id"`
	Name       string         `json:"name"`
	CRS        string         `json:"crs"`
	Type       string         `json:"type"`
	Extent     datatypes.JSON `json:"extent"`
	Options    datatypes.JSON `json:"options"`
	Properties datatypes.JSON `json:"properties"`
}

// TableName 设置表名
func (DmapLayer) TableName() string {
	return "pipal.dmap_layer"
}
