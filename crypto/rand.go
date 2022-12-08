package crypto

import "crypto/rand"

var (
	randStr = []byte("0123456789abcdefghijklmnopqrstuvwxyz")
)

func generateRand(size int, targetLength byte) string {
	result := make([]byte, size)
	rand.Read(result)
	for singleIndex, singleByte := range result {
		result[singleIndex] = randStr[singleByte%targetLength]
	}
	return string(result)
}

// RandString 获取随机字符串
func RandString(size int) string {
	return generateRand(size, byte(len(randStr)))
}

// RandDigit 获取随机数字
func RandDigit(size int) string {
	return generateRand(size, 10)
}
