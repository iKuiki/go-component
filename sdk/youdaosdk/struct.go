package youdaosdk

type TranslateResult struct {
	ErrorCode   int64           `json:"errorCode"`
	Query       string          `json:"query"`
	Translation []string        `json:"translation"` // 有道翻译
	Basic       BasicDictResult `json:"basic"`
	Web         []WebDictResult `json:"web"`
}

type BasicDictResult struct { // 有道词典-基本词典
	Phonetic   string   `json:"phonetic"`
	UkPhonetic string   `json:"uk-phonetic"`
	UsPhonetic string   `json:"us-phonetic"`
	Explains   []string `json:"explains"`
}

type WebDictResult struct { // 有道词典-网络释义
	Key   string   `json:"key"`
	Value []string `json:"value"`
}
