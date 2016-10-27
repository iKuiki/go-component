package umengsdk

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/yinhui87/go-component/encoding"
	"github.com/yinhui87/go-component/util"
)

type UmengSdk struct {
	AccessKey string
	SecretKey string
}

const (
	sendUrl   = "http://msg.umeng.com/api/send"
	statusUrl = "http://msg.umeng.com/api/status"
	uploadUrl = "http://msg.umeng.com/upload"
)

//初始化UmengSdk
func NewUmengSdk(appKey, appSecret string) (umengSdk *UmengSdk, initError error) {
	if appKey == "" {
		return nil, errors.New("appKey不能为空！")
	}
	if appSecret == "" {
		return nil, errors.New("appSecret不能为空！")
	}
	server := &UmengSdk{
		AccessKey: appKey,
		SecretKey: appSecret,
	}
	return server, nil
}

func (this *UmengSdk) SendAndroid(umengAndroid UmengAndroid) (UmengResult, error) {
	sign := ""
	method := "POST"

	body, err := encoding.EncodeJson(umengAndroid)
	if err != nil {
		return UmengResult{}, err
	}
	sign = this.getSign(method, sendUrl, string(body))
	url := sendUrl + "?sign=" + sign

	var result []byte
	err = util.DefaultAjaxPool.Post(&util.Ajax{
		Url:          url,
		Data:         body,
		ResponseData: &result,
	})
	if err != nil {
		if _, ok := err.(*util.AjaxStatusCodeError); !ok {
			return UmengResult{}, err
		}
	}

	var finalResult UmengResult
	err = encoding.DecodeJson(result, &finalResult)
	if err != nil {
		return UmengResult{}, err
	}
	return finalResult, nil
}

func (this *UmengSdk) SendIOS(umengIOS UmengIOS) (UmengResult, error) {
	sign := ""
	method := "POST"

	body, err := encoding.EncodeJson(umengIOS)
	if err != nil {
		return UmengResult{}, err
	}
	sign = this.getSign(method, sendUrl, string(body))
	url := sendUrl + "?sign=" + sign

	var result []byte
	err = util.DefaultAjaxPool.Post(&util.Ajax{
		Url:          url,
		Data:         body,
		ResponseData: &result,
	})
	if err != nil {
		if _, ok := err.(*util.AjaxStatusCodeError); !ok {
			return UmengResult{}, err
		}
	}

	var finalResult UmengResult
	err = json.Unmarshal(result, &finalResult)
	if err != nil {
		return UmengResult{}, err
	}
	return finalResult, nil
}

func (this *UmengSdk) GetFileId(deviceToken string) (UmengFileResult, error) {
	sign := ""
	method := "POST"

	body, err := json.Marshal(UmengFile{
		Appkey:    this.AccessKey,
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
		Content:   deviceToken,
	})
	if err != nil {
		return UmengFileResult{}, err
	}
	sign = this.getSign(method, uploadUrl, string(body))
	url := uploadUrl + "?sign=" + sign

	var result []byte
	err = util.DefaultAjaxPool.Post(&util.Ajax{
		Url:          url,
		Data:         body,
		ResponseData: &result,
	})
	if err != nil {
		return UmengFileResult{}, err
	}

	var finalResult UmengFileResult
	err = json.Unmarshal(result, &finalResult)
	if err != nil {
		return UmengFileResult{}, err
	}
	return finalResult, nil
}

func (this *UmengSdk) GetStatus(taskId string) (UmengStatusResult, error) {
	sign := ""
	method := "POST"

	body, err := json.Marshal(UmengStatus{
		Appkey:    this.AccessKey,
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
		TaskId:    taskId,
	})
	if err != nil {
		return UmengStatusResult{}, err
	}
	sign = this.getSign(method, statusUrl, string(body))
	url := statusUrl + "?sign=" + sign

	var result []byte
	err = util.DefaultAjaxPool.Post(&util.Ajax{
		Url:          url,
		Data:         body,
		ResponseData: &result,
	})
	if err != nil {
		return UmengStatusResult{}, err
	}

	var finalResult UmengStatusResult
	err = json.Unmarshal(result, &finalResult)
	if err != nil {
		return UmengStatusResult{}, err
	}
	return finalResult, nil
}

func (this *UmengSdk) getSign(method, url, body string) string {
	signStr := strings.ToUpper(method) + url + body + this.SecretKey
	return fmt.Sprintf("%x", md5.Sum([]byte(signStr)))
}
