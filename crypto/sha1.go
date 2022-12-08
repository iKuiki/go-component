package crypto

import (
	"crypto/sha1"
	"encoding/hex"
)

// Sha1 将给定字符串通过sha1加密
func Sha1(data string) string {
	hash := sha1.New()
	hash.Write([]byte(data))
	etag := hash.Sum(nil)
	etagString := hex.EncodeToString(etag)
	return etagString
}
