package utils

import (
	"reflect"
	"sync"
)

const (
	// TypeKindBool 布尔
	TypeKindBool int = 1
	// TypeKindInt 有符号整数
	TypeKindInt int = 2
	// TypeKindUint 无符号整数
	TypeKindUint int = 3
	// TypeKindFloat 浮点数
	TypeKindFloat int = 4
	// TypeKindPtr 指针
	TypeKindPtr int = 5
	// TypeKindString 字符串
	TypeKindString int = 6
	// TypeKindArray 数组
	TypeKindArray int = 7
	// TypeKindMap 映射
	TypeKindMap int = 8
	// TypeKindStruct 结构体
	TypeKindStruct int = 9
	// TypeKindInterface 接口
	TypeKindInterface int = 10
	// TypeKindFunc 函数
	TypeKindFunc int = 11
	// TypeKindChan 通道
	TypeKindChan int = 12
	// TypeKindOther 其他
	TypeKindOther int = 13
)

// // TypeKind kind类型
// var TypeKind struct {
// 	EnumStruct
// 	BOOL      int `enum:"1,布尔"`
// 	INT       int `enum:"2,有符号整数"`
// 	UINT      int `enum:"3,无符号整数"`
// 	FLOAT     int `enum:"4,浮点数"`
// 	PTR       int `enum:"5,指针"`
// 	STRING    int `enum:"6,字符串"`
// 	ARRAY     int `enum:"7,数组"`
// 	MAP       int `enum:"8,映射"`
// 	STRUCT    int `enum:"9,结构体"`
// 	INTERFACE int `enum:"10,接口"`
// 	FUNC      int `enum:"11,函数"`
// 	CHAN      int `enum:"12,通道"`
// 	OTHER     int `enum:"13,其他"`
// }

// func init() {
// 	InitEnumStruct(&TypeKind)
// }

// GetTypeKind 获取reflect下的Type
// 此处获取的kind是一个近似kind
// 如int\int8\int16\int32都归为int
func GetTypeKind(t reflect.Type) int {
	switch t.Kind() {
	case reflect.Bool:
		return TypeKindBool
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return TypeKindInt
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return TypeKindUint
	case reflect.Float32, reflect.Float64:
		return TypeKindFloat
	case reflect.Ptr:
		return TypeKindPtr
	case reflect.String:
		return TypeKindString
	case reflect.Array, reflect.Slice:
		return TypeKindArray
	case reflect.Map:
		return TypeKindMap
	case reflect.Struct:
		return TypeKindStruct
	case reflect.Interface:
		return TypeKindInterface
	case reflect.Func:
		return TypeKindFunc
	case reflect.Chan:
		return TypeKindChan
	default:
		return TypeKindOther
	}
}

type zeroable interface {
	IsZero() bool
}

// IsEmptyValue 判断对象是否为空
func IsEmptyValue(v reflect.Value) bool {
	k := v.Interface()
	switch k.(type) {
	case int:
		return k.(int) == 0
	case int8:
		return k.(int8) == 0
	case int16:
		return k.(int16) == 0
	case int32:
		return k.(int32) == 0
	case int64:
		return k.(int64) == 0
	case uint:
		return k.(uint) == 0
	case uint8:
		return k.(uint8) == 0
	case uint16:
		return k.(uint16) == 0
	case uint32:
		return k.(uint32) == 0
	case uint64:
		return k.(uint64) == 0
	case float32:
		return k.(float32) == 0
	case float64:
		return k.(float64) == 0
	case bool:
		return k.(bool) == false
	case string:
		return k.(string) == ""
	case zeroable:
		return k.(zeroable).IsZero()
	}
	return false
}

type getFieldByNameResult struct {
	structField reflect.StructField
	isExist     bool
}

// 获取指定的struct列的反射type
// 传入的name可以接受.分割
func getFieldByNameInner(t reflect.Type, name string) (reflect.StructField, bool) {
	nameArray := Explode(name, ".")
	if len(nameArray) == 0 {
		return reflect.StructField{}, false
	}
	var isExist bool
	var resultStruct reflect.StructField
	resultIndex := []int{}
	for _, singleName := range nameArray {
		resultStruct, isExist = t.FieldByName(singleName)
		if !isExist {
			return reflect.StructField{}, false
		}
		resultIndex = append(resultIndex, resultStruct.Index...)
		t = resultStruct.Type
	}
	resultStruct.Index = resultIndex
	return resultStruct, true
}

var ( // 反射机制获取field比较慢，增加一个缓存节约性能
	getFieldByNameCache = map[reflect.Type]map[string]getFieldByNameResult{}
	getFieldByNameMutex = sync.RWMutex{}
)

// 获取指定的struct列的反射type
// 传入的name可以接受.分割
// 此方法是getFieldByNameInner的包装，带有缓存机制
func getFieldByName(t reflect.Type, name string) (reflect.StructField, bool) {
	getFieldByNameMutex.RLock() //先查询缓存
	result, isExist := getFieldByNameCache[t][name]
	getFieldByNameMutex.RUnlock()

	if isExist {
		return result.structField, result.isExist
	}
	result.structField, result.isExist = getFieldByNameInner(t, name)

	getFieldByNameMutex.Lock()
	typeInfo, isExist := getFieldByNameCache[t]
	if !isExist {
		typeInfo = map[string]getFieldByNameResult{}
	}
	typeInfo[name] = result
	getFieldByNameCache[t] = typeInfo
	getFieldByNameMutex.Unlock()

	return result.structField, result.isExist
}
