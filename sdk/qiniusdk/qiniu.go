package qiniusdk

import (
	"errors"
	"github.com/google/uuid"
	"github.com/qiniu/api.v7/conf"
	"github.com/yinhui87/go-component/crypto"
	// "qiniupkg.com/api.v7/conf"
	"qiniupkg.com/api.v7/kodo"
	"time"
)

func NewQiniuSdk(accessKey, secretKey, bucket string, tokenExpire uint32, maxUploadSize int64) (qiniuSdk *QiniuSdk, err error) {
	if accessKey == "" {
		return nil, errors.New("AccessKey empty")
	}
	conf.ACCESS_KEY = accessKey
	conf.SECRET_KEY = secretKey
	qiniuSdk = &QiniuSdk{
		Client:        kodo.New(0, nil),
		Bucket:        bucket,
		TokenExpires:  tokenExpire,
		MaxUploadSize: maxUploadSize,
	}
	domains, err := qiniuSdk.GetDomains()
	if err != nil {
		return nil, errors.New("Get bucket domain fail: " + err.Error())
	}
	if len(domains) == 0 {
		return nil, errors.New("Domain list empty")
	}
	qiniuSdk.Domain = domains[0]
	return qiniuSdk, nil
}

type QiniuSdk struct {
	Client        *kodo.Client
	Bucket        string
	Domain        string
	TokenExpires  uint32
	MaxUploadSize int64
}

func (this *QiniuSdk) GenFilename(ext string) (filename string) {
	uu := uuid.New()
	filename = time.Now().Format("20060102-150405-") + crypto.GetMd5String(uu.String()) + ext
	return filename
}

func (this *QiniuSdk) GetUploadToken(key string, sizeLimit ...int64) string {
	// 设置上传的策略
	policy := &kodo.PutPolicy{
		Scope: this.Bucket + ":" + key,
		// 设置Token过期时间
		Expires: this.TokenExpires,
		// 不允许覆盖
		InsertOnly: 1,
	}
	if len(sizeLimit) == 0 {
		policy.FsizeLimit = this.MaxUploadSize
	} else {
		if sizeLimit[0] > 0 {
			policy.FsizeLimit = sizeLimit[0]
		}
	}
	// 生成一个上传token
	token := this.Client.MakeUptoken(policy)
	return token
}

func (this *QiniuSdk) MoveFile(keySrc, keyDest string) error {
	return this.Client.Bucket(this.Bucket).Move(nil, keySrc, keyDest)
}

func (this *QiniuSdk) GetFileStat(key string) (entry kodo.Entry, err error) {
	entry, err = this.Client.Bucket(this.Bucket).Stat(nil, key)
	return entry, err
}

func (this *QiniuSdk) GetDomains() (domains []string, err error) {
	err = this.Client.Call(nil, &domains, "GET", "http://api.qiniu.com/v6/domain/list?tbl="+this.Bucket)
	return domains, err
}

func (this *QiniuSdk) FetchFile(url, key string) (err error) {
	err = this.Client.Bucket(this.Bucket).Fetch(nil, key, url)
	return err
}

type ImageBaseInfo struct {
	Format     string `json:"format"`
	Width      uint32 `json:"width"`
	Height     uint32 `json:"height"`
	ColorModel string `json:"colorModel"`
}

func (this *QiniuSdk) GetImageBaseInfo(key string) (imageInfo ImageBaseInfo, err error) {
	err = this.Client.Call(nil, &imageInfo, "GET", "http://"+this.Domain+"/"+key+"?imageInfo")
	return
}
