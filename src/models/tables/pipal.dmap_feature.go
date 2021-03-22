package tables

import (
	"github.com/mocheer/charon/src/models/ormtypes"
	"gorm.io/datatypes"
)

// DmapFeature 应用配置
type DmapFeature struct {
	LayerID    int               `json:"layer_id"`
	ID         int               `json:"id"`
	Geometry   ormtypes.Geometry `json:"geometry"`
	Options    datatypes.JSON    `json:"options"`
	Properties datatypes.JSON    `json:"properties"`
}

// TableName 设置表名
func (DmapFeature) TableName() string {
	return "pipal.dmap_feature"
}
