package main

import (
	"github.com/liuyehcf/common-gtools/log"
	"os"
	"time"
)

func main() {
	stdoutAppender := log.NewWriterAppender(&log.AppenderConfig{
		Layout:  "%d{2006-01-02 15:04:05.999} [%p]-[%c]-[%L] --- %m%n",
		Filters: nil,
		Writer:  os.Stdout,
	})

	additivityLogger := log.NewLogger("additivityLogger", log.InfoLevel, true, []log.Appender{stdoutAppender})
	additivityLogger.Info("you can see this twice")

	nonAdditivityLogger := log.NewLogger("nonAdditivityLogger", log.InfoLevel, false, []log.Appender{stdoutAppender})
	nonAdditivityLogger.Info("you can see this once")

	time.Sleep(time.Second)
}
