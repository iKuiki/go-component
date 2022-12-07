package utils_test

import (
	"testing"

	"github.com/iKuiki/go-component/utils"
)

func TestQuerySort(t *testing.T) {
	a := []struct {
		ID   int
		Name string
	}{
		{ID: 2, Name: "b"},
		{ID: 1, Name: "a"},
		{ID: 3, Name: "c"},
	}
	as := utils.QuerySort(a, "ID")
	t.Log(as)
}
