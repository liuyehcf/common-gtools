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
	lock           = new(sync.Mutex)
	virtualLoggers = make(map[string]*virtualLogger, 0)
	virtualLock    = new(sync.Mutex)
	rootLogger     *loggerImpl
)

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
	logger, ok := virtualLoggers[name]

	if ok {
		return logger
	}

	virtualLock.Lock()
	defer virtualLock.Unlock()
	logger, ok = virtualLoggers[name]

	if ok {
		return logger
	}

	logger = &virtualLogger{
		name:   name,
		target: nil,
	}

	virtualLoggers[name] = logger

	return logger
}

func getTargetLogger(name string) Logger {
	logger, ok := loggers[name]

	if ok {
		return logger
	}

	return rootLogger
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
	lock.Lock()
	defer lock.Unlock()

	if _, ok := loggers[name]; ok {
		fmt.Printf("logger '%s' is replaced\n", name)
	}

	// clean bind status between virtual logger and true logger
	// this bind status will be rebuild later automatically
	for _, virtualLogger := range virtualLoggers {
		if virtualLogger != nil {
			virtualLogger.target = nil
		}
	}

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
		loggers[name] = logger
		rootLogger = logger

		for key, value := range loggers {
			if !isRoot(key) {
				value.parent = rootLogger
			}
		}

		return logger
	} else {
		logger = &loggerImpl{
			name:       name,
			level:      level,
			additivity: additivity,
			appenders:  appenders,
			parent:     rootLogger,
		}
		loggers[name] = logger

		return logger
	}
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
	for _, appender := range logger.appenders {
		appender.DoAppend(event)
	}
}

// virtual logger
// user may get logger before true logger created
// virtual logger will guarantee true logger will be bound at the right time
type virtualLogger struct {
	name   string
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

	logger.target = getTargetLogger(logger.name)
	return
}

func init() {
	initConversion()

	stdoutAppender := NewWriterAppender(&AppenderConfig{
		Layout:    "%d{2006-01-02 15:04:05.999} [%p] %m%n",
		Filters:   nil,
		Writer:    os.Stdout,
		NeedClose: false,
	})

	rootLogger = newLoggerImpl(Root, InfoLevel, false, []Appender{stdoutAppender})
}
