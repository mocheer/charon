package dmap

import (
	"bytes"
	"image/jpeg"

	"github.com/mocheer/charon/src/core/fs"
	"github.com/mocheer/charon/src/models/types"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/rasterizer"
	"github.com/twpayne/go-geom"
)

// DynamicLayer 动态图层
type DynamicLayer struct {
	ID   int
	Bbox Bbox
}

// NewDynamicLayer 动态图层服务
func NewDynamicLayer(id int, tile *types.Tile) *DynamicLayer {
	return &DynamicLayer{ID: id}
}

// Draw 绘制
func (layer *DynamicLayer) Draw(data *geom.GeometryCollection) []byte {
	for _, g := range data.Geoms() {
		switch g.(type) {
		case *geom.Point:
		case *geom.LineString:
		case *geom.LinearRing:
		case *geom.MultiLineString:
		case *geom.MultiPoint:
		case *geom.MultiPolygon:
		case *geom.Polygon:
		case *geom.GeometryCollection:
		default:

		}

		geom.NewGeometryCollection()
	}
	return layer.GetData()
}

// DrawImage 绘制图片
func (layer *DynamicLayer) DrawImage(ctx *canvas.Context, path string, x float64, y float64) []byte {
	image, _ := fs.GetImageFromFilePath(path)
	ctx.DrawImage(x, y, image, 1)

	return layer.GetData()
}

// GetData 获取图层数据
func (layer *DynamicLayer) GetData() []byte {
	c := canvas.New(256, 256)
	ctx := canvas.NewContext(c)
	//
	layer.DrawImage(ctx, "xxx.png", 0, 0)

	image := rasterizer.Draw(c, 1)
	image.Bounds()
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, image, nil)
	// global.Db.First()
	// DrawImage(ctx, "xxx.png", 0, 0)
	return buf.Bytes()
}
