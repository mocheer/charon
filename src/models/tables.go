package models

import (
	"reflect"

	"github.com/mocheer/charon/src/models/tables"
)

// Has 判断对象是否包含某一个类型
func Has(obj interface{}, key string) bool {
	t := reflect.TypeOf(obj).Elem()
	_, flag := t.FieldByName(key)
	return flag
}

// NewTableStruct 实例化结构体
func NewTableStruct(name string) IEntity {
	return tableMapFunc[name]()
}

// NewTableStructArray
func NewTableStructArray(name string) interface{} {
	switch name {
	case "app":
		var result []tables.AppConfig
		return &result
	case "page":
		var result []tables.PageConfig
		return &result
	case "view":
		var result []tables.ViewConfig
		return &result
	case "feature":
		var result []tables.DmapFeature
		return &result
	default:
		return []map[string]interface{}{}
	}
}

// GenerateTableStruct 实例化结构体
func GenerateTableStruct(name string) reflect.Value {
	entity := NewTableStruct(name)
	t := reflect.ValueOf(entity).Type()
	v := reflect.New(t).Elem()
	return v
}

// tableMapFunc 映射
var tableMapFunc = map[string]func() IEntity{
	"app":     func() IEntity { return &tables.AppConfig{} },
	"page":    func() IEntity { return &tables.PageConfig{} },
	"view":    func() IEntity { return &tables.ViewConfig{} },
	"dmap":    func() IEntity { return &tables.Dmap{} },
	"layer":   func() IEntity { return &tables.DmapLayer{} },
	"feature": func() IEntity { return &tables.DmapFeature{} },
}
