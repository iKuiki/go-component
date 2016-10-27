package youdaosdk

import (
	"github.com/yinhui87/go-component/config"
	"testing"
)

func TestTranslate(t *testing.T) {
	sdk, err := NewYoudaoSdk(config.Env("YOUDAO_APP_KEY"), config.Env("YOUDAO_APP_SECRET"))
	if err != nil {
		t.Fatalf("SDK Init Error: %s\n", err.Error())
	}
	tRet, err := sdk.Translate("饿饿饿饿饿")
	if err != nil {
		t.Fatalf("Translate Error: %s\n", err.Error())
	}
	for k, v := range tRet {
		t.Logf("retult%d: %s\n", k, v)
	}
}
