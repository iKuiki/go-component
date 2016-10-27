package qiniusdk

import (
	"bytes"
	"errors"
	"net/url"
	"qiniupkg.com/api.v7/kodo"
)

const (
	maxSize = 128 * 1024 * 1024
)

type QiniuSdk struct {
	AccessKey string
	SecretKey string
}

func (this *QiniuSdk) UploadString(bucketName string, data []byte) (string, error) {
	//判断图片文件大小
	fsize := int64(len(data))
	if fsize > maxSize {
		return "", errors.New("上传图片太大！")
	}

	uploadReader := bytes.NewReader(data)
	cfg := &kodo.Config{
		AccessKey: this.AccessKey,
		SecretKey: this.SecretKey,
	}
	client := kodo.New(0, cfg)
	bucket := client.Bucket(bucketName)
	putRet := kodo.PutRet{}
	err := bucket.PutWithoutKey(nil, &putRet, uploadReader, fsize, &kodo.PutExtra{})
	if err != nil {
		return "", err
	}

	return putRet.Hash, nil
}

func (this *QiniuSdk) UploadFile(bucketName string, fileAddr string) (string, error) {
	cfg := &kodo.Config{
		AccessKey: this.AccessKey,
		SecretKey: this.SecretKey,
	}
	client := kodo.New(0, cfg)
	bucket := client.Bucket(bucketName)
	putRet := kodo.PutRet{}
	err := bucket.PutFileWithoutKey(nil, &putRet, fileAddr, &kodo.PutExtra{})
	if err != nil {
		return "", err
	}

	return putRet.Hash, nil
}

func (this *QiniuSdk) MoveFile(bucketName string, keySrc, keyDest string) error {
	cfg := &kodo.Config{
		AccessKey: this.AccessKey,
		SecretKey: this.SecretKey,
	}
	client := kodo.New(0, cfg)
	bucket := client.Bucket(bucketName)
	return bucket.Move(nil, keySrc, keyDest)
}

func (this *QiniuSdk) MakeBaseUrl(domain, key string) string {
	return kodo.MakeBaseUrl(domain, key)
}

func (this *QiniuSdk) GetUploadToken(bucketName string) (string, error) {
	cfg := &kodo.Config{
		AccessKey: this.AccessKey,
		SecretKey: this.SecretKey,
	}
	client := kodo.New(0, cfg)

	putPolicy := &kodo.PutPolicy{
		Scope:   bucketName,
		Expires: 3600,
	}

	token := client.MakeUptoken(putPolicy)

	return token, nil
}

func (this *QiniuSdk) GetDownloadUrl(inUrl string) (string, error) {
	urlStruct, err := url.Parse(inUrl)
	if err != nil {
		return "", err
	}
	domain := "http://" + urlStruct.Host
	key := urlStruct.Path

	cfg := &kodo.Config{
		AccessKey: this.AccessKey,
		SecretKey: this.SecretKey,
	}
	client := kodo.New(0, cfg)

	getPolicy := &kodo.GetPolicy{
		Expires: 3600,
	}
	baseUrl := kodo.MakeBaseUrl(domain, key)
	privateUrl := client.MakePrivateUrl(baseUrl, getPolicy)

	return privateUrl, nil
}
