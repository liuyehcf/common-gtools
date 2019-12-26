package log

import (
	"errors"
	"fmt"
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

func NewFileAppender(config *AppenderConfig) (*fileAppender, error) {
	policy := config.FileRollingPolicy

	if strings.Contains(policy.FileName, ".") {
		return nil, errors.New("file name contains '.'")
	}
	if TimeGranularityHour != policy.TimeGranularity && TimeGranularityDay != policy.TimeGranularity {
		return nil, errors.New("TimeGranularity only support 1(TimeGranularityHour) and 2(TimeGranularityDay)")
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

	switch policy.TimeGranularity {
	case TimeGranularityHour:
		_, err := appender.cron.AddFunc("@hourly", func() {
			appender.rollingByTimer()
		})
		if err != nil {
			return nil, err
		}
		break
	case TimeGranularityDay:
		_, err := appender.cron.AddFunc("@daily", func() {
			appender.rollingByTimer()
		})
		if err != nil {
			return nil, err
		}
		break
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
	fileMetas := make([]*fileMeta, 0)

	files, err := ioutil.ReadDir(appender.policy.Directory)
	if err != nil {
		return fileMetas
	}

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

	_ = appender.createFileIfNecessary()
}

func (appender *fileAppender) rollingFilesByDayGranularity(allRollingFileMetas fileMetaSlice) {
	now := time.Now()
	dayFormatted := now.Format(formatDay)
	dayTime, _ := time.Parse(formatDay, dayFormatted)
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
