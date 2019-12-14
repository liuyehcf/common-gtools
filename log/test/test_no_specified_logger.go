package main

import (
	"github.com/liuyehcf/common-gtools/log"
	"time"
)

func main() {
	logger := log.GetLogger("notExist")

	logger.Info("you can see this")

	time.Sleep(time.Second)
}
