package tables

import "gorm.io/datatypes"

// Dmap 应用配置
type Dmap struct {
	Name    string         `json:"name" gorm:"primary_key"`
	Props   datatypes.JSON `json:"props"`
	Options datatypes.JSON `json:"options"`
	Widgets datatypes.JSON `json:"widgets"`
	Layers  datatypes.JSON `json:"layers"`
}

// TableName 设置表名
func (Dmap) TableName() string {
	return "pipal.dmap"
}
