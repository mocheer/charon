package fn

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	methodResNum = 2
)

const (
	OptIgnore    = "-"         // 忽略当前这个域
	OptOmitempty = "omitempty" // 当这个域的值为空，忽略这个域
	OptDive      = "dive"      // 递归地遍历这个结构体，将所有字段作为键
	OptWildcard  = "wildcard"  // 只适用于字符串类型，返回"%"+值+"%"，这是为了方便数据库的模糊查询
)

const (
	flagIgnore = 1 << iota
	flagOmiEmpty
	flagDive
	flagWildcard
)

// StructToMap 结构体转map
func StructToMap(s interface{}, tag string, methodName string) (res map[string]interface{}, err error) {
	v := getValue(s)
	t := getType(v)
	//
	if t == nil {
		return nil, fmt.Errorf("结构体参数有问题")
	}
	//
	res = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)

		// 忽略没有导出的字段
		if fieldType.PkgPath != "" {
			continue
		}
		// read tag
		tagVal, flag := readTag(fieldType, tag)
		// 忽略字段
		if flag&flagIgnore != 0 {
			continue
		}

		fieldValue := v.Field(i)
		// 为空时忽略
		if flag&flagOmiEmpty != 0 && fieldValue.IsZero() {
			continue
		}

		// 忽略空指针
		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			continue
		}
		if fieldValue.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}

		// get kind
		switch fieldValue.Kind() {
		case reflect.Slice, reflect.Array:
			if methodName != "" {
				_, ok := fieldValue.Type().MethodByName(methodName)
				if ok {
					key, value, err := callFunc(fieldValue, methodName)
					if err != nil {
						return nil, err
					}
					res[key] = value
					continue
				}
			}
			res[tagVal] = fieldValue
		case reflect.Struct:
			if methodName != "" {
				_, ok := fieldValue.Type().MethodByName(methodName)
				if ok {
					key, value, err := callFunc(fieldValue, methodName)
					if err != nil {
						return nil, err
					}
					res[key] = value
					continue
				}
			}

			// recursive
			deepRes, deepErr := StructToMap(fieldValue.Interface(), tag, methodName)
			if deepErr != nil {
				return nil, deepErr
			}
			if flag&flagDive != 0 {
				for k, v := range deepRes {
					res[k] = v
				}
			} else {
				res[tagVal] = deepRes
			}
		case reflect.Map:
			res[tagVal] = fieldValue
		case reflect.Chan:
			res[tagVal] = fieldValue
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64:
			res[tagVal] = fieldValue.Int()
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
			res[tagVal] = fieldValue.Uint()
		case reflect.Float32, reflect.Float64:
			res[tagVal] = fieldValue.Float()
		case reflect.String:
			if flag&flagWildcard != 0 {
				res[tagVal] = "%" + fieldValue.String() + "%"
			} else {
				res[tagVal] = fieldValue.String()
			}
		case reflect.Bool:
			res[tagVal] = fieldValue.Bool()
		case reflect.Complex64, reflect.Complex128:
			res[tagVal] = fieldValue.Complex()
		case reflect.Interface:
			res[tagVal] = fieldValue.Interface()
		default:
		}
	}
	return
}

func getValue(s interface{}) reflect.Value {
	return reflect.ValueOf(s)
}

func getType(v reflect.Value) reflect.Type {
	// 空指针
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	// 只支持结构体
	if v.Kind() != reflect.Struct {
		return nil
	}

	t := v.Type()
	return t
}

// readTag 读取标签值
func readTag(f reflect.StructField, tag string) (string, int) {
	val, ok := f.Tag.Lookup(tag)
	fieldTag := ""
	flag := 0

	// no tag, skip this field
	if !ok {
		flag |= flagIgnore
		return "", flag
	}
	opts := strings.Split(val, ",")

	fieldTag = opts[0]
	for i := 0; i < len(opts); i++ {
		switch opts[i] {
		case OptIgnore:
			flag |= flagIgnore
		case OptOmitempty:
			flag |= flagOmiEmpty
		case OptDive:
			flag |= flagDive
		case OptWildcard:
			flag |= flagWildcard
		}
	}

	return fieldTag, flag
}

// 当结构体自身实现了toMap方法时，
func callFunc(fv reflect.Value, methodName string) (string, interface{}, error) {
	methodRes := fv.MethodByName(methodName).Call([]reflect.Value{})
	if len(methodRes) != methodResNum {
		return "", nil, fmt.Errorf("wrong method %s, should have 2 output: (string,interface{})", methodName)
	}
	if methodRes[0].Kind() != reflect.String {
		return "", nil, fmt.Errorf("wrong method %s, first output should be string", methodName)
	}
	key := methodRes[0].String()
	return key, methodRes[1], nil
}
