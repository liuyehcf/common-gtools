package utils

import "reflect"

func IsNil(v interface{}) bool {
	if v == nil {
		return true
	}

	value := reflect.ValueOf(v)

	switch value.Kind() {
	case reflect.Invalid:
		return true
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer,
		reflect.Interface, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

func IsNotNil(v interface{}) bool {
	return !IsNil(v)
}
