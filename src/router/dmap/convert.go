package dmap

import (
	"math"

	"github.com/mocheer/charon/src/models/types"
)

// LonLat2Tile
func LonLat2Tile(lon float64, lat float64, z float64) (tilePoint *types.TilePoint) {
	scaleZ := math.Exp2(z)
	tileTempX := (lon + 180.0) / 360.0 * scaleZ
	tileTempY := math.Log(math.Tan(lat*math.Pi/180*0.5+0.25*math.Pi)) / (2.0 * math.Pi)
	//
	tileX := math.Floor(tileTempX)
	tileY := math.Floor((0.5 - tileTempY) * scaleZ)
	//
	pixelX := int(tileTempX*256.0) % 256
	pixelY := int((1.0-tileTempY)*scaleZ*256.0) % 256
	//
	offsetPoint := &types.Point{
		X: float64(pixelX),
		Y: float64(pixelY),
	}

	tilePoint = &types.TilePoint{
		Tile: types.Tile{
			X: int(tileX),
			Y: int(tileY),
			Z: int(z),
		},
		Offset: offsetPoint,
	}

	return
}
