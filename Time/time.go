package Time

import (
	"time"
)

// Time time.Time的别名，但拥有不同的序列化到json方法
// 本Struct会将其在序列化为json时格式化为2006-01-02 15:04:05这样的格式
type Time time.Time

const (
	timeFormart = "2006-01-02 15:04:05"
)

// UnmarshalJSON 从json反序列号
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormart+`"`, string(data), time.Local)
	*t = Time(now)
	return
}

// MarshalJSON 序列化到json
func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormart)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormart)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(timeFormart)
}
