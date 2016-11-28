package facebooksdk

import (
	"testing"
)

func TestGetUserInfo(t *testing.T) {
	accessToken := ""
	sdk := FacebookSdk{}
	userInfo, err := sdk.GetUserInfo(accessToken)
	if err != nil {
		t.Fatalf("GetUserInfo Error: %s\n", err.Error())
	}
	t.Logf("userInfo: %v\n", userInfo)
}
