package utils_test

import (
	"testing"

	"github.com/iKuiki/go-component/utils"
)

func TestQuerySort(t *testing.T) {
	a := []DemoBase{
		{ID: 2, Nickname: "b"},
		{ID: 1, Nickname: "a"},
		{ID: 3, Nickname: "d"},
		{ID: 3, Nickname: "c"},
	}
	as := utils.QuerySort(a, "ID")
	t.Log(as)
}
