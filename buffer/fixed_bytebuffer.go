package buffer

import (
	"fmt"
	"github.com/liuyehcf/common-gtools/assert"
)

type FixedByteBuffer struct {
	// target byte slice
	mem []byte

	// capacity of byte buffer
	capacity int

	// how many bytes can be read
	readableBytes int

	// first byte can be read
	readIndex int

	// next index to be written
	writeIndex int

	// marked readableBytes
	markReadableBytes int

	// marked readIndex
	markReadIndex int

	// marked writeIndex
	markWriteIndex int
}

func (buffer *FixedByteBuffer) Write(src []byte) {
	srcLen := len(src)

	remainSpace := buffer.capacity - buffer.readableBytes
	assert.AssertFalse(srcLen > remainSpace, fmt.Sprintf("no enough space, required=%d, remain=%d", srcLen, remainSpace))

	var actualWritableLen int

	if remainSpace <= srcLen {
		actualWritableLen = remainSpace
	} else {
		actualWritableLen = srcLen
	}

	tailLen := buffer.capacity - buffer.writeIndex

	// in this condition, we can derive writeIndex <= readIndex
	// so all bytes from writeIndex to tail can be written
	if tailLen <= actualWritableLen {
		tailWritableLen := tailLen
		headWritableLen := actualWritableLen - tailWritableLen

		// write all tail spaces
		copy(buffer.mem[buffer.writeIndex:], src[:tailWritableLen])

		// write remain read spaces
		copy(buffer.mem[:headWritableLen], src[tailWritableLen:actualWritableLen])

		buffer.writeIndex = headWritableLen
	} else {
		copy(buffer.mem[buffer.writeIndex:buffer.writeIndex+actualWritableLen], src[:actualWritableLen])
		buffer.writeIndex = buffer.writeIndex + actualWritableLen
	}

	buffer.readableBytes += actualWritableLen
}

func (buffer *FixedByteBuffer) Read(dst []byte) int {
	if buffer.readableBytes <= 0 {
		return 0
	}

	dstLen := len(dst)
	var actualReadLen int

	if buffer.readableBytes <= dstLen {
		actualReadLen = buffer.readableBytes
	} else {
		actualReadLen = dstLen
	}

	tailLen := buffer.capacity - buffer.readIndex

	// in this condition, we can derive writeIndex <= readIndex
	// so all bytes from readIndex to tail can be read
	if tailLen <= actualReadLen {
		tailReadableLen := tailLen
		headReadableLen := actualReadLen - tailReadableLen

		// read all tail bytes
		copy(dst[:tailReadableLen], buffer.mem[buffer.readIndex:])

		// read remain head bytes
		copy(dst[tailReadableLen:actualReadLen], buffer.mem[:headReadableLen])

		buffer.readIndex = headReadableLen
	} else {
		copy(dst[:actualReadLen], buffer.mem[buffer.readIndex:buffer.readIndex+actualReadLen])
		buffer.readIndex = buffer.readIndex + actualReadLen
	}

	buffer.readableBytes -= actualReadLen

	return actualReadLen
}

func (buffer *FixedByteBuffer) Capacity() int {
	return buffer.capacity
}

func (buffer *FixedByteBuffer) ReadableBytes() int {
	return buffer.readableBytes
}

func (buffer *FixedByteBuffer) ReadIndex() int {
	return buffer.readIndex
}

func (buffer *FixedByteBuffer) WriteIndex() int {
	return buffer.writeIndex
}

func (buffer *FixedByteBuffer) Mark() {
	buffer.markReadableBytes = buffer.readableBytes
	buffer.markReadIndex = buffer.readIndex
	buffer.markWriteIndex = buffer.writeIndex
}

func (buffer *FixedByteBuffer) Recover() {
	buffer.readableBytes = buffer.markReadableBytes
	buffer.readIndex = buffer.markReadIndex
	buffer.writeIndex = buffer.markWriteIndex
}

func (buffer *FixedByteBuffer) Clean() {
	buffer.readableBytes = 0
	buffer.readIndex = 0
	buffer.writeIndex = 0
}

func NewFixedByteBuffer(size int) Buffer {
	return &FixedByteBuffer{
		mem:           make([]byte, size),
		capacity:      size,
		readableBytes: 0,
		readIndex:     0,
		writeIndex:    0,
	}
}
