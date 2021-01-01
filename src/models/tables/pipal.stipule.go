package tables

import (
	"time"

	"gorm.io/datatypes"
)

// Stipule 应用配置
type Stipule struct {
	Name       string         `json:"name"`
	Title      string         `json:"title"`
	Enabled    bool           `json:"enabled"`
	Theme      string         `json:"theme"`
	Options    datatypes.JSON `json:"options"`
	Version    string         `json:"version"`
	CreateTime time.Time      `json:"create_time"`
	UpdateTime time.Time      `json:"update_time"`
	Remark     string         `json:"remark"`
}

// TableName 设置表名
func (Stipule) TableName() string {
	return "pipal.stipule"
}
