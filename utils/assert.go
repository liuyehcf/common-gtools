package utils

import (
	"errors"
	"fmt"
)

func AssertNil(v interface{}, message string) {
	if IsNotNil(v) {
		if err, ok := v.(error); ok {
			panic(errors.New(fmt.Sprintf("%v, %s", err.Error(), message)))
		} else {
			panic(errors.New(message))
		}
	}
}

func AssertNotNil(v interface{}, message string) {
	if v == nil {
		panic(errors.New(message))
	}
}

func AssertTrue(b bool, message string) {
	if !b {
		panic(errors.New(message))
	}
}

func AssertFalse(b bool, message string) {
	if b {
		panic(errors.New(message))
	}
}
