package main

import (
	"fmt"
	buf "github.com/liuyehcf/common-gtools/buffer"
)

func main() {
	buffer := buf.NewRecycleByteBuffer(10)

	buffer.Write([]byte{1, 2, 3, 4, 5})

	buffer.Mark()

	fmt.Printf("after write, readableBytes=%d\n", buffer.ReadableBytes())
	bytes := make([]byte, 5)
	buffer.Read(bytes)
	fmt.Println(bytes)
	fmt.Printf("after read, readableBytes=%d\n", buffer.ReadableBytes())

	buffer.Recover()
	fmt.Printf("after recover, readableBytes=%d\n", buffer.ReadableBytes())
	bytes = make([]byte, 5)
	buffer.Read(bytes)
	fmt.Println(bytes)
	fmt.Printf("after read, readableBytes=%d\n", buffer.ReadableBytes())

	buffer.Write([]byte{6, 7, 8, 9, 10})
	fmt.Printf("after write, readableBytes=%d\n", buffer.ReadableBytes())
	bytes = make([]byte, 5)
	buffer.Read(bytes)
	fmt.Println(bytes)
	fmt.Printf("after read, readableBytes=%d\n", buffer.ReadableBytes())
}
