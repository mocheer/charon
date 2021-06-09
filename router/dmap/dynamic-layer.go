package dmap

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"sync"

	"github.com/mocheer/charon/global"
	"github.com/mocheer/pluto/calc"
	"github.com/mocheer/pluto/fn"
	"github.com/mocheer/pluto/fs"
	"github.com/mocheer/pluto/ts"
	"github.com/mocheer/pluto/ts/img"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/rasterizer"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"gorm.io/datatypes"
)

// DynamicLayer 动态图层
type DynamicLayer struct {
	id       int
	canvas   *canvas.Canvas
	ctx      *canvas.Context
	tile     *ts.Tile
	minTile  *ts.Tile
	maxTile  *ts.Tile
	numTileX int
	numTileY int
	NumTile  int
	data     *geom.GeometryCollection
	//
	Options struct {
		MinZoom int `json:"minZoom"`
		MaxZoom int `json:"maxZoom"`
		Feature struct {
			Src    string `json:"src"`
			Radius int    `json:"radius"`
		}
	}
}

// NewDynamicLayer 动态图层服务
func NewDynamicLayer(id int, tile *ts.Tile) *DynamicLayer {
	data := geom.NewGeometryCollection()
	return &DynamicLayer{
		id:      id,
		tile:    tile,
		minTile: tile,
		data:    data,
	}
}

// SetOptions
func (layer *DynamicLayer) SetOptions(options interface{}) {
	json.Unmarshal(options.(datatypes.JSON), &layer.Options)
}

// Add 添加数据
func (layer *DynamicLayer) Add(data []byte) (err error) {
	var g geom.T
	geojson.Unmarshal(data, &g)
	switch g := g.(type) {
	case *geom.GeometryCollection:
		for _, gi := range g.Geoms() {
			layer.data.Push(gi)
		}
	default:
		layer.data.Push(g)
	}
	return
}

// DrawTile 绘制单张瓦片 => 只支持点图层，因为线和面边界稍复杂
func (layer *DynamicLayer) DrawTile() (err error) {

	layer.canvas = canvas.New(256, 256)
	layer.ctx = canvas.NewContext(layer.canvas)
	//
	for _, g := range layer.data.Geoms() {
		switch g := g.(type) {
		case *geom.Point:
			layer.drawPoint(g)
		}
	}
	return
}

// Draw
func (layer *DynamicLayer) Draw() (err error) {
	bounds := layer.data.Bounds()
	minLon := bounds.Min(0)
	minLat := bounds.Min(1)
	maxLon := bounds.Max(0)
	maxLat := bounds.Max(1)
	//
	minTilePoint := calc.LonLat2TilePoint(minLon, maxLat, float64(layer.tile.Z)) //左上瓦片
	maxTilePoint := calc.LonLat2TilePoint(maxLon, minLat, float64(layer.tile.Z)) //右下瓦片
	//
	if minTilePoint.Offset.X < 16 {
		minTilePoint.X -= 1
	}

	if minTilePoint.Offset.Y < 16 {
		minTilePoint.Y -= 1
	}

	if maxTilePoint.Offset.X > (256 - 16) {
		maxTilePoint.X += 1
	}

	if maxTilePoint.Offset.Y > (256 - 16) {
		maxTilePoint.Y += 1
	}

	//
	numTileX := maxTilePoint.X - minTilePoint.X + 1
	numTileY := maxTilePoint.Y - minTilePoint.Y + 1
	layer.numTileX = numTileX
	layer.numTileY = numTileY
	layer.NumTile = numTileX * numTileY
	//
	layer.minTile = &ts.Tile{X: minTilePoint.X, Y: minTilePoint.Y, Z: minTilePoint.Z}
	layer.maxTile = &ts.Tile{X: maxTilePoint.X, Y: maxTilePoint.Y, Z: maxTilePoint.Z}
	//
	layer.canvas = canvas.New(float64(numTileX)*256, float64(numTileY)*256)
	layer.ctx = canvas.NewContext(layer.canvas)

	for _, g := range layer.data.Geoms() {
		switch g := g.(type) {
		case *geom.Point:
			layer.drawPoint(g)
		case *geom.LineString:
			layer.drawPolyline(g)
		case *geom.LinearRing:
		case *geom.Polygon:
			layer.drawPolygon(g)
		case *geom.MultiLineString:
		case *geom.MultiPoint:
		case *geom.MultiPolygon:
		}
	}

	return
}

// Coor2Pixel
func (layer *DynamicLayer) Coor2Pixel(coor geom.Coord) *ts.Point {
	tilePoint := calc.LonLat2TilePoint(coor.X(), coor.Y(), float64(layer.tile.Z))
	offset := tilePoint.Offset
	offset.X += float64(tilePoint.X-layer.minTile.X) * 256
	offset.Y += float64(tilePoint.Y-layer.minTile.Y) * 256
	// canvas的Y坐标轴方向跟浏览器是相反的
	offset.Y = layer.canvas.H - offset.Y
	return offset
}

// drawPoint
func (layer *DynamicLayer) drawPoint(p *geom.Point) {
	offset := layer.Coor2Pixel(p.Coords())
	//
	switch layer.Options.Feature.Src {
	case "rect":
	case "circle":
		layer.drawCircle(offset.X, offset.Y, float64(layer.Options.Feature.Radius))
	default:
		tolerance := 16.0 // 图片大小的一半
		if offset.X >= -tolerance && offset.Y >= -tolerance && offset.X <= (layer.canvas.W+tolerance) && offset.Y <= (layer.canvas.H+tolerance) {
			layer.drawImage(layer.Options.Feature.Src, offset.X, offset.Y)
		}
	}
}

// drawPolyline
func (layer *DynamicLayer) drawPolyline(line *geom.LineString) {
	coors := line.Coords()
	p := &canvas.Path{}
	//
	layer.ctx.StrokeWidth = 1
	layer.ctx.SetStrokeColor(color.RGBA{0, 0, 255, 255})
	layer.ctx.SetFillColor(color.Transparent)
	for i, coor := range coors {
		offset := layer.Coor2Pixel(coor)
		if i == 0 {
			p.MoveTo(offset.X, offset.Y)
		}
		p.LineTo(offset.X, offset.Y)
	}
	//
	layer.ctx.DrawPath(0, 0, p)
}

// drawPolygon
func (layer *DynamicLayer) drawPolygon(polygon *geom.Polygon) {
	coors2 := polygon.Coords()
	p := &canvas.Path{}
	// 边框大小
	layer.ctx.SetStrokeWidth(1)
	// 边框颜色
	layer.ctx.SetStrokeColor(color.RGBA{0, 0, 255, 255})
	// 填充颜色
	layer.ctx.SetFillColor(color.RGBA{255, 0, 0, 100})
	//
	for _, coors := range coors2 {
		for j, coor := range coors {
			offset := layer.Coor2Pixel(coor)
			if j == 0 {
				p.MoveTo(offset.X, offset.Y)
			}
			p.LineTo(offset.X, offset.Y)
		}
	}
	p.Close()
	//
	layer.ctx.DrawPath(0, 0, p)
}

// drawCircle 绘制圆
func (layer *DynamicLayer) drawCircle(x float64, y float64, r float64) {
	layer.ctx.DrawPath(x, y, canvas.Circle(r))
}

// drawImage 绘制图片
func (layer *DynamicLayer) drawImage(path string, x float64, y float64) {
	im, err := img.FromFile(filepath.Join(global.Config.FirstStaticDir(), path))

	if err != nil {
		global.Log.Error(err)
		return
	}
	layer.ctx.DrawImage(x-16, y-16, im.Image, 1)
}

// GetData 获取画布数据
func (layer *DynamicLayer) GetData() []byte {
	image := rasterizer.Draw(layer.canvas, 1)
	bs, _ := img.ToBytes(image, "png")
	return bs
}

// savingTile 保存瓦片
func (layer *DynamicLayer) savingTile(imgRGBA *image.RGBA, i, j int) string {
	minTile := layer.minTile
	imgTilePath := fmt.Sprintf(ImageTilePathFormat, layer.id, minTile.Z, j, i)
	if !fs.IsExist(imgTilePath) {
		x0 := (i - minTile.X) * 256
		y0 := (j - minTile.Y) * 256
		subImg := imgRGBA.SubImage(image.Rect(x0, y0, x0+256, y0+256))
		f, _ := fs.OpenOrCreate(imgTilePath, os.O_CREATE|os.O_WRONLY, 0666)
		png.Encode(f, subImg)
		f.Close()
	}
	return imgTilePath
}

// SaveTiles 保存瓦片
func (layer *DynamicLayer) SaveTiles() *sync.WaitGroup {
	img := rasterizer.Draw(layer.canvas, 1)
	minTile := layer.minTile
	maxTile := layer.maxTile

	tasks := make([]func(), layer.NumTile)
	createTask := func(img *image.RGBA, i, j int) func() {
		return func() {
			layer.savingTile(img, i, j)
		}
	}
	count := 0
	for i := minTile.X; i <= maxTile.X; i++ {
		for j := minTile.Y; j <= maxTile.Y; j++ {
			tasks[count] = createTask(img, i, j)
			count++
		}
	}
	// 并发执行
	return fn.GoFns(16, tasks)
}

// GetTile
func (layer *DynamicLayer) GetTile(i, j int) string {
	img := rasterizer.Draw(layer.canvas, 1)
	return layer.savingTile(img, i, j)
}
