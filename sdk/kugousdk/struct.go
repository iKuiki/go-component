package kugousdk

type SearchResult struct {
	Data struct {
		Correctiontip   string           `json:"correctiontip"`
		Correctiontype  int              `json:"correctiontype"`
		Forcecorrection int              `json:"forcecorrection"`
		Info            []SearchSongItem `json:"info"`
		Istag           int              `json:"istag"`
		Istagresult     int              `json:"istagresult"`
		Timestamp       uint32           `json:"timestamp"`
		Total           uint32           `json:"total"`
	} `json:"data"`
	Errcode int    `json:"errcode"`
	Error   string `json:"error"`
	Status  int    `json:"status"`
}

type SearchSongItem struct {
	Three20filesize  int                   `json:"320filesize"`
	Three20hash      string                `json:"320hash"`
	Three20privilege int                   `json:"320privilege"`
	Accompany        int                   `json:"Accompany"`
	AlbumID          string                `json:"album_id"`
	AlbumName        string                `json:"album_name"`
	Bitrate          int                   `json:"bitrate"`
	Duration         int                   `json:"duration"`
	Extname          string                `json:"extname"`
	Feetype          int                   `json:"feetype"`
	Filename         string                `json:"filename"`
	Filesize         int                   `json:"filesize"`
	Group            []SearchSongGroupItem `json:"group"`
	Hash             string                `json:"hash"`
	Isnew            int                   `json:"isnew"`
	M4afilesize      int                   `json:"m4afilesize"`
	Mvhash           string                `json:"mvhash"`
	Othername        string                `json:"othername"`
	Ownercount       int                   `json:"ownercount"`
	Privilege        int                   `json:"privilege"`
	Singername       string                `json:"singername"`
	Songname         string                `json:"songname"`
	Source           string                `json:"source"`
	Sourceid         int                   `json:"sourceid"`
	Sqfilesize       int                   `json:"sqfilesize"`
	Sqhash           string                `json:"sqhash"`
	Sqprivilege      int                   `json:"sqprivilege"`
	Srctype          int                   `json:"srctype"`
	Topic            string                `json:"topic"`
}

type SearchSongGroupItem struct {
	Three20filesize  int    `json:"320filesize"`
	Three20hash      string `json:"320hash"`
	Three20privilege int    `json:"320privilege"`
	Accompany        int    `json:"Accompany"`
	AlbumID          string `json:"album_id"`
	AlbumName        string `json:"album_name"`
	Bitrate          int    `json:"bitrate"`
	Duration         int    `json:"duration"`
	Extname          string `json:"extname"`
	Feetype          int    `json:"feetype"`
	Filename         string `json:"filename"`
	Filesize         int    `json:"filesize"`
	Hash             string `json:"hash"`
	Isnew            int    `json:"isnew"`
	M4afilesize      int    `json:"m4afilesize"`
	Mvhash           string `json:"mvhash"`
	Othername        string `json:"othername"`
	Ownercount       int    `json:"ownercount"`
	Privilege        int    `json:"privilege"`
	Singername       string `json:"singername"`
	Songname         string `json:"songname"`
	Source           string `json:"source"`
	Sourceid         int    `json:"sourceid"`
	Sqfilesize       int    `json:"sqfilesize"`
	Sqhash           string `json:"sqhash"`
	Sqprivilege      int    `json:"sqprivilege"`
	Srctype          int    `json:"srctype"`
	Topic            string `json:"topic"`
}

type DetailExtra struct {
	One28filesize   int    `json:"128filesize"`
	One28hash       string `json:"128hash"`
	Three20filesize int    `json:"320filesize"`
	Three20hash     string `json:"320hash"`
	Sqfilesize      int    `json:"sqfilesize"`
	Sqhash          string `json:"sqhash"`
}

type DetailResult struct {
	BitRate     int         `json:"bitRate"`
	Ctype       int         `json:"ctype"`
	Errcode     int         `json:"errcode"`
	Error       string      `json:"error"`
	ExtName     string      `json:"extName"`
	Extra       DetailExtra `json:"extra"`
	FileHead    int         `json:"fileHead"`
	FileName    string      `json:"fileName"`
	FileSize    int         `json:"fileSize"`
	Hash        string      `json:"hash"`
	Privilege   int         `json:"privilege"`
	Q           int         `json:"q"`
	ReqHash     string      `json:"req_hash"`
	SingerHead  string      `json:"singerHead"`
	Status      int         `json:"status"`
	Stype       int         `json:"stype"`
	TimeLength  int         `json:"timeLength"`
	TopicRemark string      `json:"topic_remark"`
	TopicURL    string      `json:"topic_url"`
	URL         string      `json:"url"`
}

type SongDetail struct {
	BitRate     int         `json:"bitRate"`
	Ctype       int         `json:"ctype"`
	ExtName     string      `json:"extName"`
	Extra       DetailExtra `json:"extra"`
	FileHead    int         `json:"fileHead"`
	FileName    string      `json:"fileName"`
	FileSize    int         `json:"fileSize"`
	Hash        string      `json:"hash"`
	Privilege   int         `json:"privilege"`
	Q           int         `json:"q"`
	ReqHash     string      `json:"req_hash"`
	SingerHead  string      `json:"singerHead"`
	Stype       int         `json:"stype"`
	TimeLength  int         `json:"timeLength"`
	TopicRemark string      `json:"topic_remark"`
	TopicURL    string      `json:"topic_url"`
	URL         string      `json:"url"`
}
