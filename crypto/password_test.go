package crypto_test

import (
	"testing"

	"github.com/iKuiki/go-component/crypto"
	"github.com/stretchr/testify/assert"
)

func TestPasswordCrypto(t *testing.T) {
	password := "aabbcc"
	hash, err := crypto.PasswordHash(password, crypto.PasswordHashKindBcrypt)
	assert.NoError(t, err)
	t.Logf("PasswordHash: %s\nlen: %d\n", hash, len(hash))
	isRight, err := crypto.PasswordVerify(password, hash, crypto.PasswordHashKindBcrypt)
	assert.NoError(t, err)
	assert.True(t, isRight)
}
