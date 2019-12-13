package log

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const (
	TraceLevel = 1
	DebugLevel = 2
	InfoLevel  = 3
	WarnLevel  = 4
	ErrorLevel = 5

	Root = "root"
)

var (
	loggerPool         = make(map[string]Logger, 0)
	loggerPoolLock     = new(sync.Mutex)
	lazyLoggerPool     = make(map[string]*lazyBoundLogger, 0)
	lazyLoggerPoolLock = new(sync.Mutex)
	rootLogger         Logger
	mLogger            Logger = new(muteLogger)
)

type Logger interface {
	// get logger name
	Name() string

	// whether debug enabled
	// if not, all the debug will be ignored
	IsTraceEnabled() bool

	// debug log
	Trace(format string, values ...interface{})

	// whether debug enabled
	// if not, all the debug will be ignored
	IsDebugEnabled() bool

	// debug log
	Debug(format string, values ...interface{})

	// whether info enabled
	// if not, all the debug/info will be ignored
	IsInfoEnabled() bool

	// info log
	Info(format string, values ...interface{})

	// whether warn enabled
	// if not, all the debug/info/warn will be ignored
	IsWarnEnabled() bool

	// warn log
	Warn(format string, values ...interface{})

	// whether error enabled
	// if not, all the debug/info/warn/error will be ignored
	IsErrorEnabled() bool

	// error log
	Error(format string, values ...interface{})
}

func GetLogger(name string) Logger {
	logger, ok := lazyLoggerPool[name]

	if ok {
		return logger
	}

	lazyLoggerPoolLock.Lock()
	defer lazyLoggerPoolLock.Unlock()
	logger, ok = lazyLoggerPool[name]

	if ok {
		return logger
	}

	logger = &lazyBoundLogger{
		name:    name,
		target:  nil,
		isBound: false,
	}

	lazyLoggerPool[name] = logger

	return logger
}

func getTargetLogger(name string) Logger {
	logger, ok := loggerPool[name]

	if ok {
		return logger
	}

	if rootLogger != nil {
		return rootLogger
	}

	return mLogger
}

type loggerImpl struct {
	name      string
	level     int
	appenders []Appender
}

func NewLogger(name string, level int, appenders []Appender) Logger {
	logger := &loggerImpl{name: name, level: level, appenders: appenders}

	loggerPoolLock.Lock()
	defer loggerPoolLock.Unlock()

	if _, ok := loggerPool[name]; ok {
		fmt.Printf("logger '%s' is replaced\n", name)
	}
	loggerPool[name] = logger

	for _, lazyLogger := range lazyLoggerPool {
		if lazyLogger != nil {
			lazyLogger.isBound = false
			lazyLogger.target = nil
		}
	}

	if name == Root {
		rootLogger = logger
	}

	// always return lazy init logger
	return GetLogger(name)
}

func (logger *loggerImpl) Name() string {
	return logger.name
}

func (logger *loggerImpl) IsTraceEnabled() bool {
	return logger.level <= TraceLevel
}

func (logger *loggerImpl) Trace(format string, values ...interface{}) {
	if logger.IsTraceEnabled() {
		logger.callAllAppenders(TraceLevel, format, values...)
	}
}

func (logger *loggerImpl) IsDebugEnabled() bool {
	return logger.level <= DebugLevel
}

func (logger *loggerImpl) Debug(format string, values ...interface{}) {
	if logger.IsDebugEnabled() {
		logger.callAllAppenders(DebugLevel, format, values...)
	}
}

func (logger *loggerImpl) IsInfoEnabled() bool {
	return logger.level <= InfoLevel
}

func (logger *loggerImpl) Info(format string, values ...interface{}) {
	if logger.IsInfoEnabled() {
		logger.callAllAppenders(InfoLevel, format, values...)
	}
}

func (logger *loggerImpl) IsWarnEnabled() bool {
	return logger.level <= WarnLevel
}

func (logger *loggerImpl) Warn(format string, values ...interface{}) {
	if logger.IsWarnEnabled() {
		logger.callAllAppenders(WarnLevel, format, values...)
	}
}

func (logger *loggerImpl) IsErrorEnabled() bool {
	return logger.level <= ErrorLevel
}

func (logger *loggerImpl) Error(format string, values ...interface{}) {
	if logger.IsErrorEnabled() {
		logger.callAllAppenders(ErrorLevel, format, values...)
	}
}

func (logger *loggerImpl) callAllAppenders(level int, format string, values ...interface{}) {
	_, file, line, _ := runtime.Caller(3)
	event := &LoggingEvent{
		Name:      logger.name,
		Level:     level,
		Timestamp: time.Now(),
		File:      file,
		Line:      line,
		Message:   format,
		Values:    values,
	}

	for _, appender := range logger.appenders {
		appender.DoAppend(event)
	}
}

// lazy bound logger
type lazyBoundLogger struct {
	name    string
	target  Logger
	isBound bool
}

func (logger *lazyBoundLogger) Name() string {
	return logger.name
}

func (logger *lazyBoundLogger) IsTraceEnabled() bool {
	logger.lazyInitIfNecessary()

	// target may be null if target logger is created or replaced
	target := logger.target
	if target == nil {
		return false
	}
	return target.IsTraceEnabled()
}

func (logger *lazyBoundLogger) Trace(format string, values ...interface{}) {
	logger.lazyInitIfNecessary()

	// target may be null if target logger is created or replaced
	target := logger.target
	if target == nil {
		return
	}
	target.Trace(format, values...)
}

func (logger *lazyBoundLogger) IsDebugEnabled() bool {
	logger.lazyInitIfNecessary()

	// target may be null if target logger is created or replaced
	target := logger.target
	if target == nil {
		return false
	}
	return target.IsDebugEnabled()
}

func (logger *lazyBoundLogger) Debug(format string, values ...interface{}) {
	logger.lazyInitIfNecessary()

	// target may be null if target logger is created or replaced
	target := logger.target
	if target == nil {
		return
	}
	target.Debug(format, values...)
}

func (logger *lazyBoundLogger) IsInfoEnabled() bool {
	logger.lazyInitIfNecessary()

	// target may be null if target logger is created or replaced
	target := logger.target
	if target == nil {
		return false
	}
	return target.IsInfoEnabled()
}

func (logger *lazyBoundLogger) Info(format string, values ...interface{}) {
	logger.lazyInitIfNecessary()

	// target may be null if target logger is created or replaced
	target := logger.target
	if target == nil {
		return
	}
	target.Info(format, values...)
}

func (logger *lazyBoundLogger) IsWarnEnabled() bool {
	logger.lazyInitIfNecessary()

	// target may be null if target logger is created or replaced
	target := logger.target
	if target == nil {
		return false
	}
	return target.IsWarnEnabled()
}

func (logger *lazyBoundLogger) Warn(format string, values ...interface{}) {
	logger.lazyInitIfNecessary()

	// target may be null if target logger is created or replaced
	target := logger.target
	if target == nil {
		return
	}
	target.Warn(format, values...)
}

func (logger *lazyBoundLogger) IsErrorEnabled() bool {
	logger.lazyInitIfNecessary()

	// target may be null if target logger is created or replaced
	target := logger.target
	if target == nil {
		return false
	}
	return target.IsErrorEnabled()
}

func (logger *lazyBoundLogger) Error(format string, values ...interface{}) {
	logger.lazyInitIfNecessary()

	// target may be null if target logger is created or replaced
	target := logger.target
	if target == nil {
		return
	}
	target.Error(format, values...)
}

func (logger *lazyBoundLogger) lazyInitIfNecessary() {
	if logger.isBound {
		return
	}

	logger.target = getTargetLogger(logger.name)
	logger.isBound = true
	return
}

// this logger will return if no logger matches
type muteLogger struct {
}

func (logger *muteLogger) Name() string {
	return "mute"
}

func (logger *muteLogger) IsTraceEnabled() bool {
	return false
}

func (logger *muteLogger) Trace(format string, values ...interface{}) {
	return
}

func (logger *muteLogger) IsDebugEnabled() bool {
	return false
}

func (logger *muteLogger) Debug(format string, values ...interface{}) {
	return
}

func (logger *muteLogger) IsInfoEnabled() bool {
	return false
}

func (logger *muteLogger) Info(format string, values ...interface{}) {
	return
}

func (logger *muteLogger) IsWarnEnabled() bool {
	return false
}

func (logger *muteLogger) Warn(format string, values ...interface{}) {
	return
}

func (logger *muteLogger) IsErrorEnabled() bool {
	return false
}

func (logger *muteLogger) Error(format string, values ...interface{}) {
	return
}
