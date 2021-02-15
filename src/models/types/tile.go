package types

// Tile 地图瓦片
type Tile struct {
	Z, Y, X int
}

// TileCache
type TileCache struct {
	TileColumnSize, TileRowSize, ColsPerFile, RowsPerFile, EpsgCode int
}
