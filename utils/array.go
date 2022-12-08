package utils

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
)

// ArrayReverse 切片逆序
// 将传入的切片逆序后输出
func ArrayReverse(data interface{}) interface{} {
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type()
	dataLen := dataValue.Len()
	result := reflect.MakeSlice(dataType, dataLen, dataLen)

	for i := 0; i != dataLen; i++ {
		result.Index(dataLen - i - 1).Set(dataValue.Index(i))
	}
	return result.Interface()
}

// ArrayIn 检查给定元素在目标切片的什么位置
// 如果未找到则返回-1
func ArrayIn(arrayData interface{}, findData interface{}) int {
	var findIndex int
	findIndex = -1
	arrayDataValue := reflect.ValueOf(arrayData)
	arrayDataValueLen := arrayDataValue.Len()
	for i := 0; i != arrayDataValueLen; i++ {
		singleArrayDataValue := arrayDataValue.Index(i).Interface()
		if singleArrayDataValue == findData {
			findIndex = i
			break
		}
	}
	return findIndex
}

// ArrayExist 返回给定的元素在切片中是否存在
func ArrayExist(arrayData interface{}, findData interface{}) bool {
	return ArrayIn(arrayData, findData) != -1
}

// ArrayUnique 将给定切片中重复元素剔除
// 保证切片中元素唯一
func ArrayUnique(arrayData interface{}) interface{} {
	arrayValue := reflect.ValueOf(arrayData)
	arrayType := arrayValue.Type()
	arrayLen := arrayValue.Len()

	result := reflect.MakeSlice(arrayType, 0, 0)
	resultTemp := map[interface{}]bool{}

	for i := 0; i != arrayLen; i++ {
		singleArrayDataValue := arrayValue.Index(i)
		singleArrayDataValueInterface := singleArrayDataValue.Interface()
		_, isExist := resultTemp[singleArrayDataValueInterface]
		if isExist == true {
			continue
		}
		resultTemp[singleArrayDataValueInterface] = true
		result = reflect.Append(result, singleArrayDataValue)
	}
	return result.Interface()
}

// 将slice转为map
// 转换后slice的item为map的key
func sliceToMap(arrayData []interface{}) map[interface{}]bool {
	result := map[interface{}]bool{}

	for _, singleArray := range arrayData {
		arrayValue := reflect.ValueOf(singleArray)
		arrayLen := arrayValue.Len()

		for i := 0; i != arrayLen; i++ {
			result[arrayValue.Index(i).Interface()] = true
		}
	}

	return result
}

// ArrayDiff 计算切片的差集
// 对给定的切片A，剔除B、C。。。中包含的元素后输出
func ArrayDiff(arrayData interface{}, arrayData2 interface{}, arrayOther ...interface{}) interface{} {
	arrayOther = append([]interface{}{arrayData2}, arrayOther...)
	arrayOtherMap := sliceToMap(arrayOther)

	arrayValue := reflect.ValueOf(arrayData)
	arrayType := arrayValue.Type()
	arrayLen := arrayValue.Len()
	result := reflect.MakeSlice(arrayType, 0, 0)

	for i := 0; i != arrayLen; i++ {
		singleArrayDataValue := arrayValue.Index(i)
		singleArrayDataValueInterface := singleArrayDataValue.Interface()

		_, isExist := arrayOtherMap[singleArrayDataValueInterface]
		if isExist == true {
			continue
		}
		result = reflect.Append(result, singleArrayDataValue)
		arrayOtherMap[singleArrayDataValueInterface] = true
	}

	return result.Interface()
}

// ArrayIntersect 对多个切片计算交集
func ArrayIntersect(arrayData interface{}, arrayData2 interface{}, arrayOther ...interface{}) interface{} {
	arrayOther = append([]interface{}{arrayData2}, arrayOther...)
	arrayOtherMap := sliceToMap(arrayOther)

	arrayValue := reflect.ValueOf(arrayData)
	arrayType := arrayValue.Type()
	arrayLen := arrayValue.Len()
	result := reflect.MakeSlice(arrayType, 0, 0)

	for i := 0; i != arrayLen; i++ {
		singleArrayDataValue := arrayValue.Index(i)
		singleArrayDataValueInterface := singleArrayDataValue.Interface()

		isFirst, isExist := arrayOtherMap[singleArrayDataValueInterface]
		if isExist == false || isFirst == false {
			continue
		}
		result = reflect.Append(result, singleArrayDataValue)
		arrayOtherMap[singleArrayDataValueInterface] = false
	}

	return result.Interface()
}

// ArrayMerge 对多个切片计算并集
func ArrayMerge(arrayData interface{}, arrayData2 interface{}, arrayOther ...interface{}) interface{} {
	arrayOther = append([]interface{}{arrayData2}, arrayOther...)
	arrayOtherMap := sliceToMap(arrayOther)

	arrayValue := reflect.ValueOf(arrayData)
	arrayType := arrayValue.Type()
	arrayLen := arrayValue.Len()
	result := reflect.MakeSlice(arrayType, 0, 0)

	for i := 0; i != arrayLen; i++ {
		singleArrayDataValue := arrayValue.Index(i)
		singleArrayDataValueInterface := singleArrayDataValue.Interface()

		isFirst, isExist := arrayOtherMap[singleArrayDataValueInterface]
		if isExist == true && isFirst == false {
			continue
		}
		result = reflect.Append(result, singleArrayDataValue)
		arrayOtherMap[singleArrayDataValueInterface] = false
	}

	for single, isFirst := range arrayOtherMap {
		if isFirst == false {
			continue
		}
		result = reflect.Append(result, reflect.ValueOf(single))
	}

	return result.Interface()
}

// ArraySort 对切片进行升序排序
// 只允许传入[]int与[]string类型
func ArraySort(data interface{}) interface{} {
	//建立一份拷贝数据
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type()
	dataValueLen := dataValue.Len()

	dataResult := reflect.MakeSlice(dataType, dataValueLen, dataValueLen)
	reflect.Copy(dataResult, dataValue)

	//排序
	dataElemType := dataType.Elem()
	var result interface{}
	if dataElemType.Kind() == reflect.Int {
		intArray := dataResult.Interface().([]int)
		sort.Sort(sort.IntSlice(intArray))
		result = intArray
	} else if dataElemType.Kind() == reflect.String {
		stringArray := dataResult.Interface().([]string)
		sort.Sort(sort.StringSlice(stringArray))
		result = stringArray
	} else {
		panic("invalid sort type " + fmt.Sprintf("%v", dataElemType))
	}
	return result
}

// ArrayShuffle 将传入切片乱序后输出
func ArrayShuffle(data interface{}) interface{} {
	//建立一份拷贝数据
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type()
	dataValueLen := dataValue.Len()

	dataResult := reflect.MakeSlice(dataType, dataValueLen, dataValueLen)

	//打乱
	perm := rand.Perm(dataValueLen)
	for index, newIndex := range perm {
		dataResult.Index(newIndex).Set(dataValue.Index(index))
	}
	return dataResult.Interface()
}

// ArraySlice 获取切片中指定位置的一段数据
// 从index为begin的位置开始截取，包含begin
// 到index为end的位置停止，不含end
// begin与end做了安全校验
func ArraySlice(data interface{}, beginIndex int, endIndexArray ...int) interface{} {
	//建立一份拷贝数据
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type()
	dataValueLen := dataValue.Len()

	//计算size
	endIndex := 0
	if len(endIndexArray) >= 1 {
		endIndex = endIndexArray[0]
	} else {
		endIndex = dataValueLen
	}
	size := 0
	// 几个特殊情况，都不返回数据
	if beginIndex >= endIndex {
		// index逆向，不返回数据
		size = 0
	} else if endIndex <= 0 {
		// 结束位置越左界，不返回数据
		size = 0
	} else if beginIndex >= dataValueLen {
		// 开始位置越右界，不返回数据
		size = 0
	} else {
		// 有交集，则调整边界为有效值
		if beginIndex <= 0 {
			beginIndex = 0
		}
		if endIndex >= dataValueLen {
			endIndex = dataValueLen
		}
		size = endIndex - beginIndex
	}

	//拷贝
	dataResult := reflect.MakeSlice(dataType, size, size)
	for i := 0; i != size; i++ {
		dataResult.Index(i).Set(dataValue.Index(i + beginIndex))
	}
	return dataResult.Interface()
}
