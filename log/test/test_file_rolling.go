package main

import (
	"github.com/liuyehcf/common-gtools/log"
	"time"
)

func main() {
	layout := "%d{2006-01-02 15:04:05.999} [%p] %m%n"

	commonFileAppender := log.NewFileAppender(&log.AppenderConfig{
		Layout:  layout,
		Filters: nil,
		FileRollingPolicy: &log.RollingPolicy{
			Directory:       "/tmp/gtools/logs",
			FileName:        "rolling",
			TimeGranularity: log.TimeGranularityHour,
			MaxHistory:      20,
			MaxFileSize:     1,
		},
	})

	logger := log.NewLogger(log.Root, log.InfoLevel, false, []log.Appender{commonFileAppender})

	go func() {
		for {
			logger.Info("now: '{}'", time.Now())

			time.Sleep(time.Millisecond)
		}
	}()

	<-make(chan interface{}, 0)
}
