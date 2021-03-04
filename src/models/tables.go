package models

import (
	"reflect"

	"github.com/mocheer/charon/src/models/tables"
)

// Entity 实体抽象类
type Entity interface {
	TableName() string
}

// Has 判断对象是否包含某一个类型
func Has(obj interface{}, key string) bool {
	t := reflect.TypeOf(obj).Elem()
	_, flag := t.FieldByName(key)
	return flag
}

// NewTableStruct 实例化结构体
func NewTableStruct(name string) Entity {
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
var tableMapFunc = map[string]func() Entity{
	"app":     func() Entity { return &tables.AppConfig{} },
	"page":    func() Entity { return &tables.PageConfig{} },
	"view":    func() Entity { return &tables.ViewConfig{} },
	"dmap":    func() Entity { return &tables.Dmap{} },
	"layer":   func() Entity { return &tables.DmapLayer{} },
	"feature": func() Entity { return &tables.DmapFeature{} },
}
