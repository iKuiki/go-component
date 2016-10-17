package sdk

import (
	"tapi/helper/encoding"
	"tapi/helper/util"
)

type WeiboSdk struct {
	AppId     string
	AppSecret string
}

type WeiboSdkOauthAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   uint32 `json:"expires_in"`
	RemindIn    string `json:"remind_in"`
	Uid         string `json:"uid"`
}

type WeiboSdkUserInfo struct {
	Id         uint64 `json:"id"`
	Idstr      string `json:"idstr"`
	ScreenName string `json:"screen_name"`
	Name       string `json:"name"`
	// 用户无用信息暂不处理？
	ProfileImageUrl string `json:"profile_image_url"`
	Gender          string `json:"gender"`
	AvatarLarge     string `json:"avatar_large"`
	AvatarHd        string `json:"avatar_hd"`
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
	var sdkErr QqSdkError
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
