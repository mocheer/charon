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
	VersionAPI int            `json:"versionApi"`
	VersionLib int            `json:"versionLib"`
	Remark     string         `json:"remark"`
	CreateAt   time.Time      `json:"createAt"`
	UpdateAt   time.Time      `json:"updateAt"`
	DeletedAt  gorm.DeletedAt `json:"deletedAt"`
}

// TableName 设置表名
func (AppConfig) TableName() string {
	return "pipal.config_app"
}
