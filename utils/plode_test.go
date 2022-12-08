package utils_test

import (
	"testing"

	"github.com/iKuiki/go-component/utils"
	"github.com/stretchr/testify/assert"
)

func TestExplode(t *testing.T) {
	s := "a,b, c"
	a := utils.Explode(s, ",")
	assert.Equal(t, []string{"a", "b", "c"}, a)

	s2 := utils.Implode(a, ",")
	assert.Equal(t, "a,b,c", s2)
}

func TestExplodeInt(t *testing.T) {
	s := "1,2 , 3"
	a := utils.ExplodeInt(s, ",")
	assert.Equal(t, []int{1, 2, 3}, a)

	s2 := utils.ImplodeInt(a, ",")
	assert.Equal(t, "1,2,3", s2)
}
