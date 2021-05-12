package model

import "github.com/mocheer/pluto/fn"

// IEntity 实体抽象类
type IEntity interface {
	TableName() string
}

type Entity struct{}

// 转成 map[string]interface
func (e *Entity) ToMap() (map[string]interface{}, error) {
	return fn.StructToMap(e, "map", "")
}

// 重置结构体
func (e *Entity) Reset() {
	*e = Entity{}
}
