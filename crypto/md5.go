package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

// GetMd5String 将给定的字符串生成MD5字符串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
