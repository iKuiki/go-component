package utils_test

import (
	"testing"

	"github.com/iKuiki/go-component/utils"
	"github.com/stretchr/testify/assert"
)

func TestMapKeyAndValue(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
	}
	k, v := utils.MapKeyAndValue(m)
	kA, ok := k.([]string)
	assert.True(t, ok)
	assert.Equal(t, []string{"a", "b"}, kA)
	vA, ok := v.([]int)
	assert.True(t, ok)
	assert.Equal(t, []int{1, 2}, vA)
}

func TestMapSafeDeleteItem(t *testing.T) {
	m := map[string]interface{}{
		"a": 1,
		"b": 2,
	}
	ok := utils.MapSafeDeleteItem(&m, "a")
	assert.True(t, ok)
	assert.Len(t, m, 1)
	ok = utils.MapSafeDeleteItem(&m, "a")
	assert.False(t, ok)
	ok = utils.MapSafeDeleteItem(&m, "c")
	assert.False(t, ok)
}
