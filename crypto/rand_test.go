package crypto_test

import (
	"testing"

	"github.com/iKuiki/go-component/crypto"
	"github.com/stretchr/testify/assert"
)

func TestRandString(t *testing.T) {
	sMap := make(map[string]bool)
	for i := 0; i < 100; i++ {
		sMap[crypto.RandString(10)] = true
	}
	assert.Len(t, sMap, 100)
}

func TestRandDigit(t *testing.T) {
	sMap := make(map[string]bool)
	for i := 0; i < 100; i++ {
		sMap[crypto.RandDigit(10)] = true
	}
	assert.Len(t, sMap, 100)
}
