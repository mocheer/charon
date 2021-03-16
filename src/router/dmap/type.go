package dmap

import "github.com/mocheer/charon/src/models/orm"

// @see https://github.com/shengzheng1981/green-gis-server
// Bbox 视野范围
type Bbox struct {
	MinX float64
	MinY float64
	MaxX float64
	MaxY float64
}

// QueryParams query服务的输入参数
type QueryParams struct {
	orm.SelectBuilder
	GeometryPrecision int    `json:"geometryPrecision"` //几何数据精度
	F                 string `json:"f"`                 //返回的数据格式
}

// IdentifyParams identify服务的输入参数
type IdentifyParams struct {
	Geometry  string   `json:"geometry"`  //几何对象
	Layers    []string `json:"layers"`    //识别的图层组
	Tolerance int      `json:"tolerance"` //能够接受的误差
	F         string   `json:"f"`         //返回的数据格式
}
