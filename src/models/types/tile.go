package types

// Tile 地图瓦片
type Tile struct {
	Row, Level, Column, EpsgCode int
}

// TileCache
type TileCache struct {
	HasTransparency                                       bool
	TileColumnSize, TileRowSize, ColsPerFile, RowsPerFile int
	EpsgCode, MinLevel, MaxLevel                          int
}
