package kugousdk

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/httplib"
	"strconv"
)

const (
	KG_PATH_SEARCH = "http://mobilecdn.kugou.com/api/v3/search/song"
	KG_PATH_DETAIL = "http://m.kugou.com/app/i/getSongInfo.php"
)

type KugouSdk struct {
}

func NewKugouSdk() (kugouSdk *KugouSdk, err error) {
	sdk := &KugouSdk{}
	return sdk, nil
}

func (this *KugouSdk) Search(keyword string, page, size uint64) ([]SearchSongItem, error) {
	req := httplib.Get(KG_PATH_SEARCH)
	req.Param("page", strconv.FormatUint(page, 10))
	req.Param("pagesize", strconv.FormatUint(size, 10))
	req.Param("keyword", keyword)
	byteData, err := req.Bytes()
	if err != nil {
		return nil, errors.New("request to kugou server error: " + err.Error())
	}
	var searchResult SearchResult
	err = json.Unmarshal(byteData, &searchResult)
	if err != nil {
		return nil, errors.New("respond parse as json error: " + err.Error())
	}
	if searchResult.Errcode != 0 {
		return nil, errors.New("kugou server errorCode " + strconv.Itoa(searchResult.Errcode) + ": " + searchResult.Error)
	}
	return searchResult.Data.Info, nil
}

func (this *KugouSdk) Detail(hash string) (SongDetail, error) {
	req := httplib.Get(KG_PATH_DETAIL)
	req.Param("cmd", "playInfo")
	req.Param("hash", hash)
	byteData, err := req.Bytes()
	if err != nil {
		return SongDetail{}, errors.New("request to kugou server error: " + err.Error())
	}
	var detail DetailResult
	err = json.Unmarshal(byteData, &detail)
	if err != nil {
		return SongDetail{}, errors.New("respond parse as json error: " + err.Error())
	}
	if detail.Errcode != 0 {
		return SongDetail{}, errors.New("kugou server errorCode " + strconv.Itoa(detail.Errcode) + ": " + detail.Error)
	}
	ret := SongDetail{
		BitRate: detail.BitRate,
		Ctype:   detail.Ctype,
		ExtName: detail.ExtName,
		Extra: DetailExtra{
			One28filesize:   detail.Extra.One28filesize,
			One28hash:       detail.Extra.One28hash,
			Three20filesize: detail.Extra.Three20filesize,
			Three20hash:     detail.Extra.Three20hash,
			Sqfilesize:      detail.Extra.Sqfilesize,
			Sqhash:          detail.Extra.Sqhash,
		},
		FileHead:    detail.FileHead,
		FileName:    detail.FileName,
		FileSize:    detail.FileSize,
		Hash:        detail.Hash,
		Privilege:   detail.Privilege,
		Q:           detail.Q,
		ReqHash:     detail.ReqHash,
		SingerHead:  detail.SingerHead,
		Stype:       detail.Stype,
		TimeLength:  detail.TimeLength,
		TopicRemark: detail.TopicRemark,
		TopicURL:    detail.TopicURL,
		URL:         detail.URL,
	}
	return ret, nil
}
