package log

import (
	"github.com/liuyehcf/common-gtools/buffer"
)

type StringWriter struct {
	buf buffer.ByteBuffer
}

func NewStringWriter(buf buffer.ByteBuffer) *StringWriter {
	return &StringWriter{buf: buf}
}

func (writer *StringWriter) Write(p []byte) (int, error) {
	writer.buf.Write(p)
	return len(p), nil
}

func (writer *StringWriter) Close() error {
	return nil
}

func (writer *StringWriter) ReadString() string {
	bytes := make([]byte, writer.buf.ReadableBytes())

	writer.buf.Read(bytes)

	return string(bytes)
}
