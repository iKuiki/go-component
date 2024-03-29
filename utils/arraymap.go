package utils

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// 名称映射
// 将传入的字符串的第一位改为小写后返回
func nameMapper(name string) string {
	return strings.ToLower(name[0:1]) + name[1:]
}

// 将两个map合并为一个
// 如果两个map都有同一个key，则以前者的为准
func combileMap(result map[string]interface{}, singleResultMap reflect.Value) {
	singleResultMapType := singleResultMap.Type()
	if singleResultMapType.Kind() != reflect.Map {
		return
	}
	singleResultMapKeys := singleResultMap.MapKeys()
	for _, singleKey := range singleResultMapKeys {
		singleResultKey := fmt.Sprintf("%v", singleKey)
		_, isExist := result[singleResultKey]
		if isExist {
			continue
		}
		result[singleResultKey] = singleResultMap.MapIndex(singleKey).Interface()
	}
}

// 映射时struct的结构信息
type arrayMappingStructInfo struct {
	name      string
	omitempty bool
	anonymous bool
	canRead   bool
	canWrite  bool
	index     []int
}

// 映射字段信息
type arrayMappingInfo struct {
	kind       int
	isTimeType bool
	field      []arrayMappingStructInfo
}

// 映射信息缓存
var arrayMappingInfoMap struct {
	mutex sync.RWMutex
	data  map[string]map[reflect.Type]arrayMappingInfo
}

func init() {
	arrayMappingInfoMap.data = map[string]map[reflect.Type]arrayMappingInfo{}
	var mm struct {
		Test interface{}
	}
	interfaceType = reflect.TypeOf(mm).Field(0).Type
}

var interfaceType reflect.Type

// 获取指定对象的结构信息的实现
// @param tag 用来读取json相关定义时，读取的tag
func getDataTagInfoInner(dataType reflect.Type, tag string) arrayMappingInfo {
	dataTypeKind := GetTypeKind(dataType)
	result := arrayMappingInfo{}
	result.kind = dataTypeKind
	if dataTypeKind == TypeKindStruct {
		if dataType == reflect.TypeOf(time.Time{}) {
			//时间类型
			result.isTimeType = true
		} else {
			//结构体类型
			result.isTimeType = false
			anonymousField := []arrayMappingStructInfo{}
			noanonymousField := []arrayMappingStructInfo{}
			for i := 0; i != dataType.NumField(); i++ {
				singleDataType := dataType.Field(i)
				if singleDataType.PkgPath != "" && singleDataType.Anonymous == false {
					continue
				}
				var singleName string
				var omitempty bool
				var canRead bool
				var canWrite bool
				singleName = nameMapper(singleDataType.Name) // TODO: 此处是否可以考虑不要用nameMapper
				canRead = true
				canWrite = true
				omitempty = false

				jsonTag := singleDataType.Tag.Get(tag)
				jsonTagList := strings.Split(jsonTag, ",")
				for singleTagIndex, singleTag := range jsonTagList {
					if singleTag == "-" {
						canRead = false
						canWrite = false
					} else if singleTag == "->" {
						canRead = false
						canWrite = true
					} else if singleTag == "<-" {
						canRead = true
						canWrite = false
					} else if singleTagIndex == 0 && singleTag != "" {
						singleName = singleTag
					} else if singleTagIndex == 1 && singleTag == "omitempty" {
						omitempty = true
					}
				}
				single := arrayMappingStructInfo{}
				single.name = singleName
				single.omitempty = omitempty
				single.canRead = canRead
				single.canWrite = canWrite
				single.index = singleDataType.Index
				single.anonymous = singleDataType.Anonymous
				if singleDataType.Anonymous {
					anonymousField = append(anonymousField, single)
				} else {
					noanonymousField = append(noanonymousField, single)
				}
			}
			result.field = append(noanonymousField, anonymousField...)
		}
	}
	return result
}

// 获取指定对象的结构信息的代理
// 通过此方法有缓存
// @param tag 用来读取json相关定义时，读取的tag
func getDataTagInfo(target reflect.Type, tag string) arrayMappingInfo {
	arrayMappingInfoMap.mutex.RLock()
	var result arrayMappingInfo
	var ok bool
	resultArray, okArray := arrayMappingInfoMap.data[tag]
	if okArray {
		result, ok = resultArray[target]
	}
	arrayMappingInfoMap.mutex.RUnlock()

	if ok {
		return result
	}
	result = getDataTagInfoInner(target, tag)

	arrayMappingInfoMap.mutex.Lock()
	if !okArray {
		resultArray = map[reflect.Type]arrayMappingInfo{}
		arrayMappingInfoMap.data[tag] = resultArray
	}
	resultArray[target] = result
	arrayMappingInfoMap.mutex.Unlock()

	return result
}

// 将目标对象转为map的实现
// 支持传入的data类型包括struct、array、map以及其指针
// @param dataValue 目标对象的反射value
// @param 用于组织返回的map的key的取值来源，从目标struct的哪个tag取值，此取值可影响data字段的name
// @return resultMap 转换后的map的value，但如果转换失败，则返回原是value
// @return isEmptyValue 是否是空值
func arrayToMapInner(dataValue reflect.Value, tag string) (reflect.Value, bool) {
	if dataValue.IsValid() == false {
		return dataValue, true
	}
	var result reflect.Value
	var isEmpty bool
	dataType := getDataTagInfo(dataValue.Type(), tag)
	if dataType.kind == TypeKindStruct && dataType.isTimeType == true {
		timeValue := dataValue.Interface().(time.Time)
		result = reflect.ValueOf(timeValue.Format("2006-01-02 15:04:05"))
		isEmpty = IsEmptyValue(dataValue)
	} else if dataType.kind == TypeKindStruct && dataType.isTimeType == false {
		resultMap := map[string]interface{}{}
		for _, singleType := range dataType.field {
			if singleType.canWrite == false {
				continue
			}
			singleResultMap, isEmptyValue := arrayToMapInner(dataValue.FieldByIndex(singleType.index), tag)
			if singleType.anonymous == false {
				if singleType.omitempty == true && isEmptyValue {
					continue
				}
				if singleResultMap.IsValid() == false {
					continue
				}
				resultMap[singleType.name] = singleResultMap.Interface()
			} else {
				combileMap(resultMap, singleResultMap)
			}
		}
		result = reflect.ValueOf(resultMap)
		isEmpty = (len(resultMap) == 0)
	} else if dataType.kind == TypeKindArray {
		resultSlice := []interface{}{}
		dataLen := dataValue.Len()
		for i := 0; i != dataLen; i++ {
			singleDataValue := dataValue.Index(i)
			singleDataResult, _ := arrayToMapInner(singleDataValue, tag)
			resultSlice = append(resultSlice, singleDataResult.Interface())
		}
		result = reflect.ValueOf(resultSlice)
		isEmpty = (len(resultSlice) == 0)
	} else if dataType.kind == TypeKindMap {
		dataKeyType := dataValue.Type().Key()
		resultMapType := reflect.MapOf(dataKeyType, interfaceType)
		resultMap := reflect.MakeMap(resultMapType)
		dataKeys := dataValue.MapKeys()
		for _, singleDataKey := range dataKeys {
			singleDataValue := dataValue.MapIndex(singleDataKey)
			singleDataResult, _ := arrayToMapInner(singleDataValue, tag)
			resultMap.SetMapIndex(singleDataKey, singleDataResult)
		}
		result = resultMap
		isEmpty = (len(dataKeys) == 0)
	} else if dataType.kind == TypeKindInterface ||
		dataType.kind == TypeKindPtr {
		result, isEmpty = arrayToMapInner(dataValue.Elem(), tag)
	} else {
		result = dataValue
		isEmpty = IsEmptyValue(dataValue)
	}
	return result, isEmpty
}

// ArrayToMap 将struct转为map
func ArrayToMap(data interface{}, tag string) interface{} {
	dataValue, _ := arrayToMapInner(reflect.ValueOf(data), tag)
	if dataValue.IsValid() == false {
		return nil
	}
	return dataValue.Interface()
}

func mapToBool(dataValue reflect.Value, target reflect.Value) error {
	dataType := dataValue.Type()
	dataKind := GetTypeKind(dataType)
	if dataKind == TypeKindBool {
		target.SetBool(dataValue.Bool())
		return nil
	} else if dataKind == TypeKindString {
		dataBool, err := strconv.ParseBool(dataValue.String())
		if err != nil {
			return errors.Errorf("不是布尔值，其值为[%s]", dataValue.String())
		}
		target.SetBool(dataBool)
		return nil
	} else {
		return errors.Errorf("不是布尔值，其类型为[%s]", dataValue.Type().String())
	}
}

func mapToUint(dataValue reflect.Value, target reflect.Value) error {
	dataType := dataValue.Type()
	dataKind := GetTypeKind(dataType)
	if dataKind == TypeKindUint {
		target.SetUint(dataValue.Uint())
		return nil
	} else if dataKind == TypeKindInt {
		target.SetUint(uint64(dataValue.Int()))
		return nil
	} else if dataKind == TypeKindFloat {
		target.SetUint(uint64(math.Floor(dataValue.Float() + 0.5)))
		return nil
	} else if dataKind == TypeKindString {
		dataUint, err := strconv.ParseUint(dataValue.String(), 10, 64)
		if err != nil {
			return errors.Errorf("不是无符号整数，其值为[%s]", dataValue.String())
		}
		target.SetUint(dataUint)
		return nil
	} else {
		return errors.Errorf("不是无符号整数，其类型为[%s]", dataValue.Type().String())
	}
}

func mapToInt(dataValue reflect.Value, target reflect.Value) error {
	dataType := dataValue.Type()
	dataKind := GetTypeKind(dataType)
	if dataKind == TypeKindInt {
		target.SetInt(dataValue.Int())
		return nil
	} else if dataKind == TypeKindUint {
		target.SetInt(int64(dataValue.Uint()))
		return nil
	} else if dataKind == TypeKindFloat {
		target.SetInt(int64(math.Floor(dataValue.Float() + 0.5)))
		return nil
	} else if dataKind == TypeKindString {
		dataInt, err := strconv.ParseInt(dataValue.String(), 10, 64)
		if err != nil {
			return errors.Errorf("不是整数，其值为[%s]", dataValue.String())
		}
		target.SetInt(dataInt)
		return nil
	} else {
		return errors.Errorf("不是整数，其类型为[%s]", dataValue.Type().String())
	}
}

func mapToFloat(dataValue reflect.Value, target reflect.Value) error {
	dataType := dataValue.Type()
	dataKind := GetTypeKind(dataType)
	if dataKind == TypeKindFloat {
		target.SetFloat(dataValue.Float())
		return nil
	} else if dataKind == TypeKindInt {
		target.SetFloat(float64(dataValue.Int()))
		return nil
	} else if dataKind == TypeKindUint {
		target.SetFloat(float64(dataValue.Uint()))
		return nil
	} else if dataKind == TypeKindString {
		dataFloat, err := strconv.ParseFloat(dataValue.String(), 64)
		if err != nil {
			return errors.Errorf("不是浮点数，其值为[%s]", dataValue.String())
		}
		target.SetFloat(dataFloat)
		return nil
	} else {
		return errors.Errorf("不是浮点数，其类型为[%s]", dataValue.Type().String())
	}
}

func mapToString(dataValue reflect.Value, target reflect.Value) error {
	stringValue := fmt.Sprintf("%v", dataValue.Interface())
	target.SetString(stringValue)
	return nil
}

func mapToArray(dataValue reflect.Value, target reflect.Value, tag string) error {
	dataType := dataValue.Type()
	dataKind := GetTypeKind(dataType)
	if dataKind != TypeKindArray {
		return errors.Errorf("不是数组，其类型为[%s]", dataValue.Type().String())
	}
	//增长空间
	dataLen := dataValue.Len()
	targetType := target.Type()
	targetLen := target.Len()
	if targetType.Kind() == reflect.Slice {
		if target.IsNil() == true {
			var newTarget reflect.Value
			newTarget = reflect.MakeSlice(targetType, dataLen, dataLen)
			target.Set(newTarget)
		} else if targetLen != dataLen {
			var newTarget reflect.Value
			newTarget = reflect.MakeSlice(targetType, dataLen, dataLen)
			reflect.Copy(newTarget, target)
			target.Set(newTarget)
		}
		targetLen = dataLen
	}
	//复制数据
	for i := 0; i != targetLen; i++ {
		if i >= dataLen {
			targetElemType := targetType.Elem()
			zeroElemType := reflect.Zero(targetElemType)
			for i := dataLen; i < targetLen; i++ {
				target.Index(i).Set(zeroElemType)
			}
			break
		} else {
			singleData := dataValue.Index(i)
			singleDataTarget := target.Index(i)
			err := mapToArrayInner(singleData, singleDataTarget, tag)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func mapToMap(dataValue reflect.Value, target reflect.Value, tag string) error {
	dataType := dataValue.Type()
	dataKind := GetTypeKind(dataType)
	if dataKind != TypeKindMap {
		return errors.Errorf("不是映射，其类型为[%s]", dataValue.Type().String())
	}
	dataKeys := dataValue.MapKeys()
	targetType := target.Type()
	targetKeyType := targetType.Key()
	targetValueType := targetType.Elem()
	if target.IsNil() == true {
		var newTarget reflect.Value
		newTarget = reflect.MakeMap(targetType)
		target.Set(newTarget)
	}
	for _, singleDataKey := range dataKeys {
		singleDataTargetKey := reflect.New(targetKeyType).Elem()
		err := mapToArrayInner(singleDataKey, singleDataTargetKey, tag)
		if err != nil {
			return err
		}

		singleDataValue := dataValue.MapIndex(singleDataKey)
		singleDataTargetValue := reflect.New(targetValueType).Elem()
		singleDataTargetValueOld := target.MapIndex(singleDataTargetKey)
		if singleDataTargetValueOld.IsValid() == true {
			singleDataTargetValue.Set(singleDataTargetValueOld)
		}
		err = mapToArrayInner(singleDataValue, singleDataTargetValue, tag)
		if err != nil {
			return errors.Errorf("参数%s%s", singleDataKey, err.Error())
		}
		target.SetMapIndex(singleDataTargetKey, singleDataTargetValue)
	}
	return nil
}

func mapToTime(dataValue reflect.Value, target reflect.Value) error {
	dataType := dataValue.Type()
	if dataType == reflect.TypeOf(time.Time{}) {
		target.Set(dataValue)
	} else if dataType.Kind() == reflect.String {
		timeValue, err := time.ParseInLocation("2006-01-02 15:04:05", dataValue.String(), time.Now().Local().Location())
		if err != nil {
			return errors.Errorf("不是时间，其值为[%s]", dataValue.String())
		}
		target.Set(reflect.ValueOf(timeValue))
		return nil
	}
	return errors.Errorf("不是时间，其类型为[%s]", dataValue.Type().String())
}

func mapToStruct(dataValue reflect.Value, target reflect.Value, targetType arrayMappingInfo, tag string) error {
	dataType := dataValue.Type()
	dataKind := GetTypeKind(dataType)
	if dataKind != TypeKindMap {
		return errors.Errorf("不是映射，其类型为[%s]", dataValue.Type().String())
	}
	dataTypeKey := dataType.Key()
	for _, singleStructInfo := range targetType.field {
		if singleStructInfo.canRead == false {
			continue
		}
		if singleStructInfo.anonymous == true {
			//FIXME 暂不考虑匿名结构体的覆盖问题
			singleDataValue := target.FieldByIndex(singleStructInfo.index)
			err := mapToArrayInner(dataValue, singleDataValue, tag)
			if err != nil {
				return errors.Errorf("参数%s%s", singleStructInfo.name, err.Error())
			}
		} else {
			singleMapKey := reflect.New(dataTypeKey)
			singleDataKey := reflect.ValueOf(singleStructInfo.name)
			err := mapToArrayInner(singleDataKey, singleMapKey, tag)
			if err != nil {
				return err
			}

			singleDataValue := target.FieldByIndex(singleStructInfo.index)
			singleMapResult := dataValue.MapIndex(singleMapKey.Elem())
			if singleMapResult.IsValid() == false {
				continue
			}
			err = mapToArrayInner(singleMapResult, singleDataValue, tag)
			if err != nil {
				return errors.Errorf("参数%s%s", singleDataKey, err.Error())
			}
		}
	}
	return nil
}

func mapToPtr(dataValue reflect.Value, target reflect.Value, tag string) error {
	targetElem := target.Elem()
	if targetElem.IsValid() == false {
		targetElem = reflect.New(target.Type().Elem())
		target.Set(targetElem)
	}
	return mapToArrayInner(dataValue, targetElem, tag)
}

func mapToInterface(dataValue reflect.Value, target reflect.Value, tag string) error {
	targetElem := target.Elem()
	if targetElem.IsValid() == false {
		target.Set(dataValue)
		return nil
	}
	newTargetElem := reflect.New(targetElem.Type()).Elem()
	newTargetElem.Set(targetElem)
	err := mapToArrayInner(dataValue, newTargetElem, tag)
	if err != nil {
		return err
	}
	target.Set(newTargetElem)
	return nil
}

func mapToArrayInner(data reflect.Value, target reflect.Value, tag string) error {
	//处理data是个nil的问题
	if data.IsValid() == false {
		target.Set(reflect.Zero(target.Type()))
		return nil
	}
	//处理data是多层嵌套的问题
	dataKind := data.Type().Kind()
	if dataKind == reflect.Interface {
		return mapToArrayInner(data.Elem(), target, tag)
	} else if dataKind == reflect.Ptr {
		return mapToArrayInner(data.Elem(), target, tag)
	}
	//根据target是多层嵌套的问题
	targetType := getDataTagInfo(target.Type(), tag)
	if targetType.kind == TypeKindPtr {
		return mapToPtr(data, target, tag)
	} else if targetType.kind == TypeKindInterface {
		return mapToInterface(data, target, tag)
	}
	//处理data是个空字符串
	if dataKind == reflect.String && data.String() == "" {
		target.Set(reflect.Zero(target.Type()))
		return nil
	}
	switch targetType.kind {

	case TypeKindBool:
		return mapToBool(data, target)
	case TypeKindInt:
		return mapToInt(data, target)
	case TypeKindUint:
		return mapToUint(data, target)
	case TypeKindFloat:
		return mapToFloat(data, target)
	case TypeKindString:
		return mapToString(data, target)
	case TypeKindArray:
		return mapToArray(data, target, tag)
	case TypeKindMap:
		return mapToMap(data, target, tag)
	case TypeKindStruct:
		if targetType.isTimeType {
			return mapToTime(data, target)
		}
		return mapToStruct(data, target, targetType, tag)
	default:
		return errors.Errorf("unkown target type %s", target.Type().String())
	}
}

// MapToArray 将map转为struct、map或array
// 此操作与ArrayToMap互为相反操作
func MapToArray(data interface{}, target interface{}, tag string) error {
	dataValue := reflect.ValueOf(data)
	targetValue := reflect.ValueOf(target)
	if targetValue.IsValid() == false {
		return errors.New("target is nil")
	}
	if targetValue.Kind() != reflect.Ptr {
		return errors.New("invalid target is not ptr")
	}
	return mapToArrayInner(dataValue, targetValue, tag)
}
