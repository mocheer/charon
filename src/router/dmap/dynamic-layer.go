package dmap

import (
	"bytes"
	"image/jpeg"

	"github.com/mocheer/charon/src/core/fs"
	"github.com/mocheer/charon/src/models/types"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/rasterizer"
)

// NewDynamicLayer 动态图层服务
func NewDynamicLayer(id int, tile *types.Tile) *DynamicLayer {

	return &DynamicLayer{ID: id}

}

// DrawImage 绘制图片
func (layer *DynamicLayer) DrawImage(ctx *canvas.Context, path string, x float64, y float64) {
	image, _ := fs.GetImageFromFilePath(path)
	ctx.DrawImage(x, y, image, 1)
}

// GetData 获取图层数据
func (layer *DynamicLayer) GetData() []byte {
	c := canvas.New(256, 256)
	ctx := canvas.NewContext(c)

	layer.DrawImage(ctx, "xxx.png", 0, 0)

	image := rasterizer.Draw(c, 1)
	image.Bounds()
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, image, nil)
	// global.Db.First()
	// DrawImage(ctx, "xxx.png", 0, 0)
	return nil
}
