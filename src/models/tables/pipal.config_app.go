package tables

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// AppConfig 应用配置
type AppConfig struct {
	Name       string         `json:"name" gorm:"primary_key"`
	Title      string         `json:"title"`
	Enabled    bool           `json:"enabled"`
	Theme      string         `json:"theme"`
	Options    datatypes.JSON `json:"options"`
	VersionAPI int            `json:"version_api"`
	VersionLib int            `json:"version_lib"`
	Remark     string         `json:"remark"`
	CreateAt   time.Time      `json:"create_time"`
	UpdateAt   time.Time      `json:"update_time"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
}

// TableName 设置表名
func (AppConfig) TableName() string {
	return "pipal.config_app"
}
