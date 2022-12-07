package utils

import "reflect"

// MapKeyAndValue 将map的key与value以数组返回
func MapKeyAndValue(data interface{}) (interface{}, interface{}) {
	//解析data
	dataType := reflect.TypeOf(data)
	if dataType.Kind() != reflect.Map {
		panic("need a map for MapKeyAndValue")
	}
	dataKeyType := dataType.Key()
	dataValueType := dataType.Elem()

	//合并数据
	dataKeySlice := reflect.MakeSlice(reflect.SliceOf(dataKeyType), 0, 0)
	dataValueSlice := reflect.MakeSlice(reflect.SliceOf(dataValueType), 0, 0)
	dataValue := reflect.ValueOf(data)
	for _, singleKey := range dataValue.MapKeys() {
		dataKeySlice = reflect.Append(dataKeySlice, singleKey)
		dataValueSlice = reflect.Append(dataValueSlice, dataValue.MapIndex(singleKey))
	}
	return dataKeySlice.Interface(), dataValueSlice.Interface()
}

// MapSafeDeleteItem 安全删除map下的元素(避免元素不存在导致报错)
func MapSafeDeleteItem(m *map[string]interface{}, field string) bool {
	if _, ok := (*m)[field]; ok {
		delete(*m, field)
		return true
	}
	return false
}
