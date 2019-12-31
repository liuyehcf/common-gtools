package main

import (
	buf "github.com/liuyehcf/common-gtools/buffer"
	"github.com/liuyehcf/common-gtools/utils"
	"math/rand"
	"testing"
)

func TestCase1(t *testing.T) {
	buffer := buf.NewRecycleByteBuffer(5)
	utils.AssertTrue(buffer.ReadableBytes() == 0, "test")
	utils.AssertTrue(buffer.ReadIndex() == 0, "test")
	utils.AssertTrue(buffer.WriteIndex() == 0, "test")

	buffer.Write([]byte{1, 2})
	utils.AssertTrue(buffer.ReadableBytes() == 2, "test")
	utils.AssertTrue(buffer.ReadIndex() == 0, "test")
	utils.AssertTrue(buffer.WriteIndex() == 2, "test")

	bytes := make([]byte, 5)

	n := buffer.Read(bytes)
	utils.AssertTrue(n == 2, "test")
	utils.AssertTrue(bytes[0] == 1, "test")
	utils.AssertTrue(bytes[1] == 2, "test")
	utils.AssertTrue(buffer.ReadableBytes() == 0, "test")
	utils.AssertTrue(buffer.ReadIndex() == 2, "test")
	utils.AssertTrue(buffer.WriteIndex() == 2, "test")

	buffer.Write([]byte{1, 2, 3})
	utils.AssertTrue(buffer.ReadableBytes() == 3, "test")
	utils.AssertTrue(buffer.ReadIndex() == 2, "test")
	utils.AssertTrue(buffer.WriteIndex() == 0, "test")

	buffer.Write([]byte{4, 5})
	utils.AssertTrue(buffer.ReadableBytes() == 5, "test")
	utils.AssertTrue(buffer.ReadIndex() == 2, "test")
	utils.AssertTrue(buffer.WriteIndex() == 2, "test")

	n = buffer.Read(bytes)
	utils.AssertTrue(n == 5, "test")
	utils.AssertTrue(bytes[0] == 1, "test")
	utils.AssertTrue(bytes[1] == 2, "test")
	utils.AssertTrue(bytes[2] == 3, "test")
	utils.AssertTrue(bytes[3] == 4, "test")
	utils.AssertTrue(bytes[4] == 5, "test")
	utils.AssertTrue(buffer.ReadableBytes() == 0, "test")
	utils.AssertTrue(buffer.ReadIndex() == 2, "test")
	utils.AssertTrue(buffer.WriteIndex() == 2, "test")
}

func TestCase2(t *testing.T) {
	buffer := buf.NewRecycleByteBuffer(100)

	for i := 0; i < 100; i += 1 {
		bytes := make([]byte, i)
		for j := 0; j < i; j += 1 {
			bytes[j] = byte(rand.Int())
		}

		buffer.Write(bytes)
		utils.AssertTrue(buffer.ReadableBytes() == i, "test")

		toBytes := make([]byte, i)
		read := buffer.Read(toBytes)
		utils.AssertTrue(read == i, "test")

		for j := 0; j < i; j += 1 {
			utils.AssertTrue(bytes[j] == toBytes[j], "test")
		}
	}
}

func TestCase3(t *testing.T) {
	buffer := buf.NewRecycleByteBuffer(6)
	buffer.Write([]byte{1, 2, 3, 4, 5, 6})

	utils.AssertTrue(buffer.ReadableBytes() == 6, "test")
	utils.AssertTrue(buffer.ReadIndex() == 0, "test")
	utils.AssertTrue(buffer.WriteIndex() == 0, "test")

	bytes := make([]byte, 6)
	buffer.Mark()
	n := buffer.Read(bytes)
	utils.AssertTrue(n == 6, "test")
	utils.AssertTrue(buffer.ReadableBytes() == 0, "test")
	utils.AssertTrue(buffer.ReadIndex() == 0, "test")
	utils.AssertTrue(buffer.WriteIndex() == 0, "test")

	buffer.Recover()
	utils.AssertTrue(buffer.ReadableBytes() == 6, "test")
	utils.AssertTrue(buffer.ReadIndex() == 0, "test")
	utils.AssertTrue(buffer.WriteIndex() == 0, "test")

	buffer.Read(bytes)
	utils.AssertTrue(n == 6, "test")
	utils.AssertTrue(buffer.ReadableBytes() == 0, "test")
	utils.AssertTrue(buffer.ReadIndex() == 0, "test")
	utils.AssertTrue(buffer.WriteIndex() == 0, "test")
}
