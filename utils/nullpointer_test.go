package utils

import (
	"testing"
)

type testInterface interface {
	DoSomething()
}

type testImpl struct {
}

func (t *testImpl) DoSomething() {
	panic("implement me")
}

func TestIsInterfaceNil(t *testing.T) {
	var nilImpl *testImpl = nil

	AssertTrue(IsNil(nilImpl), "test")
	AssertFalse(IsNotNil(nilImpl), "test")

	var nilInterface testInterface = nil
	AssertTrue(IsNil(nilInterface), "test")
	AssertFalse(IsNotNil(nilInterface), "test")
}

func TestIsInterfaceNotNil(t *testing.T) {
	var nilImpl = new(testImpl)

	AssertTrue(IsNotNil(nilImpl), "test")
	AssertFalse(IsNil(nilImpl), "test")

	var nilInterface testInterface = nilImpl
	AssertTrue(IsNotNil(nilInterface), "test")
	AssertFalse(IsNil(nilInterface), "test")
}

func TestBasicType(t *testing.T) {
	AssertTrue(IsNotNil(true), "test")
	AssertFalse(IsNil(true), "test")

	AssertTrue(IsNotNil(int(1)), "test")
	AssertFalse(IsNil(int(1)), "test")

	AssertTrue(IsNotNil(int8(1)), "test")
	AssertFalse(IsNil(int8(1)), "test")

	AssertTrue(IsNotNil(int16(1)), "test")
	AssertFalse(IsNil(int16(1)), "test")

	AssertTrue(IsNotNil(int32(1)), "test")
	AssertFalse(IsNil(int32(1)), "test")

	AssertTrue(IsNotNil(int64(1)), "test")
	AssertFalse(IsNil(int64(1)), "test")

	AssertTrue(IsNotNil(uint(1)), "test")
	AssertFalse(IsNil(uint(1)), "test")

	AssertTrue(IsNotNil(uint8(1)), "test")
	AssertFalse(IsNil(uint8(1)), "test")

	AssertTrue(IsNotNil(uint16(1)), "test")
	AssertFalse(IsNil(uint16(1)), "test")

	AssertTrue(IsNotNil(uint32(1)), "test")
	AssertFalse(IsNil(uint32(1)), "test")

	AssertTrue(IsNotNil(uint64(1)), "test")
	AssertFalse(IsNil(uint64(1)), "test")

	AssertTrue(IsNotNil(float32(1)), "test")
	AssertFalse(IsNil(float32(1)), "test")

	AssertTrue(IsNotNil(float64(1)), "test")
	AssertFalse(IsNil(float64(1)), "test")

	AssertTrue(IsNotNil(complex64(1)), "test")
	AssertFalse(IsNil(complex64(1)), "test")

	AssertTrue(IsNotNil(complex128(1)), "test")
	AssertFalse(IsNil(complex128(1)), "test")

	AssertTrue(IsNotNil([1]testInterface{nil}), "test")
	AssertFalse(IsNil([1]testInterface{nil}), "test")

	AssertTrue(IsNotNil([1]*testImpl{new(testImpl)}), "test")
	AssertFalse(IsNil([1]*testImpl{new(testImpl)}), "test")
}
