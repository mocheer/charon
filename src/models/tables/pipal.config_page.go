package tables

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// PageConfig 页面配置
type PageConfig struct {
	AppName   string         `json:"app_name" gorm:"primary_key"`
	Name      string         `json:"name" gorm:"primary_key"`
	Dep       string         `json:"dep"`
	Guard     bool           `json:"guard"`
	Props     datatypes.JSON `json:"props"`
	Options   datatypes.JSON `json:"options"`
	Mime      string         `json:"mime"`
	Log       string         `json:"log"`
	Remark    string         `json:"remark"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

// TableName 设置表名
func (PageConfig) TableName() string {
	return "pipal.config_page"
}
