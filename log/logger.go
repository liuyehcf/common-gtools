package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	TraceLevel = 1
	DebugLevel = 2
	InfoLevel  = 3
	WarnLevel  = 4
	ErrorLevel = 5

	Root = "ROOT"
)

var (
	loggers        = make(map[string]*loggerImpl, 0)
	lock           = new(sync.RWMutex)
	virtualLoggers = make(map[string]*virtualLogger, 0)
	virtualLock    = new(sync.RWMutex)
	rootLogger     *loggerImpl
)

func getLogger(name string) (*loggerImpl, bool) {
	lock.RLock()
	defer lock.RUnlock()

	logger, ok := loggers[name]
	return logger, ok
}

func setOrReplaceLogger(name string, logger *loggerImpl) {
	lock.Lock()
	defer lock.Unlock()

	if _, ok := loggers[name]; ok {
		fmt.Printf("logger '%s' is replaced\n", name)
	}

	loggers[name] = logger
}

func foreachLogger(f func(key string, value *loggerImpl)) {
	lock.RLock()
	defer lock.RUnlock()

	for key, value := range loggers {
		f(key, value)
	}
}

func getVirtualLogger(name string) (*virtualLogger, bool) {
	virtualLock.RLock()
	defer virtualLock.RUnlock()

	logger, ok := virtualLoggers[name]
	return logger, ok
}

func setVirtualLoggerIfNotExist(name string, logger *virtualLogger) bool {
	virtualLock.Lock()
	defer virtualLock.Unlock()

	if _, ok := virtualLoggers[name]; ok {
		return false
	}

	virtualLoggers[name] = logger

	return true
}

func foreachVirtualLogger(f func(key string, value *virtualLogger)) {
	virtualLock.RLock()
	defer virtualLock.RUnlock()

	for key, value := range virtualLoggers {
		f(key, value)
	}
}

// logger interface
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
	logger, ok := getVirtualLogger(name)

	if ok {
		return logger
	}

	setVirtualLoggerIfNotExist(name, &virtualLogger{
		name:   name,
		target: nil,
	})

	logger, _ = getVirtualLogger(name)

	return logger
}

func getTargetLogger(name string, level int) Logger {
	logger, ok := getLogger(name)

	if ok {
		return logger
	}

	return newLoggerImpl(name, level, true, nil)
}

type loggerImpl struct {
	name       string
	level      int
	additivity bool
	appenders  []Appender
	parent     *loggerImpl
}

func NewLogger(name string, level int, additivity bool, appenders []Appender) Logger {
	// create logger impl
	newLoggerImpl(name, level, additivity, appenders)

	// always return virtual logger
	return GetLogger(name)
}

func newLoggerImpl(name string, level int, additivity bool, appenders []Appender) *loggerImpl {
	var logger *loggerImpl

	if isRoot(name) {
		name = Root
		logger = &loggerImpl{
			name:       name,
			level:      level,
			additivity: false,
			appenders:  appenders,
			parent:     nil,
		}
		rootLogger = logger

		// reset all non-root loggers' parent field
		foreachLogger(func(key string, value *loggerImpl) {
			if !isRoot(key) {
				value.parent = rootLogger
			}
		})
	} else {
		logger = &loggerImpl{
			name:       name,
			level:      level,
			additivity: additivity,
			appenders:  appenders,
			parent:     rootLogger,
		}
	}

	setOrReplaceLogger(name, logger)

	// clean bind status between virtual logger and target logger
	// this bind status will be rebuild later automatically
	foreachVirtualLogger(func(key string, value *virtualLogger) {
		if value != nil {
			value.target = nil
		}
	})

	return logger
}

func isRoot(name string) bool {
	return strings.ToUpper(name) == Root
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

	for l := logger; l != nil; l = l.parent {
		l.appendLoopOnAppenders(event)
		if !l.additivity {
			break
		}
	}
}

func (logger *loggerImpl) appendLoopOnAppenders(event *LoggingEvent) {
	if logger.appenders != nil {
		for _, appender := range logger.appenders {
			if appender != nil {
				appender.DoAppend(event)
			}
		}
	}
}

// virtual logger
// user may get logger before target logger created
// virtual logger will guarantee target logger will be bound at the right time
type virtualLogger struct {
	name   string
	level  int
	target Logger
}

func (logger *virtualLogger) Name() string {
	return logger.name
}

func (logger *virtualLogger) IsTraceEnabled() bool {
	logger.buildBoundStatusIfNecessary()

	// target may be null if target logger is created or replaced
	target := logger.target
	if target == nil {
		return false
	}
	return target.IsTraceEnabled()
}

func (logger *virtualLogger) Trace(format string, values ...interface{}) {
	logger.buildBoundStatusIfNecessary()

	// target may be null if target logger is created or replaced
	target := logger.target
	if target == nil {
		return
	}
	target.Trace(format, values...)
}

func (logger *virtualLogger) IsDebugEnabled() bool {
	logger.buildBoundStatusIfNecessary()

	// target may be null if target logger is created or replaced
	target := logger.target
	if target == nil {
		return false
	}
	return target.IsDebugEnabled()
}

func (logger *virtualLogger) Debug(format string, values ...interface{}) {
	logger.buildBoundStatusIfNecessary()

	// target may be null if target logger is created or replaced
	target := logger.target
	if target == nil {
		return
	}
	target.Debug(format, values...)
}

func (logger *virtualLogger) IsInfoEnabled() bool {
	logger.buildBoundStatusIfNecessary()

	// target may be null if target logger is created or replaced
	target := logger.target
	if target == nil {
		return false
	}
	return target.IsInfoEnabled()
}

func (logger *virtualLogger) Info(format string, values ...interface{}) {
	logger.buildBoundStatusIfNecessary()

	// target may be null if target logger is created or replaced
	target := logger.target
	if target == nil {
		return
	}
	target.Info(format, values...)
}

func (logger *virtualLogger) IsWarnEnabled() bool {
	logger.buildBoundStatusIfNecessary()

	// target may be null if target logger is created or replaced
	target := logger.target
	if target == nil {
		return false
	}
	return target.IsWarnEnabled()
}

func (logger *virtualLogger) Warn(format string, values ...interface{}) {
	logger.buildBoundStatusIfNecessary()

	// target may be null if target logger is created or replaced
	target := logger.target
	if target == nil {
		return
	}
	target.Warn(format, values...)
}

func (logger *virtualLogger) IsErrorEnabled() bool {
	logger.buildBoundStatusIfNecessary()

	// target may be null if target logger is created or replaced
	target := logger.target
	if target == nil {
		return false
	}
	return target.IsErrorEnabled()
}

func (logger *virtualLogger) Error(format string, values ...interface{}) {
	logger.buildBoundStatusIfNecessary()

	// target may be null if target logger is created or replaced
	target := logger.target
	if target == nil {
		return
	}
	target.Error(format, values...)
}

func (logger *virtualLogger) buildBoundStatusIfNecessary() {
	if logger.target != nil {
		return
	}

	logger.target = getTargetLogger(logger.name, logger.level)
	return
}

func init() {
	initConversion()

	stdoutAppender := NewWriterAppender(&AppenderConfig{
		Layout:    "%d{2006-01-02 15:04:05.999} [%p]-[%c]-[%L] --- %m%n",
		Filters:   nil,
		Writer:    os.Stdout,
		NeedClose: false,
	})

	rootLogger = newLoggerImpl(Root, InfoLevel, false, []Appender{stdoutAppender})
}
