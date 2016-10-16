package crypto

import (
	"testing"
)

func TestPasswordCrypto(t *testing.T) {
	password := "aabbcc"
	hash, err := PasswordHash(password, PASSWORD_ALGO_BCRYPT)
	if err != nil {
		t.Fatalf("PasswordHash Error: %s\n", err.Error())
	}
	t.Logf("PasswordHash: %s\nlen: %d\n", hash, len(hash))
	isRight, err := PasswordVerify(password, hash, PASSWORD_ALGO_BCRYPT)
	if err != nil {
		t.Fatalf("PasswordVerify Error: %s\n", err.Error())
	}
	t.Logf("PasswordVerify Result: %t\n", isRight)
}
