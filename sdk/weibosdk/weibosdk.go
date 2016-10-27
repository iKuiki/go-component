package weibosdk

import (
	"github.com/yinhui87/go-component/encoding"
	"github.com/yinhui87/go-component/sdk/common"
	"github.com/yinhui87/go-component/util"
)

type WeiboSdk struct {
	AppId     string
	AppSecret string
}

func (this *WeiboSdk) api(method string, url string, query interface{}) ([]byte, error) {
	queryInfo, err := encoding.EncodeUrlQuery(query)
	if err != nil {
		return nil, err
	}
	url = "https://api.weibo.com" + url
	if len(queryInfo) != 0 {
		url += "?" + string(queryInfo)
	}
	var result []byte
	ajaxOption := &util.Ajax{
		Url:          url,
		ResponseData: &result,
	}
	if method == "GET" {
		err = util.DefaultAjaxPool.Get(ajaxOption)
	} else {
		err = util.DefaultAjaxPool.Post(ajaxOption)
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (this *WeiboSdk) apiUrl(method string, url string, query interface{}, responseData interface{}) error {
	result, err := this.api(method, url, query)
	if err != nil {
		return err
	}
	var sdkErr common.SdkError
	_, err = encoding.DecodeJsonp(result, &sdkErr)
	if err == nil && sdkErr.Code != 0 {
		return &sdkErr
	}
	err = encoding.DecodeJson(result, responseData)
	if err != nil {
		return err
	}
	return nil
}

func (this *WeiboSdk) GetOauthAccessToken(callback string, code string) (WeiboSdkOauthAccessToken, error) {
	var result WeiboSdkOauthAccessToken
	err := this.apiUrl("POST", "/oauth2/access_token", map[string]interface{}{
		"grant_type":    "authorization_code",
		"client_id":     this.AppId,
		"redirect_uri":  callback,
		"client_secret": this.AppSecret,
		"code":          code,
	}, &result)
	if err != nil {
		return WeiboSdkOauthAccessToken{}, err
	}
	return result, nil
}

//用户接口
func (this *WeiboSdk) GetUserInfo(accessToken, uid string) (WeiboSdkUserInfo, error) {
	result := WeiboSdkUserInfo{}
	err := this.apiUrl("GET", "/2/users/show.json", map[string]string{
		"access_token": accessToken,
		"uid":          uid,
	}, &result)
	if err != nil {
		return WeiboSdkUserInfo{}, err
	}
	return result, nil
}
