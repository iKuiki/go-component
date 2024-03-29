package logger

var defaultLogger Logger

// Log 实际的日志输出器
type Log interface {
	Print(...interface{})
	Printf(format string, args ...interface{})
	Println(...interface{})
	Error(...interface{})
	Errorf(format string, args ...interface{})
	Warn(...interface{})
	Warnf(format string, args ...interface{})
	Info(...interface{})
	Infof(format string, args ...interface{})
	Debug(...interface{})
	Debugf(format string, args ...interface{})
	Fatal(...interface{})
	Fatalf(format string, args ...interface{})
	// 将内存中的日志同步到磁盘
	Sync()
}

// Logger 日志记录者
type Logger interface {
	// 包含v2以兼容
	Log

	With(values Values) Logger
}

// SetDefaultLogger 设置全局缺省logger对象
func SetDefaultLogger(logger Logger) {
	defaultLogger = logger
}
