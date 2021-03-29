package fn

import (
	"reflect"
)

// GetType 通过反射获取数据类型
// GetType(string) == "string"
func GetType(val interface{}) reflect.Type {
	typ := reflect.TypeOf(val)                   // typ 有可能是指针或者nil
	if typ != nil && typ.Kind() == reflect.Ptr { //
		typ = typ.Elem()
	}
	return typ
}

// IsString
func IsString(val interface{}) bool {
	typ := GetType(val)
	return typ != nil && typ.Kind() == reflect.String
}
