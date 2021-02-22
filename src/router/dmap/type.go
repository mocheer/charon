package dmap

// DynamicLayer 动态图层
type DynamicLayer struct {
	ID   int
	Bbox Bbox
}

// Bbox 视野范围
type Bbox struct {
	minX float64
	minY float64
	maxX float64
	maxY float64
}
