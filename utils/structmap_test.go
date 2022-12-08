package utils_test

import (
	"testing"

	"github.com/iKuiki/go-component/utils"
	"github.com/stretchr/testify/assert"
)

func TestStructToMap(t *testing.T) {
	a := DemoStruct{
		ID:       3,
		Nickname: "hello",
	}
	m := utils.StructToMap(a, "Name")
	assert.Equal(t, map[string]interface{}{
		"ID":       3,
		"Nickname": "hello",
	}, m)
}

func TestStructToMapViaJSON(t *testing.T) {
	a := DemoBase{
		ID:       3,
		Nickname: "hello",
	}
	m := utils.StructToMapViaJSON(a)
	assert.Equal(t, map[string]interface{}{
		"ID":       3.0, // 通过json时此处为float64
		"Nickname": "hello",
	}, m)
}
