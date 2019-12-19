package properties

import (
	"fmt"
	"github.com/liuyehcf/common-gtools/assert"
	"os"
	"testing"
	"time"
)

var name string

func TestBoolWithDefault(t *testing.T) {
	var value bool

	resetName()

	_ = os.Setenv(name, "")
	value = GetBoolOrDefaultFromEnv(name, true)
	assert.AssertTrue(value, "test")
	value = GetBoolOrDefaultFromEnv(name, false)
	assert.AssertFalse(value, "test")

	_ = os.Setenv(name, "wrong")
	value = GetBoolOrDefaultFromEnv(name, true)
	assert.AssertTrue(value, "test")
	value = GetBoolOrDefaultFromEnv(name, false)
	assert.AssertFalse(value, "test")

	_ = os.Setenv(name, "true")
	value = GetBoolOrDefaultFromEnv(name, true)
	assert.AssertTrue(value, "test")
	value = GetBoolOrDefaultFromEnv(name, false)
	assert.AssertTrue(value, "test")

	_ = os.Setenv(name, "false")
	value = GetBoolOrDefaultFromEnv(name, true)
	assert.AssertFalse(value, "test")
	value = GetBoolOrDefaultFromEnv(name, false)
	assert.AssertFalse(value, "test")
}

func TestBoolWithoutDefault(t *testing.T) {
	var value bool
	var err error

	resetName()

	_ = os.Setenv(name, "")
	value, err = GetBoolFromEnv(name)
	assert.AssertNotNil(err, "test")

	_ = os.Setenv(name, "wrong")
	value, err = GetBoolFromEnv(name)
	assert.AssertNotNil(err, "test")

	_ = os.Setenv(name, "true")
	value, err = GetBoolFromEnv(name)
	assert.AssertNil(err, "test")
	assert.AssertTrue(value, "test")

	_ = os.Setenv(name, "false")
	value, err = GetBoolFromEnv(name)
	assert.AssertNil(err, "test")
	assert.AssertFalse(value, "test")
}

func TestStringWithDefault(t *testing.T) {
	var value string

	resetName()

	_ = os.Setenv(name, "")
	value = GetStringOrDefaultFromEnv(name, "something")
	assert.AssertTrue(value == "something", "test")

	_ = os.Setenv(name, "true")
	value = GetStringOrDefaultFromEnv(name, "something")
	assert.AssertTrue(value == "true", "test")
}

func TestStringWithoutDefault(t *testing.T) {
	var value string
	var err error

	resetName()

	_ = os.Setenv(name, "")
	value, err = GetStringFromEnv(name)
	assert.AssertNotNil(err, "test")

	_ = os.Setenv(name, "true")
	value, err = GetStringFromEnv(name)
	assert.AssertNil(err, "test")
	assert.AssertTrue(value == "true", "test")
}

func TestIntWithDefault(t *testing.T) {
	var value int

	resetName()

	_ = os.Setenv(name, "")
	value = GetIntOrDefaultFromEnv(name, 6)
	assert.AssertTrue(6 == value, "test")

	_ = os.Setenv(name, "abc")
	value = GetIntOrDefaultFromEnv(name, 6)
	assert.AssertTrue(6 == value, "test")

	_ = os.Setenv(name, "3")
	value = GetIntOrDefaultFromEnv(name, 6)
	assert.AssertTrue(3 == value, "test")

	_ = os.Setenv(name, "0777")
	value = GetIntOrDefaultFromEnv(name, 6)
	assert.AssertTrue(511 == value, "test")

	_ = os.Setenv(name, "0xff")
	value = GetIntOrDefaultFromEnv(name, 6)
	assert.AssertTrue(255 == value, "test")
}

func TestIntWithoutDefault(t *testing.T) {
	var value int
	var err error

	resetName()

	_ = os.Setenv(name, "")
	value, err = GetIntFromEnv(name)
	assert.AssertNotNil(err, "test")

	_ = os.Setenv(name, "abc")
	value, err = GetIntFromEnv(name)
	assert.AssertNotNil(err, "test")

	_ = os.Setenv(name, "3")
	value, err = GetIntFromEnv(name)
	assert.AssertNil(err, "test")
	assert.AssertTrue(3 == value, "test")

	_ = os.Setenv(name, "0777")
	value, err = GetIntFromEnv(name)
	assert.AssertNil(err, "test")
	assert.AssertTrue(511 == value, "test")

	_ = os.Setenv(name, "0xff")
	value, err = GetIntFromEnv(name)
	assert.AssertNil(err, "test")
	assert.AssertTrue(255 == value, "test")
}

func TestInt8WithDefault(t *testing.T) {
	var value int8

	resetName()

	_ = os.Setenv(name, "")
	value = GetInt8OrDefaultFromEnv(name, 6)
	assert.AssertTrue(6 == value, "test")

	_ = os.Setenv(name, "abc")
	value = GetInt8OrDefaultFromEnv(name, 6)
	assert.AssertTrue(6 == value, "test")

	_ = os.Setenv(name, "3")
	value = GetInt8OrDefaultFromEnv(name, 6)
	assert.AssertTrue(3 == value, "test")

	_ = os.Setenv(name, "077")
	value = GetInt8OrDefaultFromEnv(name, 6)
	assert.AssertTrue(63 == value, "test")

	_ = os.Setenv(name, "0xf")
	value = GetInt8OrDefaultFromEnv(name, 6)
	assert.AssertTrue(15 == value, "test")
}

func TestInt8WithoutDefault(t *testing.T) {
	var value int8
	var err error

	resetName()

	_ = os.Setenv(name, "")
	value, err = GetInt8FromEnv(name)
	assert.AssertNotNil(err, "test")

	_ = os.Setenv(name, "abc")
	value, err = GetInt8FromEnv(name)
	assert.AssertNotNil(err, "test")

	_ = os.Setenv(name, "3")
	value, err = GetInt8FromEnv(name)
	assert.AssertNil(err, "test")
	assert.AssertTrue(3 == value, "test")

	_ = os.Setenv(name, "077")
	value, err = GetInt8FromEnv(name)
	assert.AssertNil(err, "test")
	assert.AssertTrue(63 == value, "test")

	_ = os.Setenv(name, "0xf")
	value, err = GetInt8FromEnv(name)
	assert.AssertNil(err, "test")
	assert.AssertTrue(15 == value, "test")
}

func TestInt16WithDefault(t *testing.T) {
	var value int16

	resetName()

	_ = os.Setenv(name, "")
	value = GetInt16OrDefaultFromEnv(name, 6)
	assert.AssertTrue(6 == value, "test")

	_ = os.Setenv(name, "abc")
	value = GetInt16OrDefaultFromEnv(name, 6)
	assert.AssertTrue(6 == value, "test")

	_ = os.Setenv(name, "3")
	value = GetInt16OrDefaultFromEnv(name, 6)
	assert.AssertTrue(3 == value, "test")

	_ = os.Setenv(name, "0777")
	value = GetInt16OrDefaultFromEnv(name, 6)
	assert.AssertTrue(511 == value, "test")

	_ = os.Setenv(name, "0xff")
	value = GetInt16OrDefaultFromEnv(name, 6)
	assert.AssertTrue(255 == value, "test")
}

func TestInt16WithoutDefault(t *testing.T) {
	var value int16
	var err error

	resetName()

	_ = os.Setenv(name, "")
	value, err = GetInt16FromEnv(name)
	assert.AssertNotNil(err, "test")

	_ = os.Setenv(name, "abc")
	value, err = GetInt16FromEnv(name)
	assert.AssertNotNil(err, "test")

	_ = os.Setenv(name, "3")
	value, err = GetInt16FromEnv(name)
	assert.AssertNil(err, "test")
	assert.AssertTrue(3 == value, "test")

	_ = os.Setenv(name, "0777")
	value = GetInt16OrDefaultFromEnv(name, 6)
	assert.AssertTrue(511 == value, "test")

	_ = os.Setenv(name, "0xff")
	value = GetInt16OrDefaultFromEnv(name, 6)
	assert.AssertTrue(255 == value, "test")
}

func TestInt32WithDefault(t *testing.T) {
	var value int32

	resetName()

	_ = os.Setenv(name, "")
	value = GetInt32OrDefaultFromEnv(name, 6)
	assert.AssertTrue(6 == value, "test")

	_ = os.Setenv(name, "abc")
	value = GetInt32OrDefaultFromEnv(name, 6)
	assert.AssertTrue(6 == value, "test")

	_ = os.Setenv(name, "3")
	value = GetInt32OrDefaultFromEnv(name, 6)
	assert.AssertTrue(3 == value, "test")

	_ = os.Setenv(name, "0777")
	value = GetInt32OrDefaultFromEnv(name, 6)
	assert.AssertTrue(511 == value, "test")

	_ = os.Setenv(name, "0xff")
	value = GetInt32OrDefaultFromEnv(name, 6)
	assert.AssertTrue(255 == value, "test")
}

func TestInt32WithoutDefault(t *testing.T) {
	var value int32
	var err error

	resetName()

	_ = os.Setenv(name, "")
	value, err = GetInt32FromEnv(name)
	assert.AssertNotNil(err, "test")

	_ = os.Setenv(name, "abc")
	value, err = GetInt32FromEnv(name)
	assert.AssertNotNil(err, "test")

	_ = os.Setenv(name, "3")
	value, err = GetInt32FromEnv(name)
	assert.AssertNil(err, "test")
	assert.AssertTrue(3 == value, "test")

	_ = os.Setenv(name, "0777")
	value = GetInt32OrDefaultFromEnv(name, 6)
	assert.AssertTrue(511 == value, "test")

	_ = os.Setenv(name, "0xff")
	value = GetInt32OrDefaultFromEnv(name, 6)
	assert.AssertTrue(255 == value, "test")
}

func TestInt64WithDefault(t *testing.T) {
	var value int64

	resetName()

	_ = os.Setenv(name, "")
	value = GetInt64OrDefaultFromEnv(name, 6)
	assert.AssertTrue(6 == value, "test")

	_ = os.Setenv(name, "abc")
	value = GetInt64OrDefaultFromEnv(name, 6)
	assert.AssertTrue(6 == value, "test")

	_ = os.Setenv(name, "3")
	value = GetInt64OrDefaultFromEnv(name, 6)
	assert.AssertTrue(3 == value, "test")

	_ = os.Setenv(name, "0777")
	value = GetInt64OrDefaultFromEnv(name, 6)
	assert.AssertTrue(511 == value, "test")

	_ = os.Setenv(name, "0xff")
	value = GetInt64OrDefaultFromEnv(name, 6)
	assert.AssertTrue(255 == value, "test")
}

func TestInt64WithoutDefault(t *testing.T) {
	var value int64
	var err error

	resetName()

	_ = os.Setenv(name, "")
	value, err = GetInt64FromEnv(name)
	assert.AssertNotNil(err, "test")

	_ = os.Setenv(name, "abc")
	value, err = GetInt64FromEnv(name)
	assert.AssertNotNil(err, "test")

	_ = os.Setenv(name, "3")
	value, err = GetInt64FromEnv(name)
	assert.AssertNil(err, "test")
	assert.AssertTrue(3 == value, "test")

	_ = os.Setenv(name, "0777")
	value = GetInt64OrDefaultFromEnv(name, 6)
	assert.AssertTrue(511 == value, "test")

	_ = os.Setenv(name, "0xff")
	value = GetInt64OrDefaultFromEnv(name, 6)
	assert.AssertTrue(255 == value, "test")
}

func TestUintWithDefault(t *testing.T) {
	var value uint

	resetName()

	_ = os.Setenv(name, "")
	value = GetUintOrDefaultFromEnv(name, 6)
	assert.AssertTrue(6 == value, "test")

	_ = os.Setenv(name, "abc")
	value = GetUintOrDefaultFromEnv(name, 6)
	assert.AssertTrue(6 == value, "test")

	_ = os.Setenv(name, "3")
	value = GetUintOrDefaultFromEnv(name, 6)
	assert.AssertTrue(3 == value, "test")

	_ = os.Setenv(name, "0777")
	value = GetUintOrDefaultFromEnv(name, 6)
	assert.AssertTrue(511 == value, "test")

	_ = os.Setenv(name, "0xff")
	value = GetUintOrDefaultFromEnv(name, 6)
	assert.AssertTrue(255 == value, "test")
}

func TestUintWithoutDefault(t *testing.T) {
	var value uint
	var err error

	resetName()

	_ = os.Setenv(name, "")
	value, err = GetUintFromEnv(name)
	assert.AssertNotNil(err, "test")

	_ = os.Setenv(name, "abc")
	value, err = GetUintFromEnv(name)
	assert.AssertNotNil(err, "test")

	_ = os.Setenv(name, "3")
	value, err = GetUintFromEnv(name)
	assert.AssertNil(err, "test")
	assert.AssertTrue(3 == value, "test")

	_ = os.Setenv(name, "0777")
	value = GetUintOrDefaultFromEnv(name, 6)
	assert.AssertTrue(511 == value, "test")

	_ = os.Setenv(name, "0xff")
	value = GetUintOrDefaultFromEnv(name, 6)
	assert.AssertTrue(255 == value, "test")
}

func TestUint8WithDefault(t *testing.T) {
	var value uint8

	resetName()

	_ = os.Setenv(name, "")
	value = GetUint8OrDefaultFromEnv(name, 6)
	assert.AssertTrue(6 == value, "test")

	_ = os.Setenv(name, "abc")
	value = GetUint8OrDefaultFromEnv(name, 6)
	assert.AssertTrue(6 == value, "test")

	_ = os.Setenv(name, "3")
	value = GetUint8OrDefaultFromEnv(name, 6)
	assert.AssertTrue(3 == value, "test")

	_ = os.Setenv(name, "077")
	value = GetUint8OrDefaultFromEnv(name, 6)
	assert.AssertTrue(63 == value, "test")

	_ = os.Setenv(name, "0xff")
	value = GetUint8OrDefaultFromEnv(name, 6)
	assert.AssertTrue(255 == value, "test")
}

func TestUint8WithoutDefault(t *testing.T) {
	var value uint8
	var err error

	resetName()

	_ = os.Setenv(name, "")
	value, err = GetUint8FromEnv(name)
	assert.AssertNotNil(err, "test")

	_ = os.Setenv(name, "abc")
	value, err = GetUint8FromEnv(name)
	assert.AssertNotNil(err, "test")

	_ = os.Setenv(name, "3")
	value, err = GetUint8FromEnv(name)
	assert.AssertNil(err, "test")
	assert.AssertTrue(3 == value, "test")

	_ = os.Setenv(name, "077")
	value = GetUint8OrDefaultFromEnv(name, 6)
	assert.AssertTrue(63 == value, "test")

	_ = os.Setenv(name, "0xff")
	value = GetUint8OrDefaultFromEnv(name, 6)
	assert.AssertTrue(255 == value, "test")
}

func TestUint16WithDefault(t *testing.T) {
	var value uint16

	resetName()

	_ = os.Setenv(name, "")
	value = GetUint16OrDefaultFromEnv(name, 6)
	assert.AssertTrue(6 == value, "test")

	_ = os.Setenv(name, "abc")
	value = GetUint16OrDefaultFromEnv(name, 6)
	assert.AssertTrue(6 == value, "test")

	_ = os.Setenv(name, "3")
	value = GetUint16OrDefaultFromEnv(name, 6)
	assert.AssertTrue(3 == value, "test")

	_ = os.Setenv(name, "0777")
	value = GetUint16OrDefaultFromEnv(name, 6)
	assert.AssertTrue(511 == value, "test")

	_ = os.Setenv(name, "0xff")
	value = GetUint16OrDefaultFromEnv(name, 6)
	assert.AssertTrue(255 == value, "test")
}

func TestUint16WithoutDefault(t *testing.T) {
	var value uint16
	var err error

	resetName()

	_ = os.Setenv(name, "")
	value, err = GetUint16FromEnv(name)
	assert.AssertNotNil(err, "test")

	_ = os.Setenv(name, "abc")
	value, err = GetUint16FromEnv(name)
	assert.AssertNotNil(err, "test")

	_ = os.Setenv(name, "3")
	value, err = GetUint16FromEnv(name)
	assert.AssertNil(err, "test")
	assert.AssertTrue(3 == value, "test")

	_ = os.Setenv(name, "0777")
	value = GetUint16OrDefaultFromEnv(name, 6)
	assert.AssertTrue(511 == value, "test")

	_ = os.Setenv(name, "0xff")
	value = GetUint16OrDefaultFromEnv(name, 6)
	assert.AssertTrue(255 == value, "test")
}

func TestUint32WithDefault(t *testing.T) {
	var value uint32

	resetName()

	_ = os.Setenv(name, "")
	value = GetUint32OrDefaultFromEnv(name, 6)
	assert.AssertTrue(6 == value, "test")

	_ = os.Setenv(name, "abc")
	value = GetUint32OrDefaultFromEnv(name, 6)
	assert.AssertTrue(6 == value, "test")

	_ = os.Setenv(name, "3")
	value = GetUint32OrDefaultFromEnv(name, 6)
	assert.AssertTrue(3 == value, "test")

	_ = os.Setenv(name, "0777")
	value = GetUint32OrDefaultFromEnv(name, 6)
	assert.AssertTrue(511 == value, "test")

	_ = os.Setenv(name, "0xff")
	value = GetUint32OrDefaultFromEnv(name, 6)
	assert.AssertTrue(255 == value, "test")
}

func TestUint32WithoutDefault(t *testing.T) {
	var value uint32
	var err error

	resetName()

	_ = os.Setenv(name, "")
	value, err = GetUint32FromEnv(name)
	assert.AssertNotNil(err, "test")

	_ = os.Setenv(name, "abc")
	value, err = GetUint32FromEnv(name)
	assert.AssertNotNil(err, "test")

	_ = os.Setenv(name, "3")
	value, err = GetUint32FromEnv(name)
	assert.AssertNil(err, "test")
	assert.AssertTrue(3 == value, "test")

	_ = os.Setenv(name, "0777")
	value = GetUint32OrDefaultFromEnv(name, 6)
	assert.AssertTrue(511 == value, "test")

	_ = os.Setenv(name, "0xff")
	value = GetUint32OrDefaultFromEnv(name, 6)
	assert.AssertTrue(255 == value, "test")
}

func TestUint64WithDefault(t *testing.T) {
	var value uint64

	resetName()

	_ = os.Setenv(name, "")
	value = GetUint64OrDefaultFromEnv(name, 6)
	assert.AssertTrue(6 == value, "test")

	_ = os.Setenv(name, "abc")
	value = GetUint64OrDefaultFromEnv(name, 6)
	assert.AssertTrue(6 == value, "test")

	_ = os.Setenv(name, "3")
	value = GetUint64OrDefaultFromEnv(name, 6)
	assert.AssertTrue(3 == value, "test")

	_ = os.Setenv(name, "0777")
	value = GetUint64OrDefaultFromEnv(name, 6)
	assert.AssertTrue(511 == value, "test")

	_ = os.Setenv(name, "0xff")
	value = GetUint64OrDefaultFromEnv(name, 6)
	assert.AssertTrue(255 == value, "test")
}

func TestUint64WithoutDefault(t *testing.T) {
	var value uint64
	var err error

	resetName()

	_ = os.Setenv(name, "")
	value, err = GetUint64FromEnv(name)
	assert.AssertNotNil(err, "test")

	_ = os.Setenv(name, "abc")
	value, err = GetUint64FromEnv(name)
	assert.AssertNotNil(err, "test")

	_ = os.Setenv(name, "3")
	value, err = GetUint64FromEnv(name)
	assert.AssertNil(err, "test")
	assert.AssertTrue(3 == value, "test")

	_ = os.Setenv(name, "0777")
	value = GetUint64OrDefaultFromEnv(name, 6)
	assert.AssertTrue(511 == value, "test")

	_ = os.Setenv(name, "0xff")
	value = GetUint64OrDefaultFromEnv(name, 6)
	assert.AssertTrue(255 == value, "test")
}

func resetName() {
	name = fmt.Sprintf("testName_%d", time.Now().Second())
}
