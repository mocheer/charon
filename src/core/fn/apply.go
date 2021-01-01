package fn

import "reflect"

// Apply 通过反射支持字符串直接调用函数
func Apply(obj_func interface{}, args []interface{}) []reflect.Value {
	fn := reflect.ValueOf(obj_func)
	in := make([]reflect.Value, len(args))
	for k, param := range args {
		in[k] = reflect.ValueOf(param)
	}
	r := fn.Call(in)
	return r
}
