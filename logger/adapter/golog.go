package adapter

import (
	"fmt"

	"github.com/iKuiki/go-component/logger"
	"github.com/kataras/golog"
)

// gologAdapter golog Adapter for logger interface
type gologAdapter struct {
	*golog.Logger
}

// NewGologAdapter 创建golog对象
func NewGologAdapter(enableDebug bool) (logger logger.Logger) {
	adapter := &gologAdapter{golog.New()}
	if enableDebug {
		adapter.SetLevel("debug")
	}
	return adapter
}

// With context support
func (adapter *gologAdapter) With(values logger.Values) logger.Logger {
	if values.GetString("golog_prefix") == "true" { // 明确启用prefix的context才尝试抓取
		vMap := make(map[string]string)
		values.RangeValues(func(k, v string) bool {
			vMap[k] = v
			return true
		})
		if len(vMap) > 0 {
			ret := &gologAdapter{
				Logger: adapter.Clone().SetPrefix(fmt.Sprint(vMap)),
			}
			return ret
		}
	}
	return adapter
}

func (adapter *gologAdapter) Sync() {
	// seelog无需sync
	return
}
