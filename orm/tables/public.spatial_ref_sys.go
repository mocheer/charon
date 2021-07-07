package tables

// SpatialRefSys 坐标系
// postgis拓展创建的表
type SpatialRefSys struct {
	Srid      int
	AuthName  string
	AuthSrid  int
	Srtext    string
	Proj4text string
}

// TableName 设置表名
func (SpatialRefSys) TableName() string {
	return "public.spatial_ref_sys"
}
