package orm

import (
	"github.com/mocheer/charon/orm/tables"
)

// tableMap 映射
var TableMap = map[string]interface{}{
	"app":     &tables.AppConfig{},
	"page":    &tables.PageConfig{},
	"view":    &tables.ViewConfig{},
	"dmap":    &tables.Dmap{},
	"layer":   &tables.DmapLayer{},
	"feature": &tables.DmapFeature{},
}
