package alidayusdk

import (
	"encoding/json"
	"errors"
	"github.com/northbright/alidayu"
	"strings"
)

// AlidayuConfig 阿里大鱼sdk配置
type AlidayuConfig struct {
	SignName string
	// TemplateID string
	// TemplateParam map[string]string
}

// AlidayuSmsCommonTemplateParam 阿里大鱼短信模版
type AlidayuSmsCommonTemplateParam struct {
	Code    string `json:"code"`
	Product string `json:"product"`
}

// AlidayuSdk 阿里大鱼sdk
type AlidayuSdk struct {
	client *alidayu.Client
	config AlidayuConfig
}

// NewAlidayu 新建阿里大鱼sdk
func NewAlidayu(appKey, appSecret string, config AlidayuConfig) (dayu *AlidayuSdk, err error) {
	dayu = &AlidayuSdk{
		client: &alidayu.Client{AppKey: appKey, AppSecret: appSecret},
		config: config,
	}
	return dayu, nil
}

// SendMsg 发送信息
func (sdk *AlidayuSdk) SendMsg(mobile string, templateID string, templateParam interface{}) (err error) {
	return sdk.SendMultiMsg([]string{mobile}, templateID, templateParam)
}

// SendMultiMsg 批量发送信息
func (sdk *AlidayuSdk) SendMultiMsg(mobiles []string, templateID string, templateParam interface{}) (err error) {
	for _, v := range mobiles {
		if len(v) != 11 {
			return errors.New("Mobile " + v + " length not equal 11")
		}
	}
	smsParam, e := json.Marshal(&templateParam)
	if e != nil {
		return errors.New("Marshal templateParam to json error: " + e.Error())
	}
	params := map[string]string{
		"method":             "alibaba.aliqin.fc.sms.num.send",
		"sms_type":           "normal",
		"sms_free_sign_name": sdk.config.SignName,
		"sms_template_code":  templateID,
		"sms_param":          string(smsParam),
		"rec_num":            strings.Join(mobiles, ","),
	}
	success, result, e := sdk.client.Exec(params)
	if e != nil {
		return errors.New("SendMsg Exec error: " + e.Error())
	}
	if !success {
		return errors.New("SendMsg Fail: " + result)
	}
	return nil
}
