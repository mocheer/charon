package orm

import (
	"github.com/mocheer/charon/orm/tables"
	"github.com/mocheer/pluto/ts"
)

// tableMap 映射
var TableMap = ts.Map{
	"app":     &tables.AppConfig{},
	"page":    &tables.PageConfig{},
	"view":    &tables.ViewConfig{},
	"dmap":    &tables.Dmap{},
	"layer":   &tables.DmapLayer{},
	"feature": &tables.DmapFeature{},
}
