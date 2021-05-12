package types

type Geometry struct {
	X, Y int
}

func (g Geometry) GormDataType() string {
	return "geometry"
}

//
func (g *Geometry) Scan(v interface{}) error {
	return nil
}
