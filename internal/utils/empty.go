package utils

import "reflect"

// IsEmptyStruct 判断结构体是否为空
func IsEmptyStruct(s interface{}) bool {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Struct {
		return v.IsZero()
	}
	return true // 如果不是结构体返回 true
}
