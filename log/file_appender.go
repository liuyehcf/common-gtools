package log

import (
	"fmt"
	"github.com/liuyehcf/common-gtools/assert"
	cr "github.com/robfig/cron/v3"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	TimeGranularityHour = int(1)
	TimeGranularityDay  = int(2)
	formatDay           = "2006-01-02"
	emptyString         = ""
	fileSuffix          = ".log"
	pathSeparator       = string(os.PathSeparator)
)

type fileMeta struct {
	// exclude directory
	abstractPath string
	day          string
	hour         string
	index        string
	dayValue     int
	hourValue    int
	indexValue   int
}

func newFileMeta(abstractPath string, day string, hour string, index string) *fileMeta {
	var dayValue int
	var hourValue int
	var indexValue int
	var err error

	if day == emptyString {
		dayValue = -1
	} else {
		dayTime, err := time.Parse(formatDay, day)
		if err != nil {
			return nil
		}

		dayValue = dayTime.Second()
	}

	if hour == emptyString {
		hourValue = -1
	} else {
		hourValue, err = strconv.Atoi(hour)
		if err != nil {
			return nil
		}
	}

	if index == emptyString {
		indexValue = -1
	} else {
		indexValue, err = strconv.Atoi(index)
		if err != nil {
			return nil
		}
	}

	return &fileMeta{
		abstractPath: abstractPath,
		day:          day,
		hour:         hour,
		index:        index,
		dayValue:     dayValue,
		hourValue:    hourValue,
		indexValue:   indexValue,
	}
}

type fileMetaSlice []*fileMeta

func (slice fileMetaSlice) Len() int {
	return len(slice)
}

// smaller means older files
func (slice fileMetaSlice) Less(i, j int) bool {
	left := slice[i]
	right := slice[j]
	if left.dayValue < right.dayValue {
		return true
	} else if left.dayValue > right.dayValue {
		return false
	} else {
		if left.hourValue < right.hourValue {
			return true
		} else if left.hourValue > right.hourValue {
			return false
		} else {
			return left.indexValue >= right.indexValue
		}
	}
}

func (slice fileMetaSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type RollingPolicy struct {
	// directory for log
	Directory string

	// file name without suffix
	// if your log file is `default.log`, then just pass `default` here
	FileName string

	// only support TimeGranularityHour and TimeGranularityDay
	TimeGranularity int

	// maximum history of rolling logs
	MaxHistory int

	// maximum size of log file
	MaxFileSize int64
}

type fileAppender struct {
	abstractAppender
	policy           *RollingPolicy
	cron             *cr.Cron
	file             *os.File
	fileAbstractPath string
	fileRelativePath string
	fileAbstractName string
}

func NewFileAppender(config *AppenderConfig) *fileAppender {
	policy := config.FileRollingPolicy
	assert.AssertFalse(strings.HasSuffix(policy.Directory, pathSeparator), "directory ends with path separator")
	assert.AssertFalse(strings.Contains(policy.FileName, "."), "file name contains '.'")
	assert.AssertTrue(TimeGranularityHour == policy.TimeGranularity || TimeGranularityDay == policy.TimeGranularity, "TimeGranularity only support 1(TimeGranularityHour) and 2(TimeGranularityDay)")
	assert.AssertFalse(policy.MaxHistory < 1, "MaxHistory must large than 0")
	assert.AssertFalse(policy.MaxFileSize < 1, "MaxFileSize must large than 0")

	fileRelativePath := policy.FileName + fileSuffix
	appender := &fileAppender{
		abstractAppender: abstractAppender{
			encoder: newPatternEncoder(config.Layout),
			filters: config.Filters,
			lock:    new(sync.Mutex),
			queue:   make(chan []byte, 1024),
		},
		policy:           policy,
		cron:             cr.New(),
		fileRelativePath: fileRelativePath,
		fileAbstractPath: policy.Directory + pathSeparator + fileRelativePath,
		fileAbstractName: policy.Directory + pathSeparator + policy.FileName,
	}

	appender.cron.Start()
	appender.createDirectoryIfNecessary()
	appender.createFileIfNecessary()

	go appender.onEventLoop()
	switch policy.TimeGranularity {
	case TimeGranularityHour:
		_, err := appender.cron.AddFunc("@hourly", func() {
			appender.rollingByTimer()
		})
		assert.AssertNil(err, "failed to add cron func")
		break
	case TimeGranularityDay:
		_, err := appender.cron.AddFunc("@daily", func() {
			appender.rollingByTimer()
		})
		assert.AssertNil(err, "failed to add cron func")
		break
	}

	return appender
}

func (appender *fileAppender) Destroy() {
	lock.Lock()
	defer lock.Unlock()
	executeIgnorePanic(func() {
		appender.cron.Stop()
	})
	executeIgnorePanic(func() {
		_ = appender.file.Close()
	})
	executeIgnorePanic(func() {
		close(appender.queue)
	})
	appender.isDestroyed = true
}

func (appender *fileAppender) onEventLoop() {
	defer func() {
		recover()
	}()
	for !appender.isDestroyed {
		content := <-appender.queue
		appender.rollingIfFileSizeExceeded()
		appender.write(content)
	}
}

func (appender *fileAppender) rollingIfFileSizeExceeded() {
	info, err := appender.file.Stat()
	assert.AssertNil(err, "failed to get file stat")

	if info.Size() >= appender.policy.MaxFileSize {
		appender.doSizeRolling()
	}
}

func (appender *fileAppender) rollingByTimer() {
	appender.doSizeRolling()
}

func (appender *fileAppender) doSizeRolling() {
	appender.lock.Lock()
	defer appender.lock.Unlock()

	fileMetas := appender.getAllRollingFileMetas()

	switch appender.policy.TimeGranularity {
	case TimeGranularityHour:
		appender.rollingFilesByHourGranularity(fileMetas)
		break
	case TimeGranularityDay:
		appender.rollingFilesByDayGranularity(fileMetas)
		break
	}
}

func (appender *fileAppender) getAllRollingFileMetas() []*fileMeta {
	files, err := ioutil.ReadDir(appender.policy.Directory)
	assert.AssertNil(err, "failed to read directory")

	fileMetas := make([]*fileMeta, 0)

	for _, file := range files {
		if strings.HasPrefix(file.Name(), appender.policy.FileName) &&
			strings.HasSuffix(file.Name(), fileSuffix) {
			fileMeta := appender.parseRollingFileInfo(file)

			if fileMeta != nil {
				fileMetas = append(fileMetas, fileMeta)
			}
		}
	}

	return fileMetas
}

func (appender *fileAppender) parseRollingFileInfo(fileInfo os.FileInfo) *fileMeta {
	abstractPath := appender.policy.Directory + pathSeparator + fileInfo.Name()

	// skip current file
	if fileInfo.Name() == appender.fileRelativePath {
		return nil
	}

	segments := strings.Split(fileInfo.Name(), ".")
	if segments == nil {
		return nil
	}
	segmentLen := len(segments)

	switch appender.policy.TimeGranularity {
	case TimeGranularityHour:
		// xxx.2006-01-02.08.1.log
		if segmentLen == 5 {
			return newFileMeta(abstractPath, segments[1], segments[2], segments[3])
		} else {
			return nil
		}
	case TimeGranularityDay:
		// xxx.2006-01-02.1.log
		if segmentLen == 4 {
			return newFileMeta(abstractPath, segments[1], "", segments[2])
		} else {
			return nil
		}
	}

	return nil
}

func (appender *fileAppender) rollingFilesByHourGranularity(allRollingFileMetas fileMetaSlice) {
	now := time.Now()
	dayFormatted := now.Format(formatDay)
	hour := now.Hour()

	policy := appender.policy

	if len(allRollingFileMetas) >= policy.MaxHistory {
		sort.Sort(allRollingFileMetas)
		maxRemainHistory := policy.MaxHistory - 1
		removedFileMetas := allRollingFileMetas[:len(allRollingFileMetas)-maxRemainHistory]

		for _, removedFileMeta := range removedFileMetas {
			_ = os.Remove(removedFileMeta.abstractPath)
		}

		allRollingFileMetas = allRollingFileMetas[len(allRollingFileMetas)-maxRemainHistory:]
	}

	fileMetasOfCurHour := make(fileMetaSlice, 0)

	for _, fileMeta := range allRollingFileMetas {
		if fileMeta.hourValue == hour {
			fileMetasOfCurHour = append(fileMetasOfCurHour, fileMeta)
		}
	}

	sort.Sort(fileMetasOfCurHour)

	for i := 0; i < fileMetasOfCurHour.Len(); i += 1 {
		fileMeta := fileMetasOfCurHour[i]

		// dir/xxx.2006-01-02.08.1.log
		_ = os.Rename(fileMeta.abstractPath,
			fmt.Sprintf("%s.%s.%s.%d%s", appender.fileAbstractName, fileMeta.day, fileMeta.hour, fileMeta.indexValue+1, fileSuffix))
	}

	_ = appender.file.Close()

	_ = os.Rename(appender.fileAbstractPath,
		fmt.Sprintf("%s.%s.%02d.%d%s", appender.fileAbstractName, dayFormatted, hour, 0, fileSuffix))

	appender.createFileIfNecessary()
}

func (appender *fileAppender) rollingFilesByDayGranularity(allRollingFileMetas fileMetaSlice) {
	now := time.Now()
	dayFormatted := now.Format(formatDay)
	dayTime, err := time.Parse(formatDay, dayFormatted)
	assert.AssertNil(err, "failed to parse day time")
	day := dayTime.Second()

	policy := appender.policy

	if len(allRollingFileMetas) >= policy.MaxHistory {
		sort.Sort(allRollingFileMetas)
		maxRemainHistory := policy.MaxHistory - 1
		removedFileMetas := allRollingFileMetas[:len(allRollingFileMetas)-maxRemainHistory]

		for _, removedFileMeta := range removedFileMetas {
			_ = os.Remove(removedFileMeta.abstractPath)
		}

		allRollingFileMetas = allRollingFileMetas[len(allRollingFileMetas)-maxRemainHistory:]
	}

	fileMetasOfCurDay := make(fileMetaSlice, 0)

	for _, fileMeta := range allRollingFileMetas {
		if fileMeta.dayValue == day {
			fileMetasOfCurDay = append(fileMetasOfCurDay, fileMeta)
		}
	}

	sort.Sort(fileMetasOfCurDay)

	for i := 0; i < fileMetasOfCurDay.Len(); i += 1 {
		fileMeta := fileMetasOfCurDay[i]

		_ = os.Rename(fileMeta.abstractPath,
			fmt.Sprintf("%s.%s.%d%s", appender.fileAbstractName, fileMeta.day, fileMeta.indexValue+1, fileSuffix))
	}

	_ = appender.file.Close()

	_ = os.Rename(appender.fileAbstractPath,
		fmt.Sprintf("%s.%s.%d%s", appender.fileAbstractName, dayFormatted, 0, fileSuffix))

	appender.createFileIfNecessary()
}

func (appender *fileAppender) createDirectoryIfNecessary() {
	err := os.MkdirAll(appender.policy.Directory, os.ModePerm)
	assert.AssertNil(err, "failed to create directory")
}

func (appender *fileAppender) createFileIfNecessary() {
	var err error
	appender.file, err = os.OpenFile(appender.fileAbstractPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	assert.AssertNil(err, "failed to open file")
}

func (appender *fileAppender) write(bytes []byte) {
	appender.lock.Lock()
	defer appender.lock.Unlock()
	_, err := appender.file.Write(bytes)
	assert.AssertNil(err, "failed to write content")
}
