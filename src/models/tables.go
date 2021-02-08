package models

import "github.com/mocheer/charon/src/models/tables"

// Entity 实体抽象类
type Entity interface {
	TableName() string
}

// TableMapGenerate 映射
var TableMapGenerate = map[string]func() Entity{
	"app":  func() Entity { return &tables.Stipule{} },
	"page": func() Entity { return &tables.Petiole{} },
	"view": func() Entity { return &tables.Blade{} },
}
