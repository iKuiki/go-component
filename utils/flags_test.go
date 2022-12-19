package utils_test

import (
	"strings"
	"testing"

	"github.com/iKuiki/go-component/utils"
	"github.com/stretchr/testify/assert"
)

func TestSplitCommandFlag(t *testing.T) {
	args := "-e pro migrate -exec grant" // 虚拟一个传入参数
	mainFlags, subCommand, subArgs := utils.SplitCommandFlag(strings.Split(args, " "))
	assert.Equal(t, []string{
		"-e",
		"pro",
	}, mainFlags)
	assert.Equal(t, "migrate", subCommand)
	assert.Equal(t, []string{
		"-exec",
		"grant",
	}, subArgs)
}
