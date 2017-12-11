package Time

import (
	"strconv"
	"time"
)

// UnixTime time.UnixTime的别名，但拥有不同的序列化到json方法
// 本Struct会将其在序列化为json时格式化为10位unix时间戳
type UnixTime time.Time

// UnmarshalJSON 从json反序列号
func (t *UnixTime) UnmarshalJSON(data []byte) (err error) {
	if i, err := strconv.ParseInt(string(data), 10, 64); err == nil {
		*t = UnixTime(time.Unix(i, 0))
	}
	return
}

// MarshalJSON 序列化到json
func (t UnixTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

func (t UnixTime) String() string {
	return time.Time(t).Format(timeFormart)
}
