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
	pn := strings.Split(pc.Name(), "-")[0]
	for _, prefix := range FuncNamePrefixes {
		if strings.HasPrefix(pn, prefix) {
			return strings.TrimPrefix(pn, prefix)
		}
	}
	return pn
}

// ReflectFields 通过反射获取指定结构体所有字段名
func ReflectFields(s any) []string {
	st := reflect.TypeOf(s)
	var fns []string
	for i := 0; i < st.NumField(); i++ {
		fns = append(fns, st.Field(i).Name)
	}
	return fns
}

// ReflectValues 通过反射获取指定结构体所有字段值
func ReflectValues(s any) []any {
	sv := reflect.ValueOf(s)
	var vs []any
	for i := 0; i < sv.NumField(); i++ {
		vs = append(vs, sv.Field(i).Interface())
	}
	return vs
}

// ReflectValuesToString 通过反射获取指定结构体所有字段值的字符串形式
func ReflectValuesToString(s any) []string {
	var vs []string
	for _, v := range ReflectValues(s) {
		vs = append(vs, fmt.Sprintf("%v", v))
	}
	return vs
}

// ReflectValueByField 获取指定结构体指定字段的值
func ReflectValueByField(s any, fn string) (any, bool) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	f := v.FieldByName(fn)
	if !f.IsValid() {
		return nil, false
	}
	return f.Interface(), true
}

// AttemptReflectValueByField 依次尝试字段名列表获取指定结构体的字段值
func AttemptReflectValueByField(s any, fns []string) (any, bool) {
	for _, fn := range fns {
		if v, ok := ReflectValueByField(s, fn); ok {
			return v, true
		}
	}
	return nil, false
}

// ReflectValueList deprecated
func ReflectValueList(s any) []any {
	return ReflectValues(s)
}

// ReflectValueStringList deprecated
func ReflectValueStringList(s any) []string {
	return ReflectValuesToString(s)
}

// ReflectFieldValue deprecated
func ReflectFieldValue(s any, fn string) (any, bool) {
	return ReflectValueByField(s, fn)
}

// AttemptReflectFieldValue deprecated
func AttemptReflectFieldValue(s any, ns []string) (any, bool) {
	return AttemptReflectValueByField(s, ns)
}
