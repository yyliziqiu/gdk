package xutil

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

var FuncNamePrefixes []string

// ReflectFuncName 通过反射获取方法名称
func ReflectFuncName(f any) string {
	pc := runtime.FuncForPC(reflect.ValueOf(f).Pointer())
	name := strings.Split(pc.Name(), "-")[0]

	for _, prefix := range FuncNamePrefixes {
		if strings.HasPrefix(name, prefix) {
			return strings.TrimPrefix(name, prefix)
		}
	}

	return name
}

// ReflectFieldList 通过反射获取指定结构体所有字段名
func ReflectFieldList(s any) []string {
	mt := reflect.TypeOf(s)
	var fields []string
	for i := 0; i < mt.NumField(); i++ {
		fields = append(fields, mt.Field(i).Name)
	}
	return fields
}

// ReflectValueList 通过反射获取指定结构体所有字段值
func ReflectValueList(s any) []any {
	mv := reflect.ValueOf(s)
	var values []any
	for i := 0; i < mv.NumField(); i++ {
		values = append(values, mv.Field(i).Interface())
	}
	return values
}

// ReflectValueStringList 通过反射获取指定结构体所有字段值的字符串形式
func ReflectValueStringList(s any) []string {
	var values []string
	for _, v := range ReflectValueList(s) {
		values = append(values, fmt.Sprintf("%v", v))
	}
	return values
}

// ReflectFieldValue 获取指定结构体指定字段的值
func ReflectFieldValue(s any, fn string) (any, bool) {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	field := val.FieldByName(fn)
	if !field.IsValid() {
		return nil, false
	}
	return field.Interface(), true
}

// AttemptReflectFieldValue 依次尝试字段名列表获取指定结构体的字段值
func AttemptReflectFieldValue(s any, fns []string) (any, bool) {
	for _, fn := range fns {
		if val, ok := ReflectFieldValue(s, fn); ok {
			return val, true
		}
	}
	return nil, false
}
