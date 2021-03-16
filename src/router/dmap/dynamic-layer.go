package dmap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/png"
	"path/filepath"

	"github.com/mocheer/charon/src/core/fs"
	"github.com/mocheer/charon/src/logger"
	"github.com/mocheer/charon/src/models/types"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/rasterizer"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"gorm.io/datatypes"
)

// DynamicLayer 动态图层
type DynamicLayer struct {
	canvas  *canvas.Canvas
	ctx     *canvas.Context
	tile    *types.Tile
	Options struct {
		MinZoom int `json:"minZoom"`
		Feature struct {
			Src string `json:"src"`
		}
	}
}

// NewDynamicLayer 动态图层服务
func NewDynamicLayer(tile *types.Tile) *DynamicLayer {
	c := canvas.New(256, 256)
	ctx := canvas.NewContext(c)
	return &DynamicLayer{
		canvas: c,
		ctx:    ctx,
		tile:   tile,
	}
}

// SetOptions
func (layer *DynamicLayer) SetOptions(options interface{}) {
	json.Unmarshal(options.(datatypes.JSON), &layer.Options)
}

// Draw
func (layer *DynamicLayer) Draw(data []byte) (err error) {
	var g geom.T
	geojson.Unmarshal(data, &g)

	switch g.(type) {
	case *geom.Point:
		layer.DrawPoint(g.(*geom.Point))
	case *geom.LineString:
	case *geom.LinearRing:
	case *geom.MultiLineString:
	case *geom.MultiPoint:
	case *geom.MultiPolygon:
	case *geom.Polygon:
	case *geom.GeometryCollection:
		layer.DrawGeometryCollection(g.(*geom.GeometryCollection))
	default:
		err = fmt.Errorf("数据解析失败，不是符合规范的数据")
	}
	return
}

// DrawGeometryCollection 绘制
func (layer *DynamicLayer) DrawGeometryCollection(data *geom.GeometryCollection) {
	for _, g := range data.Geoms() {
		switch g.(type) {
		case *geom.Point:
			layer.DrawPoint(g.(*geom.Point))
		case *geom.LineString:
		case *geom.LinearRing:
		case *geom.MultiLineString:
		case *geom.MultiPoint:
		case *geom.MultiPolygon:
		case *geom.Polygon:
		}
	}
}

// DrawFeature
func (layer *DynamicLayer) DrawPoint(point *geom.Point) {

	coor := point.Coords()
	tilePoint := LonLat2Tile(coor.X(), coor.Y(), float64(layer.tile.Z))
	offset := tilePoint.Offset
	if tilePoint.X == layer.tile.X && tilePoint.Y == layer.tile.Y && offset.X >= 0 && offset.Y >= 0 && offset.X <= 256 && offset.Y <= 256 {
		layer.drawImage(layer.Options.Feature.Src, offset.X, offset.Y)
	}
	return
}

// DrawImage 绘制图片
func (layer *DynamicLayer) drawCircle(x float64, y float64, r float64) {
	layer.ctx.DrawPath(x, 256-y, canvas.Circle(r))
}

// DrawImage 绘制图片
func (layer *DynamicLayer) drawImage(path string, x float64, y float64) {
	image, err := fs.GetImageFromFilePath(filepath.Join("./public", path))

	if err != nil {
		logger.Error(err)
		return
	}
	layer.ctx.DrawImage(x-16, 256-(y+16), image, 1)
	// layer.ctx.DrawPath(x, 256-y, canvas.Circle(5))
}

func (layer *DynamicLayer) getData() []byte {
	image := rasterizer.Draw(layer.canvas, 1)
	image.Bounds()
	buf := new(bytes.Buffer)
	// jpeg.Encode(buf, image, nil)
	png.Encode(buf, image)

	return buf.Bytes()
}
