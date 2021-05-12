package types

// Tile 地图瓦片
type Tile struct {
	X, Y, Z int
}

type Point struct {
	X, Y float64
}

type TilePoint struct {
	Tile
	Offset *Point
}
