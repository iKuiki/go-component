package test

import (
	"github.com/yinhui87/go-component/sdk/qiniusdk"
	"testing"
)

const (
	ACCESS_KEY = "Your ACCESS_KEY"
	SECRET_KEY = "Your SECRET_KEY"
	BUCKET     = "Your bucket"
)

func TestQiniu(t *testing.T) {
	qiniuSdk, err := qiniusdk.NewQiniuSdk(ACCESS_KEY, SECRET_KEY, BUCKET, 3600, 0)
	if err != nil {
		t.Fatal("qiniusdk.NewQiniuSdk error: " + err.Error())
	}
	info, err := qiniuSdk.GetImageBaseInfo("Your Filename")
	if err != nil {
		t.Fatal("qiniusdk.GetImageBaseInfo error: " + err.Error())
	}
	t.Logf("ImageInfo: %v", info)
}
