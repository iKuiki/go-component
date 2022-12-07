package utils

import (
	"fmt"
	"reflect"
	"sort"
)

// ArrayColumnSort 对传入的struct为元素的切片
// 根据给定的column对其进行排序
// columnNames支持类sql查询，如ID asc,Name desc
// columnNames也支持struct嵌套查询，如Name.First
func ArrayColumnSort(data interface{}, columnNames string) interface{} {
	return QuerySort(data, columnNames)
}

// ArrayColumnUnique 对传入的以struct为元素的切片
// 根据给定的column去重
// 传入的columnNames可以是多个字段，以,分隔
// 如果传入多个字段，则多个字段都相同的元素才会认为是重复元素
func ArrayColumnUnique(data interface{}, columnNames string) interface{} {
	return QueryDistinct(data, columnNames)
}

type arrayColumnMapInfo struct {
	Index   []int
	Type    reflect.Type
	MapType reflect.Type
}

// ArrayColumnKey 对传入的以struct为元素的切片
// 获取给定的列名下的字段重组为切片返回
func ArrayColumnKey(data interface{}, columnName string) interface{} {
	return QueryColumn(data, columnName)
}

// ArrayColumnMap 对传入的以struct为元素的切片
// 生成对元素取其指定字段为key，原struct为value的map后返回
func ArrayColumnMap(data interface{}, columnNames string) interface{} {
	//提取信息
	name := Explode(columnNames, ",")
	nameInfo := []arrayColumnMapInfo{}
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type().Elem()
	for _, singleName := range name {
		singleField, ok := getFieldByName(dataType, singleName)
		if !ok {
			panic(dataType.Name() + " struct has not field " + singleName)
		}
		nameInfo = append(nameInfo, arrayColumnMapInfo{
			Index: singleField.Index,
			Type:  singleField.Type,
		})
	}
	prevType := dataType
	for i := len(nameInfo) - 1; i >= 0; i-- {
		nameInfo[i].MapType = reflect.MapOf(
			nameInfo[i].Type,
			prevType,
		)
		prevType = nameInfo[i].MapType
	}

	//整合map
	result := reflect.MakeMap(nameInfo[0].MapType)
	dataLen := dataValue.Len()
	for i := 0; i != dataLen; i++ {
		singleValue := dataValue.Index(i)
		prevValue := result
		for singleNameIndex, singleNameInfo := range nameInfo {
			var nextValue reflect.Value
			singleField := singleValue.FieldByIndex(singleNameInfo.Index)
			nextValue = prevValue.MapIndex(singleField)
			if !nextValue.IsValid() {
				if singleNameIndex+1 < len(nameInfo) {
					nextValue = reflect.MakeMap(nameInfo[singleNameIndex+1].MapType)
				} else {
					nextValue = singleValue
				}
				prevValue.SetMapIndex(singleField, nextValue)
			}
			prevValue = nextValue
		}
	}
	return result.Interface()
}

// ArrayColumnTable 将传入的数据转换为表格(二维数组)的形式展示
// @param column map[string]string 列ID:列名
// @param data 存放数据的以struct为元素的切片
func ArrayColumnTable(column interface{}, data interface{}) [][]string {
	result := [][]string{}

	columnMap := column.(map[string]string)
	columnKeys, _ := MapKeyAndValue(column)
	columnKeysReal := columnKeys.([]string)
	columnValuesReal := []string{}
	sort.Sort(sort.StringSlice(columnKeysReal))
	for _, singleKey := range columnKeysReal {
		columnValuesReal = append(columnValuesReal, columnMap[singleKey])
	}
	result = append(result, columnValuesReal)

	dataValue := reflect.ValueOf(data)
	dataLen := dataValue.Len()
	for i := 0; i != dataLen; i++ {
		singleDataValue := dataValue.Index(i)
		singleDataStringValue := ArrayToMap(singleDataValue.Interface(), "json")
		singleDataStringValueData := reflect.ValueOf(singleDataStringValue)
		singleResult := []string{}
		for _, singleColumn := range columnKeysReal {
			singleResultString := ""
			singleValue := singleDataStringValueData.MapIndex(reflect.ValueOf(singleColumn))
			if singleValue.IsValid() == false {
				singleResultString = ""
			} else {
				singleResultString = fmt.Sprintf("%v", singleValue)
			}
			singleResult = append(
				singleResult,
				singleResultString,
			)
		}
		result = append(result, singleResult)
	}
	return result
}
