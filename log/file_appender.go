package log

import (
	"errors"
	"fmt"
	"github.com/liuyehcf/common-gtools/utils"
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
	TimeGranularityNone = int(0)
	TimeGranularityHour = int(1)
	TimeGranularityDay  = int(2)
	sizeRolling         = int(0)
	timerRolling        = int(1)
	formatDay           = "2006-01-02"
	emptyString         = ""
	fileSuffix          = ".log"
	pathSeparator       = string(os.PathSeparator)
)

var (
	timeGranularityMap = map[int]string{
		TimeGranularityHour: "@hourly",
		TimeGranularityDay:  "@daily",
	}
)

type fileMeta struct {
	// exclude directory
	abstractPath string
	day          string
	hour         string
	index        string
	dayValue     int64
	hourValue    int
	indexValue   int
}

func newFileMeta(abstractPath string, day string, hour string, index string) *fileMeta {
	var dayValue int64
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

		dayValue = dayTime.Unix()
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

func NewFileAppender(config *AppenderConfig) (*fileAppender, error) {
	policy := config.FileRollingPolicy

	if strings.Contains(policy.FileName, ".") {
		return nil, errors.New("file name contains '.'")
	}
	if TimeGranularityNone != policy.TimeGranularity &&
		TimeGranularityHour != policy.TimeGranularity &&
		TimeGranularityDay != policy.TimeGranularity {
		return nil, errors.New("TimeGranularity only support 0(TimeGranularityNone) or 1(TimeGranularityHour) or 2(TimeGranularityDay)")
	}
	if policy.MaxHistory < 1 {
		return nil, errors.New("MaxHistory must large than 0")
	}
	if policy.MaxFileSize < 1 {
		return nil, errors.New("MaxFileSize must large than 0")
	}

	for strings.HasSuffix(policy.Directory, pathSeparator) {
		size := len(policy.Directory)
		policy.Directory = policy.Directory[0 : size-1]
	}

	fileRelativePath := policy.FileName + fileSuffix
	encoder, err := newPatternEncoder(config.Layout)
	if err != nil {
		return nil, err
	}
	appender := &fileAppender{
		abstractAppender: abstractAppender{
			encoder: encoder,
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
	err = appender.createDirectoryIfNecessary()
	if err != nil {
		return nil, err
	}
	err = appender.createFileIfNecessary()
	if err != nil {
		return nil, err
	}

	timeGranularity, exist := timeGranularityMap[policy.TimeGranularity]
	if exist {
		_, err := appender.cron.AddFunc(timeGranularity, func() {
			appender.rollingByTimer()
		})
		if err != nil {
			return nil, err
		}
	}

	go appender.onEventLoop()

	return appender, nil
}

func (appender *fileAppender) Destroy() {
	lock.Lock()
	defer lock.Unlock()
	appender.isDestroyed = true
	executeIgnorePanic(func() {
		appender.cron.Stop()
	})
	executeIgnorePanic(func() {
		_ = appender.file.Close()
	})
	executeIgnorePanic(func() {
		close(appender.queue)
	})
}

func (appender *fileAppender) onEventLoop() {
	defer func() {
		recover()
	}()

	var content []byte
	var ok bool
	for !appender.isDestroyed {
		if content, ok = <-appender.queue; !ok {
			// channel is closed
			break
		}
		appender.rollingIfFileSizeExceeded()
		appender.write(content)
	}
}

func (appender *fileAppender) rollingIfFileSizeExceeded() {
	info, err := appender.file.Stat()
	if err != nil {
		return
	}

	appender.lock.Lock()
	defer appender.lock.Unlock()

	if info.Size() >= appender.policy.MaxFileSize {
		appender.doRolling(sizeRolling)
	}
}

func (appender *fileAppender) rollingByTimer() {
	info, err := appender.file.Stat()
	if err != nil {
		return
	}

	appender.lock.Lock()
	defer appender.lock.Unlock()

	if info.Size() > 0 {
		appender.doRolling(timerRolling)
	}
}

func (appender *fileAppender) doRolling(rollingType int) {
	fileMetas := appender.getAllRollingFileMetas()

	switch appender.policy.TimeGranularity {
	case TimeGranularityHour:
		appender.rollingFilesByHourGranularity(rollingType, fileMetas)
		break
	case TimeGranularityNone, TimeGranularityDay:
		appender.rollingFilesByDayGranularity(rollingType, fileMetas)
		break
	}
}

func (appender *fileAppender) getAllRollingFileMetas() []*fileMeta {
	fileMetas := make([]*fileMeta, 0)

	files, err := ioutil.ReadDir(appender.policy.Directory)
	if err != nil {
		return fileMetas
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), appender.policy.FileName) &&
			strings.HasSuffix(file.Name(), fileSuffix) {
			fileMeta := appender.parseRollingFileInfo(file)

			if utils.IsNotNil(fileMeta) {
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
	case TimeGranularityNone, TimeGranularityDay:
		// xxx.2006-01-02.1.log
		if segmentLen == 4 {
			return newFileMeta(abstractPath, segments[1], "", segments[2])
		} else {
			return nil
		}
	}

	return nil
}

func (appender *fileAppender) rollingFilesByHourGranularity(rollingType int, allRollingFileMetas fileMetaSlice) {
	var t time.Time
	if rollingType == timerRolling {
		t = time.Now().Add(-time.Hour)
	} else {
		t = time.Now()
	}
	dayFormatted := t.Format(formatDay)
	dayTime, _ := time.Parse(formatDay, dayFormatted)
	day := dayTime.Unix()
	hour := t.Hour()

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
		if fileMeta.dayValue == day &&
			fileMeta.hourValue == hour {
			fileMetasOfCurHour = append(fileMetasOfCurHour, fileMeta)
		}
	}

	_ = appender.file.Close()

	_ = os.Rename(appender.fileAbstractPath,
		fmt.Sprintf("%s.%s.%02d.%d%s", appender.fileAbstractName, dayFormatted, hour, len(fileMetasOfCurHour)+1, fileSuffix))

	_ = appender.createFileIfNecessary()
}

func (appender *fileAppender) rollingFilesByDayGranularity(rollingType int, allRollingFileMetas fileMetaSlice) {
	var t time.Time
	if rollingType == timerRolling {
		t = time.Now().Add(-24 * time.Hour)
	} else {
		t = time.Now()
	}
	dayFormatted := t.Format(formatDay)
	dayTime, _ := time.Parse(formatDay, dayFormatted)
	day := dayTime.Unix()

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

	_ = appender.file.Close()

	_ = os.Rename(appender.fileAbstractPath,
		fmt.Sprintf("%s.%s.%d%s", appender.fileAbstractName, dayFormatted, len(fileMetasOfCurDay)+1, fileSuffix))

	_ = appender.createFileIfNecessary()
}

func (appender *fileAppender) createDirectoryIfNecessary() error {
	return os.MkdirAll(appender.policy.Directory, os.ModePerm)
}

func (appender *fileAppender) createFileIfNecessary() error {
	var err error
	appender.file, err = os.OpenFile(appender.fileAbstractPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	return err
}

func (appender *fileAppender) write(bytes []byte) {
	appender.lock.Lock()
	defer appender.lock.Unlock()
	_, _ = appender.file.Write(bytes)
}
