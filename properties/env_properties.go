package properties

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	empty      = ""
	trueValue  = "true"
	falseValue = "false"
)

func GetBoolOrDefaultFromEnv(name string, defaultValue bool) bool {
	value := os.Getenv(name)

	if empty == value {
		return defaultValue
	}

	if trueValue == value {
		return true
	} else if falseValue == value {
		return false
	} else {
		return defaultValue
	}
}

func GetBoolFromEnv(name string) (bool, error) {
	value := os.Getenv(name)

	if empty == value {
		return false, errors.New(fmt.Sprintf("empty value of env '%s'", name))
	}

	if trueValue == value {
		return true, nil
	} else if falseValue == value {
		return false, nil
	} else {
		return false, errors.New(fmt.Sprintf("invalid value '%s' of env '%s'", value, name))
	}
}

func GetStringOrDefaultFromEnv(name string, defaultValue string) string {
	value := os.Getenv(name)

	if empty == value {
		return defaultValue
	}

	return value
}

func GetStringFromEnv(name string) (string, error) {
	value := os.Getenv(name)

	if empty == value {
		return "", errors.New(fmt.Sprintf("empty value of env '%s'", name))
	}

	return value, nil
}

func GetIntOrDefaultFromEnv(name string, defaultValue int) int {
	return int(getSignedNumWithDefault(name, int64(defaultValue)))
}

func GetIntFromEnv(name string) (int, error) {
	value, err := getSignedNumWithoutDefault(name)

	if err != nil {
		return 0, err
	}

	return int(value), nil
}

func GetInt8OrDefaultFromEnv(name string, defaultValue int8) int8 {
	return int8(getSignedNumWithDefault(name, int64(defaultValue)))
}

func GetInt8FromEnv(name string) (int8, error) {
	value, err := getSignedNumWithoutDefault(name)

	if err != nil {
		return int8(0), err
	}

	return int8(value), nil
}

func GetInt16OrDefaultFromEnv(name string, defaultValue int16) int16 {
	return int16(getSignedNumWithDefault(name, int64(defaultValue)))
}

func GetInt16FromEnv(name string) (int16, error) {
	value, err := getSignedNumWithoutDefault(name)

	if err != nil {
		return int16(0), err
	}

	return int16(value), nil
}

func GetInt32OrDefaultFromEnv(name string, defaultValue int32) int32 {
	return int32(getSignedNumWithDefault(name, int64(defaultValue)))
}

func GetInt32FromEnv(name string) (int32, error) {
	value, err := getSignedNumWithoutDefault(name)

	if err != nil {
		return int32(0), err
	}

	return int32(value), nil
}

func GetInt64OrDefaultFromEnv(name string, defaultValue int64) int64 {
	return getSignedNumWithDefault(name, defaultValue)
}

func GetInt64FromEnv(name string) (int64, error) {
	value, err := getSignedNumWithoutDefault(name)

	if err != nil {
		return int64(0), err
	}

	return value, nil
}

func GetUintOrDefaultFromEnv(name string, defaultValue int) uint {
	return uint(getUnsignedNumWithDefault(name, uint64(defaultValue)))
}

func GetUintFromEnv(name string) (uint, error) {
	value, err := getUnsignedNumWithoutDefault(name)

	if err != nil {
		return uint(0), err
	}

	return uint(value), nil
}

func GetUint8OrDefaultFromEnv(name string, defaultValue uint8) uint8 {
	return uint8(getUnsignedNumWithDefault(name, uint64(defaultValue)))
}

func GetUint8FromEnv(name string) (uint8, error) {
	value, err := getUnsignedNumWithoutDefault(name)

	if err != nil {
		return uint8(0), err
	}

	return uint8(value), nil
}

func GetUint16OrDefaultFromEnv(name string, defaultValue uint16) uint16 {
	return uint16(getUnsignedNumWithDefault(name, uint64(defaultValue)))
}

func GetUint16FromEnv(name string) (uint16, error) {
	value, err := getUnsignedNumWithoutDefault(name)

	if err != nil {
		return uint16(0), err
	}

	return uint16(value), nil
}

func GetUint32OrDefaultFromEnv(name string, defaultValue uint32) uint32 {
	return uint32(getUnsignedNumWithDefault(name, uint64(defaultValue)))
}

func GetUint32FromEnv(name string) (uint32, error) {
	value, err := getUnsignedNumWithoutDefault(name)

	if err != nil {
		return uint32(0), err
	}

	return uint32(value), nil
}

func GetUint64OrDefaultFromEnv(name string, defaultValue uint64) uint64 {
	return getUnsignedNumWithDefault(name, defaultValue)
}

func GetUint64FromEnv(name string) (uint64, error) {
	value, err := getUnsignedNumWithoutDefault(name)

	if err != nil {
		return uint64(0), err
	}

	return value, nil
}

func getSignedNumWithDefault(name string, defaultValue int64) int64 {
	value := os.Getenv(name)

	if empty == value {
		return defaultValue
	}

	var intValue int64
	var err error

	if strings.HasPrefix(value, "0x") {
		intValue, err = strconv.ParseInt(value[2:], 16, 64)
	} else if strings.HasPrefix(value, "0") {
		intValue, err = strconv.ParseInt(value[1:], 8, 64)
	} else {
		intValue, err = strconv.ParseInt(value, 10, 64)
	}

	if err != nil {
		return defaultValue
	}

	return intValue
}

func getUnsignedNumWithDefault(name string, defaultValue uint64) uint64 {
	value := os.Getenv(name)

	if empty == value {
		return defaultValue
	}

	var uintValue uint64
	var err error

	if strings.HasPrefix(value, "0x") {
		uintValue, err = strconv.ParseUint(value[2:], 16, 64)
	} else if strings.HasPrefix(value, "0") {
		uintValue, err = strconv.ParseUint(value[1:], 8, 64)
	} else {
		uintValue, err = strconv.ParseUint(value, 10, 64)
	}

	if err != nil {
		return defaultValue
	}

	return uintValue
}

func getSignedNumWithoutDefault(name string) (int64, error) {
	value := os.Getenv(name)

	if empty == value {
		return 0, errors.New(fmt.Sprintf("empty value of env '%s'", name))
	}

	var intValue int64
	var err error

	if strings.HasPrefix(value, "0x") {
		intValue, err = strconv.ParseInt(value[2:], 16, 64)
	} else if strings.HasPrefix(value, "0") {
		intValue, err = strconv.ParseInt(value[1:], 8, 64)
	} else {
		intValue, err = strconv.ParseInt(value, 10, 64)
	}

	if err != nil {
		return 0, err
	}

	return intValue, nil
}

func getUnsignedNumWithoutDefault(name string) (uint64, error) {
	value := os.Getenv(name)

	if empty == value {
		return 0, errors.New(fmt.Sprintf("empty value of env '%s'", name))
	}

	var uintValue uint64
	var err error

	if strings.HasPrefix(value, "0x") {
		uintValue, err = strconv.ParseUint(value[2:], 16, 64)
	} else if strings.HasPrefix(value, "0") {
		uintValue, err = strconv.ParseUint(value[1:], 8, 64)
	} else {
		uintValue, err = strconv.ParseUint(value, 10, 64)
	}

	if err != nil {
		return 0, err
	}

	return uintValue, nil
}
