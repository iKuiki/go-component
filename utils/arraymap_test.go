package utils_test

import (
	"testing"

	"github.com/iKuiki/go-component/utils"
	"github.com/stretchr/testify/assert"
)

func TestArrayToMap(t *testing.T) {
	a := DemoBase{
		ID:       1,
		Nickname: "hello",
	}
	m := utils.ArrayToMap(a, "json")
	assert.Equal(t, map[string]interface{}{
		"iD":       1,
		"nickname": "hello",
	}, m)

	aa := []DemoBase{a}
	ma := utils.ArrayToMap(aa, "json")
	assert.Equal(t, []interface{}{
		map[string]interface{}{
			"iD":       1,
			"nickname": "hello",
		},
	}, ma)

	am := map[string]DemoBase{
		"a": a,
	}
	mm := utils.ArrayToMap(am, "json")
	assert.Equal(t, map[string]interface{}{
		"a": map[string]interface{}{
			"iD":       1,
			"nickname": "hello",
		},
	}, mm)
}

func TestMapToArray(t *testing.T) {
	m := map[string]interface{}{
		"iD":       1,
		"nickname": "hello",
	}
	var a DemoBase
	utils.MapToArray(m, &a, "json")
	assert.Equal(t, DemoBase{
		ID:       1,
		Nickname: "hello",
	}, a)

	ma := []interface{}{
		map[string]interface{}{
			"iD":       1,
			"nickname": "hello",
		},
	}
	var aa []DemoBase
	utils.MapToArray(ma, &aa, "json")
	assert.Equal(t, []DemoBase{a}, aa)

	mm := map[string]interface{}{
		"a": map[string]interface{}{
			"iD":       1,
			"nickname": "hello",
		},
	}
	var am map[string]DemoBase
	utils.MapToArray(mm, &am, "json")
	assert.Equal(t, map[string]DemoBase{
		"a": a,
	}, am)
}
