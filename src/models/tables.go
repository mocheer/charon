package models

import "github.com/mocheer/charon/src/models/tables"

// Entity 实体抽象类
type Entity interface {
	TableName() string
}

// func createEntityFactory(EntityStruct Entity) func() Entity {
// 	return func() Entity {
// 		return &EntityStruct{}
// 	}
// }

// TableMapGenerate 映射
var TableMapGenerate = map[string]func() Entity{
	"app":     func() Entity { return &tables.AppConfig{} },
	"page":    func() Entity { return &tables.PageConfig{} },
	"view":    func() Entity { return &tables.ViewConfig{} },
	"dmap":    func() Entity { return &tables.Dmap{} },
	"layer":   func() Entity { return &tables.DmapLayer{} },
	"feature": func() Entity { return &tables.DmapFeature{} },
}
