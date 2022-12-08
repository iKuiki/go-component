package utils

import (
	"strconv"
	"strings"
)

// Explode 将字符串按分隔符分割为字符串切片
// 与stirngs.Split相比，多了Trim空格的逻辑
// 如果Trim空格后字符串为空，则不会被添加到结果集
func Explode(input string, separator string) []string {
	dataResult := strings.Split(input, separator)
	dataResultNew := []string{}
	for _, singleResult := range dataResult {
		singleResult = strings.Trim(singleResult, " ")
		if len(singleResult) == 0 {
			continue
		}
		dataResultNew = append(dataResultNew, singleResult)
	}
	return dataResultNew
}

// Implode Explode的反向操作，与strings.Join相同
func Implode(data []string, separator string) string {
	return strings.Join(data, separator)
}

// ExplodeInt 将字符串按分隔符分割为int切片
// 要求传入的必须是可以识别为int的数据，否则panic
func ExplodeInt(input string, separator string) []int {
	dataResult := strings.Split(input, separator)
	dataResultNew := []int{}
	for _, singleResult := range dataResult {
		singleResult = strings.Trim(singleResult, " ")
		if len(singleResult) == 0 {
			continue
		}
		singleResultInt, err := strconv.Atoi(singleResult)
		if err != nil {
			panic(err)
		}
		dataResultNew = append(dataResultNew, singleResultInt)
	}
	return dataResultNew
}

// ImplodeInt ExplodeInt的反向操作
// 将int数组合并为指定分隔符的字符串
func ImplodeInt(data []int, separator string) string {
	result := []string{}
	for _, singleData := range data {
		result = append(result, strconv.Itoa(singleData))
	}
	return strings.Join(result, separator)
}
