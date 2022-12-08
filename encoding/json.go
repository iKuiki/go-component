package encoding

import (
	"encoding/json"

	"github.com/iKuiki/go-component/utils"
)

// EncodeJSON 将给定数据转为json
// 通过Map，先转为Map后再用json Marshal
func EncodeJSON(data interface{}) ([]byte, error) {
	changeValue := utils.ArrayToMap(data, "json")
	return json.Marshal(changeValue)
}

// DecodeJSON 解析json到给定结构
// 通过Map，先将json解析到map后再解析到结构
func DecodeJSON(data []byte, value interface{}) error {
	var valueDynamic interface{}
	err := json.Unmarshal(data, &valueDynamic)
	if err != nil {
		return err
	}
	return utils.MapToArray(valueDynamic, value, "json")
}
