package tables

import (
	"time"

	"gorm.io/gorm"
)

// FsFile 应用配置
type FsFile struct {
	UID       string         `json:"uid"`
	PID       string         `json:"pid"`
	Name      string         `json:"name"`
	Data      string         `json:"data"`
	CreateAt  time.Time      `json:"createAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
}

// TableName 设置表名
func (FsFile) TableName() string {
	return "charon.fs_file"
}
