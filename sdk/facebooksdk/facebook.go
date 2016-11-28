package facebooksdk

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/httplib"
	"log"
	"net/url"
	"strconv"
	"time"
)

const (
	GetOauthAccessTokenURL = "https://graph.facebook.com/oauth/access_token"
	GetOauthUserInfoURL    = "https://graph.facebook.com/me/"
)

type FacebookSdk struct {
	AppId       string
	AppSecret   string
	RedirectURI string
}

func (this *FacebookSdk) GetOauthAccessToken(authcode string) (*FacebookAccessToken, error) {
	req := httplib.Get(GetOauthAccessTokenURL)
	req.Param("client_id", this.AppId)
	req.Param("client_secret", this.AppSecret)
	req.Param("code", authcode)
	req.Param("scope", "public_profile")
	byteData, err := req.Bytes()
	if err != nil {
		return nil, errors.New("Request to Facebook Server Error: " + err.Error())
	}
	log.Printf("Facebook GetOauthAccessToken Response: %s\n", string(byteData))
	qs, err := url.ParseQuery(string(byteData))
	if err != nil {
		return nil, errors.New("ParseQuery Facebook respond Error: " + err.Error())
	}
	accessToken := FacebookAccessToken{}
	if accessToken.AccessToken = qs.Get("access_token"); accessToken.AccessToken == "" {
		return nil, errors.New("Facebook respond Not contain AccessToken: " + string(byteData))
	}
	if expires := qs.Get("expires"); expires != "" {
		expireInSecs, err := strconv.ParseInt(expires, 10, 64)
		if err != nil {
			return nil, errors.New("ExpiresIn ParseInt " + expires + " Error: " + err.Error())
		}
		accessToken.ExpiresAt = time.Unix(expireInSecs+time.Now().Unix(), 0)
	}
	return &accessToken, nil
}

func (this *FacebookSdk) GetUserInfo(accessToken string) (*FacebookUserInfo, error) {
	req := httplib.Get(GetOauthUserInfoURL)
	req.Param("access_token", accessToken)
	byteData, err := req.Bytes()
	if err != nil {
		return nil, errors.New("Request to Facebook Server Error: " + err.Error())
	}
	userInfo := FacebookUserInfo{}
	err = json.Unmarshal(byteData, &userInfo)
	if err != nil {
		return nil, errors.New("Unmarshal Facebook respond Error: " + err.Error())
	}
	return &userInfo, nil
}
