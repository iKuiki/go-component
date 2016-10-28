package umengsdk

import (
	"fmt"
	"github.com/yinhui87/go-component/config"
	"testing"
	"time"
)

func TestSendMessage(t *testing.T) {
	sdk, err := NewUmengSdk(config.Env("UMENG_APP_KEY"), config.Env("UMENG_APP_SECRET"))
	if err != nil {
		t.Fatalf("Umeng Sdk Init Fail: %s\n", err.Error())
	}
	now_time := uint32(time.Now().Unix())
	umengCommon := UmengCommon{
		Appkey:         config.Env("UMENG_APP_KEY"),
		Timestamp:      fmt.Sprintf("%d", now_time),
		Type:           "unicast",
		DeviceTokens:   config.Env("TEST_ANDROID_DEVICE_TOKEN"),
		ProductionMode: "false",
		Description:    "sdk test message",
	}
	umengAndroid := UmengAndroid{
		UmengCommon: umengCommon,
		Payload: UmengAndroidPayload{
			DisplayType: "notification",
			Body: UmengAndroidPayloadBody{
				Ticker:    "TestMessage",
				Title:     "Test",
				Text:      "This is a test message",
				AfterOpen: "go_app",
			},
		},
	}
	msgId, err := sdk.SendAndroid(umengAndroid)
	if err != nil {
		t.Errorf("SendAndroid Error: %s\n", err.Error())
	} else {
		t.Logf("SendAndroid Success, MsgId: %s\n", msgId)
	}
}
