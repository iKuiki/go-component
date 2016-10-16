package crypto

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type HashKind uint

const (
	Invalid HashKind = iota
	PASSWORD_ALGO_BCRYPT
)

func PasswordHash(password string, algo HashKind) (string, error) {
	var result []byte
	var err error
	switch algo {
	case PASSWORD_ALGO_BCRYPT:
		result, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	default:
		return "", errors.New("invalid password hash algo")
	}
	if err != nil {
		return "", err
	}
	return string(result), err
}

func PasswordVerify(password string, hash string, algo HashKind) (result bool, err error) {
	switch algo {
	case PASSWORD_ALGO_BCRYPT:
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
