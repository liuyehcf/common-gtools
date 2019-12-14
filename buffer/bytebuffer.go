package buffer

type ByteBuffer interface {
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
