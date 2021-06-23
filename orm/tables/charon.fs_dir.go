package tables

import (
	"time"

	"gorm.io/gorm"
)

// FsFile 应用配置
type FsDir struct {
	UID       string         `json:"uid"`
	PID       string         `json:"pid"`
	Name      string         `json:"name"`
	SNO       int            `json:"sno"`
	CreateAt  time.Time      `json:"createAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
}

// TableName 设置表名
func (FsDir) TableName() string {
	return "charon.fs_dir"
}
