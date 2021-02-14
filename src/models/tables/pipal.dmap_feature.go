package tables

// DmapFeature 应用配置
type DmapFeature struct {
	LayerID int `json:"layer_id"`
	ID      int `json:"name"`
}

// TableName 设置表名
func (DmapFeature) TableName() string {
	return "pipal.dmap_feature"
}
