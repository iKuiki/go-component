package utils_test

import (
	"testing"

	"github.com/iKuiki/go-component/utils"
	"github.com/stretchr/testify/assert"
)

type DemoBase struct {
	ID       int
	Nickname string
}

type DemoStruct struct {
	ID       int
	Nickname string
	Name     struct {
		First string
		Last  string
	}
}

func TestArrayColumnSort(t *testing.T) {
	a := []DemoStruct{
		{ID: 3, Nickname: "c", Name: struct {
			First string
			Last  string
		}{First: "c1", Last: "c2"}},

		{ID: 1, Nickname: "a", Name: struct {
			First string
			Last  string
		}{First: "a1", Last: "a2"}},

		{ID: 2, Nickname: "b", Name: struct {
			First string
			Last  string
		}{First: "b1", Last: "b2"}},
	}
	as := utils.ArrayColumnSort(a, "Name.First,ID")
	assert.Equal(t, []DemoStruct{
		{ID: 1, Nickname: "a", Name: struct {
			First string
			Last  string
		}{First: "a1", Last: "a2"}},

		{ID: 2, Nickname: "b", Name: struct {
			First string
			Last  string
		}{First: "b1", Last: "b2"}},

		{ID: 3, Nickname: "c", Name: struct {
			First string
			Last  string
		}{First: "c1", Last: "c2"}},
	}, as)
}

func TestArrayColumnUnique(t *testing.T) {
	a := []DemoStruct{
		{ID: 1, Nickname: "a"},
		{ID: 2, Nickname: "a"},
		{ID: 3, Nickname: "d"},
		{ID: 3, Nickname: "d"},
	}
	as := utils.ArrayColumnUnique(a, "Nickname")
	assert.Equal(t, []DemoStruct{
		{ID: 1, Nickname: "a"},
		{ID: 3, Nickname: "d"},
	}, as, "Nickname重复的字段会被过滤")

	as = utils.ArrayColumnUnique(a, "ID,Nickname")
	assert.Equal(t, []DemoStruct{
		{ID: 1, Nickname: "a"},
		{ID: 2, Nickname: "a"},
		{ID: 3, Nickname: "d"},
	}, as)
}

func TestArrayColumnKey(t *testing.T) {
	a := []DemoStruct{
		{ID: 1, Nickname: "a"},
		{ID: 2, Nickname: "a"},
		{ID: 3, Nickname: "d"},
	}
	ac := utils.ArrayColumnKey(a, "Nickname")
	assert.Equal(t, []string{"a", "a", "d"}, ac)
}

func TestArrayColumnMap(t *testing.T) {
	a := []DemoStruct{
		{ID: 1, Nickname: "a"},
		{ID: 2, Nickname: "a"},
		{ID: 3, Nickname: "d"},
	}
	am := utils.ArrayColumnMap(a, "Nickname")
	assert.Equal(t, map[string]DemoStruct{
		"a": {ID: 1, Nickname: "a"},
		"d": {ID: 3, Nickname: "d"},
	}, am)
}

func TestArrayColumnTable(t *testing.T) {
	column := map[string]string{
		"nickname": "昵称",
		"iD":       "ID",
	}
	data := []DemoStruct{
		{ID: 1, Nickname: "a"},
		{ID: 2, Nickname: "a"},
		{ID: 3, Nickname: "d"},
	}
	at := utils.ArrayColumnTable(column, data)
	t.Log(at)
}
