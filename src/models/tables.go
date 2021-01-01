package models

import "github.com/mocheer/charon/src/models/tables"

type Entity interface {
	TableName() string
}

var TableMapGenerate = map[string]func() Entity{
	"stipule": func() Entity { return &tables.Stipule{} },
	"petiole": func() Entity { return &tables.Petiole{} },
	"blade":   func() Entity { return &tables.Blade{} },
}
