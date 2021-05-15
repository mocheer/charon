package tables

import (
	"gorm.io/datatypes"
)

// ViewConfig 模块配置
type ViewConfig struct {
	Name    string         `json:"name" gorm:"primary_key"`
	Props   datatypes.JSON `json:"props"`
	Options datatypes.JSON `json:"options"`
	Remark  string         `json:"remark"`
}

// TableName 设置表名
func (ViewConfig) TableName() string {
	return "pipal.config_view"
}
