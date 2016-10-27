package weibosdk

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
