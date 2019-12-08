package buffer

import (
	"fmt"
	"github.com/liuyehcf/common-gtools/assert"
)

type Buffer interface {
	// write data from src to buffer
	// panic if no enough space
	Write(src []byte)

	// read data from buffer to dst
	// returns the number of bytes actually read
	Read(dst []byte) int

	// capacity of this buffer
	Capacity() int

	// how many bytes can be read
	ReadableBytes() int

	// read index
	ReadIndex() int

	// write index for next byte
	WriteIndex() int

	// mark current buffer status for subsequent recovery
	Mark()

	// recover buffer status to mark point
	// if the Mark method has not been called before, the behavior is unknown
	Recover()

	// clean status
	Clean()
}

type ByteBuffer struct {
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

func (buffer *ByteBuffer) Write(src []byte) {
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

func (buffer *ByteBuffer) Read(dst []byte) int {
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

func (buffer *ByteBuffer) Capacity() int {
	return buffer.capacity
}

func (buffer *ByteBuffer) ReadableBytes() int {
	return buffer.readableBytes
}

func (buffer *ByteBuffer) ReadIndex() int {
	return buffer.readIndex
}

func (buffer *ByteBuffer) WriteIndex() int {
	return buffer.writeIndex
}

func (buffer *ByteBuffer) Mark() {
	buffer.markReadableBytes = buffer.readableBytes
	buffer.markReadIndex = buffer.readIndex
	buffer.markWriteIndex = buffer.writeIndex
}

func (buffer *ByteBuffer) Recover() {
	buffer.readableBytes = buffer.markReadableBytes
	buffer.readIndex = buffer.markReadIndex
	buffer.writeIndex = buffer.markWriteIndex
}

func (buffer *ByteBuffer) Clean() {
	buffer.readableBytes = 0
	buffer.readIndex = 0
	buffer.writeIndex = 0
}

func NewByteBuffer(size int) Buffer {
	return &ByteBuffer{
		mem:           make([]byte, size),
		capacity:      size,
		readableBytes: 0,
		readIndex:     0,
		writeIndex:    0,
	}
}
