package encoding

import (
	"encoding/base64"
)

// EncodeBase64 将给定数据编码为Base64
func EncodeBase64(in []byte) (string, error) {
	return base64.StdEncoding.EncodeToString(in), nil
}

// DecodeBase64 将给定base64解码为原始数据
func DecodeBase64(in string) ([]byte, error) {
	result, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return nil, err
	}
	return result, nil
}
