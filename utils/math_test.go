package utils_test

import (
	"testing"

	"github.com/iKuiki/go-component/utils"
	"github.com/stretchr/testify/assert"
)

func TestAbsInt(t *testing.T) {
	a := utils.AbsInt(1)
	assert.Equal(t, 1, a)
	a = utils.AbsInt(-1)
	assert.Equal(t, 1, a)
}
