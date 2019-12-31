package utils

import (
	"github.com/liuyehcf/common-gtools/assert"
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

	assert.AssertTrue(IsNil(nilImpl), "test")
	assert.AssertFalse(IsNotNil(nilImpl), "test")

	var nilInterface testInterface = nil
	assert.AssertTrue(IsNil(nilInterface), "test")
	assert.AssertFalse(IsNotNil(nilInterface), "test")
}

func TestIsInterfaceNotNil(t *testing.T) {
	var nilImpl = new(testImpl)

	assert.AssertTrue(IsNotNil(nilImpl), "test")
	assert.AssertFalse(IsNil(nilImpl), "test")

	var nilInterface testInterface = nilImpl
	assert.AssertTrue(IsNotNil(nilInterface), "test")
	assert.AssertFalse(IsNil(nilInterface), "test")
}

func TestBasicType(t *testing.T) {
	assert.AssertTrue(IsNotNil(true), "test")
	assert.AssertFalse(IsNil(true), "test")
}
