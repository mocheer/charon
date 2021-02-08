package tables

import (
	"gorm.io/datatypes"
)

// Blade 模块配置
type Blade struct {
	Name    string         `json:"name" gorm:"primary_key"`
	Props   datatypes.JSON `json:"props"`
	Options datatypes.JSON `json:"options"`
	Remark  string         `json:"remark"`
}

// TableName 设置表名
func (Blade) TableName() string {
	return "pipal.config_view"
}
