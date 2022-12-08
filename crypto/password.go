package crypto

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// PasswordHashKind 密码加密类型
type PasswordHashKind uint

const (
	// PasswordHashKindInvalid 非法加密类型
	PasswordHashKindInvalid PasswordHashKind = iota
	// PasswordHashKindBcrypt bcrypt加密
	PasswordHashKindBcrypt
)

// PasswordHash 密码加密
func PasswordHash(password string, algo PasswordHashKind) (string, error) {
	var result []byte
	var err error
	switch algo {
	case PasswordHashKindBcrypt:
		result, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	default:
		return "", errors.New("invalid password hash algo")
	}
	if err != nil {
		return "", err
	}
	return string(result), err
}

// PasswordVerify 密码解密认证
func PasswordVerify(password string, hash string, algo PasswordHashKind) (result bool, err error) {
	switch algo {
	case PasswordHashKindBcrypt:
		if len(hash) <= 4 {
			return false, errors.New("invalid password hash format [" + hash + "]")
		}
		hashAlgo := hash[1:3]
		if hashAlgo != "2a" && hashAlgo != "2y" {
			return false, errors.New("invalid password hash algo [" + hashAlgo + "]")
		}
		fail := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		result = (fail == nil)
	}
	return result, err
}
