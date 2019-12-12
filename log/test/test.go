package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	file, _ := os.OpenFile("go.mod", os.O_RDONLY, 0666)

	fmt.Println(file.Name())

	fmt.Println(time.Now().Format("2006-01-02 15:04:05.999"))
}
