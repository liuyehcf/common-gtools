package main

import (
	"github.com/liuyehcf/common-gtools/assert"
	buf "github.com/liuyehcf/common-gtools/buffer"
	"math/rand"
)

func main() {
	case1()
	case2()
	case3()
}

func case1() {
	buffer := buf.NewFixedByteBuffer(5)
	assert.AssertTrue(buffer.ReadableBytes() == 0, "test")
	assert.AssertTrue(buffer.ReadIndex() == 0, "test")
	assert.AssertTrue(buffer.WriteIndex() == 0, "test")

	buffer.Write([]byte{1, 2})
	assert.AssertTrue(buffer.ReadableBytes() == 2, "test")
	assert.AssertTrue(buffer.ReadIndex() == 0, "test")
	assert.AssertTrue(buffer.WriteIndex() == 2, "test")

	bytes := make([]byte, 5)

	n := buffer.Read(bytes)
	assert.AssertTrue(n == 2, "test")
	assert.AssertTrue(bytes[0] == 1, "test")
	assert.AssertTrue(bytes[1] == 2, "test")
	assert.AssertTrue(buffer.ReadableBytes() == 0, "test")
	assert.AssertTrue(buffer.ReadIndex() == 2, "test")
	assert.AssertTrue(buffer.WriteIndex() == 2, "test")

	buffer.Write([]byte{1, 2, 3})
	assert.AssertTrue(buffer.ReadableBytes() == 3, "test")
	assert.AssertTrue(buffer.ReadIndex() == 2, "test")
	assert.AssertTrue(buffer.WriteIndex() == 0, "test")

	buffer.Write([]byte{4, 5})
	assert.AssertTrue(buffer.ReadableBytes() == 5, "test")
	assert.AssertTrue(buffer.ReadIndex() == 2, "test")
	assert.AssertTrue(buffer.WriteIndex() == 2, "test")

	n = buffer.Read(bytes)
	assert.AssertTrue(n == 5, "test")
	assert.AssertTrue(bytes[0] == 1, "test")
	assert.AssertTrue(bytes[1] == 2, "test")
	assert.AssertTrue(bytes[2] == 3, "test")
	assert.AssertTrue(bytes[3] == 4, "test")
	assert.AssertTrue(bytes[4] == 5, "test")
	assert.AssertTrue(buffer.ReadableBytes() == 0, "test")
	assert.AssertTrue(buffer.ReadIndex() == 2, "test")
	assert.AssertTrue(buffer.WriteIndex() == 2, "test")
}

func case2() {
	buffer := buf.NewFixedByteBuffer(100)

	for i := 0; i < 100; i += 1 {
		bytes := make([]byte, i)
		for j := 0; j < i; j += 1 {
			bytes[j] = byte(rand.Int())
		}

		buffer.Write(bytes)
		assert.AssertTrue(buffer.ReadableBytes() == i, "test")

		toBytes := make([]byte, i)
		read := buffer.Read(toBytes)
		assert.AssertTrue(read == i, "test")

		for j := 0; j < i; j += 1 {
			assert.AssertTrue(bytes[j] == toBytes[j], "test")
		}
	}
}

func case3() {
	buffer := buf.NewFixedByteBuffer(6)
	buffer.Write([]byte{1, 2, 3, 4, 5, 6})

	assert.AssertTrue(buffer.ReadableBytes() == 6, "test")
	assert.AssertTrue(buffer.ReadIndex() == 0, "test")
	assert.AssertTrue(buffer.WriteIndex() == 0, "test")

	bytes := make([]byte, 6)
	buffer.Mark()
	n := buffer.Read(bytes)
	assert.AssertTrue(n == 6, "test")
	assert.AssertTrue(buffer.ReadableBytes() == 0, "test")
	assert.AssertTrue(buffer.ReadIndex() == 0, "test")
	assert.AssertTrue(buffer.WriteIndex() == 0, "test")

	buffer.Recover()
	assert.AssertTrue(buffer.ReadableBytes() == 6, "test")
	assert.AssertTrue(buffer.ReadIndex() == 0, "test")
	assert.AssertTrue(buffer.WriteIndex() == 0, "test")

	buffer.Read(bytes)
	assert.AssertTrue(n == 6, "test")
	assert.AssertTrue(buffer.ReadableBytes() == 0, "test")
	assert.AssertTrue(buffer.ReadIndex() == 0, "test")
	assert.AssertTrue(buffer.WriteIndex() == 0, "test")
}
