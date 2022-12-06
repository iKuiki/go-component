package encoding

import (
	"encoding/json"

	"github.com/iKuiki/go-component/language"
)

func EncodeJson(data interface{}) ([]byte, error) {
	changeValue := language.ArrayToMap(data, "json")
	return json.Marshal(changeValue)
}

func DecodeJson(data []byte, value interface{}) error {
	var valueDynamic interface{}
	err := json.Unmarshal(data, &valueDynamic)
	if err != nil {
		return err
	}
	return language.MapToArray(valueDynamic, value, "json")
}
