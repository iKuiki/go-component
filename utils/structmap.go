package utils

import (
	"encoding/json"
	"reflect"
	"strings"
)

// StructToMap 将struct转为map
// @param model 传入的待转换的struct模型
// @param except 要忽略的字段
func StructToMap(model interface{}, except ...string) map[string]interface{} {
	ret := map[string]interface{}{}

	modelReflect := reflect.ValueOf(model)

	if modelReflect.Kind() == reflect.Ptr {
		modelReflect = modelReflect.Elem()
	}
	if modelReflect.Kind() == reflect.Invalid {
		return nil
	}
	modelRefType := modelReflect.Type()
	fieldsCount := modelReflect.NumField()

	var fieldData interface{}

	for i := 0; i < fieldsCount; i++ {
		var fieldName string
		if tag := modelRefType.Field(i).Tag.Get("json"); tag == "" {
			fieldName = modelRefType.Field(i).Name
		} else {
			fieldName = strings.Split(tag, ",")[0]
		}
		if ArrayIn(except, fieldName) > -1 {
			continue
		}

		field := modelReflect.Field(i)

		switch field.Kind() {
		case reflect.Struct:
			fallthrough
		case reflect.Ptr:
			if field.Interface() != nil {
				fieldData = StructToMap(field.Interface())
			}
		default:
			fieldData = field.Interface()
		}

		ret[fieldName] = fieldData
	}

	return ret
}

// StructToMapViaJSON 通过json将struct转为map
func StructToMapViaJSON(model interface{}) map[string]interface{} {
	data, err := json.Marshal(model)
	if err != nil {
		return nil
	}
	ret := map[string]interface{}{}
	if json.Unmarshal(data, &ret) != nil {
		return nil
	}
	return ret
}
