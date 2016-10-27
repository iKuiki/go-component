package umengsdk

type AnaymonusMap map[string]string

type UmengCommon struct {
	Appkey         string      `json:"appkey"`
	Timestamp      string      `json:"timestamp"`
	Type           string      `json:"type"`
	DeviceTokens   string      `json:"device_tokens"`
	AliasType      string      `json:"alias_type"`
	Alias          string      `json:"alias"`
	FileId         string      `json:"file_id"`
	Filter         interface{} `json:"filter"`
	ProductionMode string      `json:"production_mode"`
	Description    string      `json:"description"`
	ThirdpartyId   string      `json:"thirdparty_id"`
}

type UmengAndroidPayloadBody struct {
	Ticker      string `json:"ticker"`
	Title       string `json:"title"`
	Text        string `json:"text"`
	Icon        string `json:"icon"`
	LargeIcon   string `json:"largeIcon"`
	Img         string `json:"img"`
	Sound       string `json:"sound"`
	BuilderId   string `json:"builder_id"`
	PlayVibrate bool   `json:"play_vibrate"`
	PlayLights  bool   `json:"play_lights"`
	PlaySound   bool   `json:"play_sound"`
	AfterOpen   string `json:"after_open"`
	Url         string `json:"url"`
	Activity    string `json:"activity"`
	Custom      string `json:"custom"`
}

type UmengAndroidPayload struct {
	DisplayType string                  `json:"display_type"`
	Body        UmengAndroidPayloadBody `json:"body"`
	Extra       map[string]string       `json:"extra"`
}

type UmengAndroidPolicy struct {
	StartTime  string `json:"start_time"`
	ExpireTime string `json:"expire_time"`
	//Max_send_num int    `json:"max_send_num"`
	OutBizNo string `json:"out_biz_no"`
}

type UmengAndroid struct {
	UmengCommon
	Payload UmengAndroidPayload `json:"payload"`
	Policy  UmengAndroidPolicy  `json:"policy"`
}

type UmengIOSPayloadAps struct {
	Alert            string `json:"alert"`
	Badge            int    `json:"badge"`
	Sound            string `json:"sound"`
	ContentAvailable string `json:"content-available"`
	Category         string `json:"category"`
}

type UmengIOSPayload struct {
	Aps UmengIOSPayloadAps `json:"aps"`
	AnaymonusMap
}

type UmengIOSPolicy struct {
	StartTime  string `json:"start_time"`
	ExpireTime string `json:"expire_time"`
	//Max_send_num int    `json:"max_send_num"`
}

type UmengIOS struct {
	UmengCommon
	Payload UmengIOSPayload `json:"payload"`
	Policy  UmengIOSPolicy  `json:"policy"`
}

type UmengResult struct {
	Ret  string `json:"ret"`
	Data struct {
		MsgId        string `json:"msg_id"`
		TaskId       string `json:"task_id"`
		ErrorCode    string `json:"error_code"`
		ThirdpartyId string `json:"thirdparty_id"`
	}
}

type UmengStatus struct {
	Appkey    string `json:"appkey"`
	Timestamp string `json:"timestamp"`
	TaskId    string `json:"task_id"`
}

type UmengStatusResult struct {
	Ret  string `json:"ret"`
	Data struct {
		TaskId string `json:"task_id"`
		Status int    `json:"status"` // 消息状态: 0-排队中, 1-发送中，2-发送完成，3-发送失败，4-消息被撤销，
		// 5-消息过期, 6-筛选结果为空，7-定时任务尚未开始处理
		TotalCount   int `json:"total_count"`   // 消息总数
		AcceptCount  int `json:"accept_count"`  // 消息受理数
		SentCount    int `json:"sent_count"`    // 消息实际发送数
		OpenCount    int `json:"open_count"`    //打开数
		DismissCount int `json:"dismiss_count"` //忽略数

		ErrorCode string `json:"error_code"`
	}
}

type UmengFile struct {
	Appkey    string `json:"appkey"`
	Timestamp string `json:"timestamp"`
	Content   string `json:"content"`
}

type UmengFileResult struct {
	Ret  string `json:"ret"`
	Data struct {
		FileId string `json:"file_id"`
	}
}
