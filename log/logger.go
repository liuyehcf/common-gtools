package log

type Logger interface {
	// whether debug enabled
	// if not, all the debug will be ignored
	IsDebugEnabled() bool

	// whether info enabled
	// if not, all the debug/info will be ignored
	IsInfoEnabled() bool

	// whether warn enabled
	// if not, all the debug/info/warn will be ignored
	IsWarnEnabled() bool

	// whether error enabled
	// if not, all the debug/info/warn/error will be ignored
	IsErrorEnabled() bool

	// debug log
	Debug(format string, values interface{})

	// info log
	Info(format string, values interface{})

	// warn log
	Warn(format string, values interface{})

	// error log
	Error(format string, values interface{})
}

type DefaultLogger struct {
}
