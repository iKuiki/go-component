package utils_test

import (
	"testing"

	"github.com/iKuiki/go-component/utils"
	"github.com/stretchr/testify/assert"
)

func TestArrayReverse(t *testing.T) {
	a := []string{"a", "b", "c"}
	ar := utils.ArrayReverse(a)
	assert.Equal(t, []string{"c", "b", "a"}, ar)
}

func TestArrayIn(t *testing.T) {
	a := []string{"a", "b", "c"}
	i := utils.ArrayIn(a, "a")
	assert.Equal(t, 0, i)
	i = utils.ArrayIn(a, "d")
	assert.Equal(t, -1, i, "should not found")
}

func TestArrayUnique(t *testing.T) {
	a := []string{"a", "b", "b", "a", "c"}
	au := utils.ArrayUnique(a)
	assert.Equal(t, []string{"a", "b", "c"}, au)
}

func TestArrayDiff(t *testing.T) {
	a := []string{"a", "b", "c"}
	b := []string{"b", "d"}
	c := utils.ArrayDiff(a, b)
	assert.Equal(t, []string{"a", "c"}, c)
}

func TestArrayIntersect(t *testing.T) {
	a := []string{"a", "b", "c"}
	b := []string{"b", "d"}
	c := utils.ArrayIntersect(a, b)
	assert.Equal(t, []string{"b"}, c)
}

func TestArrayMerge(t *testing.T) {
	a := []string{"a", "b", "c"}
	b := []string{"b", "d"}
	c := utils.ArrayMerge(a, b)
	assert.Equal(t, []string{"a", "b", "c", "d"}, c)
}

func TestArraySort(t *testing.T) {
	a := []string{"a", "b", "c"}
	ar := utils.ArrayReverse(a)
	as := utils.ArraySort(ar)
	assert.Equal(t, a, as)

	ai := []int{1, 2, 3}
	air := utils.ArrayReverse(ai)
	ais := utils.ArraySort(air)
	assert.Equal(t, ai, ais)
}

func TestArrayShuffle(t *testing.T) {
	a := []string{"a", "b", "c", "d", "e", "f", "g"}
	as := utils.ArrayShuffle(a).([]string)
	assert.Equal(t, len(a), len(as), "乱序后长度应当相等")
	assert.Empty(t, utils.ArrayDiff(a, as), "两切片差集应当为空")
	assert.Empty(t, utils.ArrayDiff(as, a), "两切片差集应当为空")
	t.Log(a, as)
	assert.NotEqual(t, a, as, "乱序后应当不为同一数组")
}

func TestArraySlice(t *testing.T) {
	a := []string{"a", "b", "c", "d", "e", "f", "g"}
	as := utils.ArraySlice(a, 1, 2)
	assert.Equal(t, []string{"b"}, as)
	// 测试一些异常逻辑
	as = utils.ArraySlice(a, -1, 2)
	assert.Equal(t, []string{"a", "b"}, as)
	as = utils.ArraySlice(a, 3, 2)
	assert.Empty(t, as, "begin与end逆序")
	as = utils.ArraySlice(a, -2, -1)
	assert.Empty(t, as, "end小于0")
	as = utils.ArraySlice(a, 20)
	assert.Empty(t, as, "begin超过length")
	as = utils.ArraySlice(a, -2)
	assert.Equal(t, a, as)
}
