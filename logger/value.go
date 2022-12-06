package logger

// Values 日志上下文存储器
type Values interface {
	GetString(key string) (value string)
	// 迭代context内所有values
	RangeValues(rangeFn func(key string, value string) bool)
}
