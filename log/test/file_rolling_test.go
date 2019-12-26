package main

import (
	"github.com/liuyehcf/common-gtools/assert"
	"github.com/liuyehcf/common-gtools/log"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestRollingByHour(t *testing.T) {
	history := 20
	rolling(log.TimeGranularityHour, history, func(fileInfo os.FileInfo) {
		name := fileInfo.Name()
		segments := strings.Split(name, ".")
		assert.AssertTrue(len(segments) == 5, "test")

		index, err := strconv.Atoi(segments[len(segments)-2])
		assert.AssertNil(err, "test")
		assert.AssertTrue(0 <= index && index < history, "test")
	})
}

func TestRollingByDay(t *testing.T) {
	history := 20
	rolling(log.TimeGranularityDay, history, func(fileInfo os.FileInfo) {
		name := fileInfo.Name()
		segments := strings.Split(name, ".")
		assert.AssertTrue(len(segments) == 4, "test")

		index, err := strconv.Atoi(segments[len(segments)-2])
		assert.AssertNil(err, "test")
		assert.AssertTrue(0 <= index && index < history, "test")
	})
}

func rolling(timeGranularity int, history int, fileAssert func(os.FileInfo)) {
	command := exec.Command("/bin/bash", "-c", "rm -rf /tmp/gtools")
	err := command.Run()
	assert.AssertNil(err, "test")

	direct := "/tmp/gtools/logs"
	fileName := "rolling"
	stop := false

	commonFileAppender, _ := log.NewFileAppender(&log.AppenderConfig{
		Layout:  "%d{2006-01-02 15:04:05.999} [%p] %m%n",
		Filters: nil,
		FileRollingPolicy: &log.RollingPolicy{
			Directory:       direct,
			FileName:        fileName,
			TimeGranularity: timeGranularity,
			MaxHistory:      history,
			MaxFileSize:     1,
		},
	})

	logger := log.NewLogger(log.Root, log.InfoLevel, false, []log.Appender{commonFileAppender})

	go func() {
		for !stop {
			logger.Info("now: '{}'", time.Now())

			time.Sleep(time.Microsecond)
		}
	}()

	time.Sleep(time.Second)

	go func() {
		for !stop {
			fileInfos, err := ioutil.ReadDir(direct)
			assert.AssertNil(err, "test")
			fileNum := len(fileInfos)
			// if a file is being renamed, we can't find it here
			assert.AssertTrue(fileNum <= history+1, "test")

			for _, fileInfo := range fileInfos {
				name := fileInfo.Name()
				if name == fileName+".log" {
					continue
				}
				fileAssert(fileInfo)
			}

			time.Sleep(time.Microsecond)
		}
	}()

	time.Sleep(time.Second * 3)

	commonFileAppender.Destroy()
	stop = true

	time.Sleep(time.Millisecond * 10)
}
