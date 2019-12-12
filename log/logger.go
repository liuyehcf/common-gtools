package log

import (
	"runtime"
	"time"
)

const (
	TraceLevel = 1
	DebugLevel = 2
	InfoLevel  = 3
	WarnLevel  = 4
	ErrorLevel = 5
)

type Logger interface {
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

	// error log
	Error(format string, values ...interface{})

	// whether error enabled
	// if not, all the debug/info/warn/error will be ignored
	IsErrorEnabled() bool
}

type DefaultLogger struct {
	Level     int
	Appenders []Appender
}

func (logger *DefaultLogger) IsTraceEnabled() bool {
	return logger.Level <= TraceLevel
}

func (logger *DefaultLogger) Trace(format string, values ...interface{}) {
	if logger.IsTraceEnabled() {
		logger.callAllAppenders(TraceLevel, format, values...)
	}
}

func (logger *DefaultLogger) IsDebugEnabled() bool {
	return logger.Level <= DebugLevel
}

func (logger *DefaultLogger) Debug(format string, values ...interface{}) {
	if logger.IsDebugEnabled() {
		logger.callAllAppenders(DebugLevel, format, values...)
	}
}

func (logger *DefaultLogger) IsInfoEnabled() bool {
	return logger.Level <= InfoLevel
}

func (logger *DefaultLogger) Info(format string, values ...interface{}) {
	if logger.IsInfoEnabled() {
		logger.callAllAppenders(InfoLevel, format, values...)
	}
}

func (logger *DefaultLogger) IsWarnEnabled() bool {
	return logger.Level <= WarnLevel
}

func (logger *DefaultLogger) Warn(format string, values ...interface{}) {
	if logger.IsWarnEnabled() {
		logger.callAllAppenders(WarnLevel, format, values...)
	}
}

func (logger *DefaultLogger) IsErrorEnabled() bool {
	return logger.Level <= ErrorLevel
}

func (logger *DefaultLogger) Error(format string, values ...interface{}) {
	if logger.IsErrorEnabled() {
		logger.callAllAppenders(ErrorLevel, format, values...)
	}
}

func (logger *DefaultLogger) callAllAppenders(level int, format string, values ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	event := &LoggingEvent{
		Level:     level,
		Timestamp: time.Now(),
		File:      file,
		Line:      line,
		Message:   format,
		Values:    values,
	}

	for _, appender := range logger.Appenders {
		appender.DoAppend(event)
	}
}
