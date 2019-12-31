package utils

import "reflect"

func IsNil(v interface{}) bool {
	// if v's type is interface and v is nil
	if v == nil {
		return true
	}

	// if v's type is struct and v is nil
	value := reflect.ValueOf(v)

	return value.IsNil()
}

func IsNotNil(v interface{}) bool {
	return !IsNil(v)
}
