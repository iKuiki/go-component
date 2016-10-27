package youdaosdk

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/httplib"
	"strconv"
)

const (
	YD_SERVER_API_URL = "http://fanyi.youdao.com"
	YD_PATH_TRANSLATE = "/openapi.do"
)

var errorCodeList map[int64]string

func init() {
	errorCodeList = map[int64]string{
		0:  "正常",
		20: "要翻译的文本过长",
		30: "无法进行有效的翻译",
		40: "不支持的语言类型",
		50: "无效的key",
		60: "无词典结果，仅在获取词典结果生效",
	}
}

type YoudaoSdk struct {
	apiUrl    string
	appKey    string
	appSecret string
}

//初始化YoudaoSdk
func NewYoudaoSdk(appKey, appSecret string) (youdaoSdk *YoudaoSdk, initError error) {
	if appKey == "" {
		return nil, errors.New("appKey不能为空！")
	} else if appSecret == "" {
		return nil, errors.New("appSecret不能为空！")
	}
	server := &YoudaoSdk{
		appKey:    appKey,
		appSecret: appSecret,
		apiUrl:    YD_SERVER_API_URL,
	}
	return server, nil
}

func (ydSdk *YoudaoSdk) Translate(text string) (result []string, err error) {
	req := httplib.Post(ydSdk.apiUrl + YD_PATH_TRANSLATE)
	req.Param("q", text)
	fillHeader(req, ydSdk)
	byteData, err := req.Bytes()
	if err != nil {
		return nil, errors.New("request to youdao server error: " + err.Error())
	}
	var translateRet TranslateResult
	err = json.Unmarshal(byteData, &translateRet)
	if err != nil {
		return nil, errors.New("respond parse as json error: " + err.Error())
	}
	if translateRet.ErrorCode != 0 {
		if errStr, ok := errorCodeList[translateRet.ErrorCode]; !ok {
			return nil, errors.New("unknow errorCode from youdao server: " + strconv.FormatInt(translateRet.ErrorCode, 10))
		} else {
			return nil, errors.New("youdao server error: " + errStr)
		}
	}
	return translateRet.Translation, nil
}
