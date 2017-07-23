package alidayusdk

import (
	"encoding/json"
	"errors"
	"github.com/northbright/alidayu"
	"strings"
)

type AlidayuConfig struct {
	SignName string
	// TemplateId string
	// TemplateParam map[string]string
}

type AlidayuSmsCommonTemplateParam struct {
	Code    string `json:"code"`
	Product string `json:"product"`
}

type AlidayuSdk struct {
	client *alidayu.Client
	config AlidayuConfig
}

func NewAlidayu(appKey, appSecret string, config AlidayuConfig) (dayu *AlidayuSdk, err error) {
	dayu = &AlidayuSdk{
		client: &alidayu.Client{AppKey: appKey, AppSecret: appSecret},
		config: config,
	}
	return dayu, nil
}

func (this *AlidayuSdk) SendMsg(mobile string, templateId string, templateParam interface{}) (err error) {
	return this.SendMultiMsg([]string{mobile}, templateId, templateParam)
}

func (this *AlidayuSdk) SendMultiMsg(mobiles []string, templateId string, templateParam interface{}) (err error) {
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
		"sms_free_sign_name": this.config.SignName,
		"sms_template_code":  templateId,
		"sms_param":          string(smsParam),
		"rec_num":            strings.Join(mobiles, ","),
	}
	success, result, e := this.client.Exec(params)
	if e != nil {
		return errors.New("SendMsg Exec error: " + e.Error())
	}
	if !success {
		return errors.New("SendMsg Fail: " + result)
	}
	return nil
}
