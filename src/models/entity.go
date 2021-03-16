package models

import "github.com/mocheer/charon/src/core/fn"

// IEntity 实体抽象类
type IEntity interface {
	TableName() string
}

type Entity struct{}

func (e *Entity) toMap() (map[string]interface{}, error) {
	return fn.StructToMap(e, "map", "")
}
