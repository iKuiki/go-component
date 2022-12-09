package sqlattr

import (
	"encoding/json"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ParseAttr 将字符串格式的attr解析为Map
func ParseAttr(attr string) (mattr map[string]string) {
	mattr = make(map[string]string)
	if attr == "" {
		return
	}
	tattr := strings.TrimSpace(attr)
	sattr := strings.Split(tattr, "\n")
	for _, skv := range sattr {
		idx := strings.Index(skv, ":")
		if idx == -1 {
			continue
		}
		mattr[skv[0:idx]] = strings.TrimSpace(skv[idx+1:]) // 此处去除可能出现的多余换行符
	}
	return mattr
}

// FormatMattr 将map格式的attr序列化为string
func FormatMattr(mattr map[string]string) string {
	var r string
	list := make([]string, 0, len(mattr))
	for k, v := range mattr {
		s := (k + ":" + v + "\n")
		list = append(list, s)
	}
	sort.Strings(list)
	for _, kv := range list {
		r += kv
	}
	return r
}

// ScanAttr 将string的attributes value解析到attributes Struct上
// @Param attrPtr 目标attributes struct的指针
// @Param sttrStr 待解析的attributes字符串
func ScanAttr(attrPtr interface{}, attrStr string) error {
	mattr := ParseAttr(attrStr)
	return parseMattrToAttr(mattr, attrPtr)
}

// attrTagConfig attr列的tag配置
type attrTagConfig struct {
	Field      string // 对应的列名
	Type       string // 类型
	Default    string // 默认值
	OmitEmpty  bool   // 如果为0值则忽略
	IsExtraMap bool   // 该字段是否为不支持的字段储存的Map
}

// 解析AttrStruct的Tag，返回TagConfig
// 本方法对错误最高容忍度(是不是不该容忍,容忍会导致错误的attrTag不易发现)
// 解析失败则返回空
func parseAttrTagConfig(attrTagStr string) (tagConfig attrTagConfig, err error) {
	tagStrs := strings.Split(attrTagStr, ";")
	for _, tagStr := range tagStrs {
		tagSs := strings.Split(tagStr, ":")
		if len(tagSs) != 2 {
			err = errors.Errorf("invalid tagStr %s", tagStr)
			return
		}
		switch tagSs[0] {
		case "field": // 配置attr的列名
			tagConfig.Field = tagSs[1]
		case "type":
			tagConfig.Type = tagSs[1]
		case "default":
			tagConfig.Default = tagSs[1]
		case "omitempty":
			if tagSs[1] == "true" {
				tagConfig.OmitEmpty = true
			}
		case "extra":
			if tagSs[1] == "true" {
				if tagConfig.IsExtraMap {
					err = errors.Errorf("duplication extra field %v", tagSs[0])
					return
				}
				tagConfig.IsExtraMap = true
			}
		default:
			err = errors.Errorf("invalid attr tag %v", tagSs[0])
			return
		}
	}
	return
}

// 将map[string]string形式的attr配置到attrStruct中
func parseMattrToAttr(mattr map[string]string, attrPtr interface{}) error {
	attrPtrValue := reflect.ValueOf(attrPtr)
	// 传入的attrPtr应当必须是一个指针，否则就无法对其赋值了
	if attrPtrValue.Kind() != reflect.Ptr {
		return errors.Errorf("attr target is not a point, it is %v", attrPtrValue.Kind())
	}
	attrValue := attrPtrValue.Elem()
	if !attrValue.CanSet() {
		return errors.New("attr target can not set")
	}
	if attrValue.Kind() != reflect.Struct {
		return errors.Errorf("attr target is not a struct, it is %v", attrValue.Kind())
	}
	numField := attrValue.NumField()
	var extraFieldName string             // 不支持的字段的附加字段的名称（为空则不启用
	parsedFields := make(map[string]bool) // 已经识别的字段集合
	for i := 0; i < numField; i++ {
		// 获取字段的配置
		attrFieldStructType := attrValue.Type().Field(i)
		attrTag := attrFieldStructType.Tag.Get("attr")
		tagConfig, e := parseAttrTagConfig(attrTag)
		if e != nil {
			return errors.Wrapf(e, "parse Attr tag Config fail: %v", attrTag)
		}
		if tagConfig.IsExtraMap {
			extraFieldName = attrFieldStructType.Name
			continue
		}
		if tagConfig.Field == "" { // 如果未配置，则使用默认值
			tagConfig.Field = attrFieldStructType.Name
		}
		// 根据tag中给出的字段配置设置字段的值
		attrFieldValue := attrValue.Field(i)
		if !attrFieldValue.CanSet() {
			return errors.Errorf("attr target field[%d] can not set", i)
		}
		value, ok := mattr[tagConfig.Field]
		parsedFields[tagConfig.Field] = ok
		if value == "" { // value为空时尝试赋值默认值
			if tagConfig.Default == "" { // 默认值仍然为空则跳过本字段解析
				continue
			}
			value = tagConfig.Default
		}
		e = parseAttrField(tagConfig, value, attrFieldValue)
		if e != nil {
			return e
		}
	}
	// 如果附加字段不为空，则检查是否有未识别的字段，丢进附加字段里
	if extraFieldName != "" {
		attrField := attrValue.FieldByName(extraFieldName)
		if attrField.Kind() != reflect.Map {
			return errors.Errorf("attr extra field[%s] type %v is not map", extraFieldName, attrField.Kind())
		}
		if !attrField.CanSet() {
			return errors.Errorf("attr extra field[%s] can not set", extraFieldName)
		}
		// 初始化map
		attrField.Set(reflect.MakeMap(attrField.Type()))
		// 赋值
		for k, v := range mattr {
			if _, ok := parsedFields[k]; !ok {
				attrField.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
			}
		}
	}
	return nil
}

// 根据attrTagConfig将配置的字符串value反序列化到struct中对应field
func parseAttrField(attrConfig attrTagConfig, value string, attrFieldValue reflect.Value) error {
	switch attrFieldValue.Type().String() { // 针对某些特殊类型预先处理
	case "time.Time": // 时间类型
		var (
			t time.Time
			e error
		)
		switch attrConfig.Type {
		case "", "datetime":
			t, e = time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		case "date":
			t, e = time.ParseInLocation("2006-01-02", value, time.Local)
		case "timestamp":
			var stamp int64
			stamp, e = strconv.ParseInt(value, 10, 64)
			if e != nil {
				return errors.Wrapf(e, "invalid timestamp %v", value)
			}
			t = time.Unix(stamp, 0)
		default:
			e = errors.Errorf("Unsupported time.Time attr type: %s", attrConfig.Type)
		}
		if e != nil {
			return errors.Wrap(e, attrConfig.Field)
		}
		attrFieldValue.Set(reflect.ValueOf(t))
		return nil
	case "model.GeoPoint":
		var (
			g GeoPoint
			e error
		)
		switch attrConfig.Type {
		case "geostring":
			// 解析lat,lng格式的数据
			vs := strings.Split(value, ",")
			if len(vs) != 2 {
				return errors.Errorf("model.GeoPoint invalid value: %s", value)
			}
			g.Lat, e = strconv.ParseFloat(vs[0], 64)
			if e == nil {
				g.Lng, e = strconv.ParseFloat(vs[1], 64)
			}
			if e != nil {
				return errors.Wrapf(e, "model.GeoPoint invalid value: %s", value)
			}
		default:
			e = errors.Errorf("Unsupported model.GeoPoint attr type: %s", attrConfig.Type)
		}
		if e != nil {
			return errors.Wrap(e, attrConfig.Field)
		}
		attrFieldValue.Set(reflect.ValueOf(g))
		return nil
	case "[]model.GeoPoint": // 针对点坐标组
		switch attrConfig.Type {
		case "json": // 传统json格式存储
			e := json.Unmarshal([]byte(value), attrFieldValue.Addr().Interface())
			if e != nil {
				return errors.Wrap(e, attrConfig.Field)
			}
		case "geoarray": // 解析如同[[113,22],[114,22]]这样的格式，经度在前
			var (
				points [][]float64
				gs     []GeoPoint
			)
			e := json.Unmarshal([]byte(value), &points)
			if e != nil {
				return errors.Wrapf(e, "[]model.GeoPoint unmarshal to [][]float64 fail: %s", value)
			}
			for i, p := range points {
				if len(p) != 2 {
					return errors.Errorf("[]model.GeoPoint points[%d] len wrong: %v", i, p)
				}
				gs = append(gs, GeoPoint{
					Lng: p[0],
					Lat: p[1],
				})
			}
			attrFieldValue.Set(reflect.ValueOf(gs))
		}
		return nil
	}
	// TODO: 指针怎么处理？
	switch attrFieldValue.Kind() { // 根据attrField类型选择不同的反序列化策略
	case reflect.String: // 字符串直接赋值即可
		switch attrConfig.Type {
		case "":
			attrFieldValue.SetString(value)
		case "json":
			var j string
			e := json.Unmarshal([]byte(value), &j)
			if e != nil {
				return errors.Errorf("%s parse value %v fail: %v", attrConfig.Field, value, e)
			}
			attrFieldValue.SetString(j)
		default:
			return errors.Errorf("%s Unsupported attr type: %s", attrConfig.Field, attrConfig.Type)
		}
	case reflect.Bool:
		switch attrConfig.Type {
		case "", "truefalse":
			switch value { // 只兼容true、false
			case "true":
				attrFieldValue.SetBool(true)
			case "false":
				attrFieldValue.SetBool(false)
			default:
				return errors.Errorf("%s unknwon [truefalse]bool value: %s", attrConfig.Field, value)
			}
		case "int":
			switch value { // int方式,1为true，0为false
			case "1":
				attrFieldValue.SetBool(true)
			case "0":
				attrFieldValue.SetBool(false)
			default:
				return errors.Errorf("%s unknwon [int]bool value: %s", attrConfig.Field, value)
			}
		case "yesno":
			switch value { // yesno方式,yes为true，no为false
			case "yes":
				attrFieldValue.SetBool(true)
			case "no":
				attrFieldValue.SetBool(false)
			default:
				return errors.Errorf("%s unknwon [yesno]bool value: %s", attrConfig.Field, value)
			}
		case "yesnone":
			switch value { // yesnone方式,yes为true，none为false
			case "yes":
				attrFieldValue.SetBool(true)
			case "none":
				attrFieldValue.SetBool(false)
			default:
				return errors.Errorf("%s unknwon [yesnone]bool value: %s", attrConfig.Field, value)
			}
		default:
			return errors.Errorf("%s Unsupported attr type: %s", attrConfig.Field, attrConfig.Type)
		}
	case reflect.Int64:
		i, e := strconv.ParseInt(value, 10, 64)
		if e != nil {
			return errors.Wrap(e, attrConfig.Field)
		}
		attrFieldValue.SetInt(i)
	case reflect.Float64:
		f, e := strconv.ParseFloat(value, 64)
		if e != nil {
			return errors.Wrap(e, attrConfig.Field)
		}
		attrFieldValue.SetFloat(f)
	case reflect.Struct, reflect.Map:
		// 根据Type需要做不同的反序列化策略
		switch attrConfig.Type {
		case "", "json": // 如果为空也默认为json
			// 兼容php的关联数组为空数组的情况
			if value == "[]" {
				return nil
			}
			e := json.Unmarshal([]byte(value), attrFieldValue.Addr().Interface())
			if e != nil {
				return errors.Wrap(e, attrConfig.Field)
			}
		default:
			return errors.Errorf("%s Unsupported attr type: %s", attrConfig.Field, attrConfig.Type)
		}
	case reflect.Slice:
		// 根据Type需要做不同的反序列化策略
		switch attrConfig.Type {
		case "", "json": // 如果为空也默认为json
			e := json.Unmarshal([]byte(value), attrFieldValue.Addr().Interface())
			if e != nil {
				return errors.Wrap(e, attrConfig.Field)
			}
		case "comma", "space": // 分隔符分隔
			var splitSign string // 判定分隔符
			switch attrConfig.Type {
			case "comma":
				splitSign = ","
			case "space":
				splitSign = " "
			}
			items := strings.Split(value, splitSign)
			// 检查元素的type
			switch attrFieldValue.Type().String() {
			case "[]string":
				for _, item := range items {
					attrFieldValue.Set(reflect.Append(attrFieldValue, reflect.ValueOf(item)))
				}
			case "[]int64":
				for _, item := range items {
					i, e := strconv.ParseInt(item, 10, 64)
					if e != nil {
						return errors.Errorf("%s Parse int64 separator data[%s] error: %v", attrConfig.Field, item, e)
					}
					attrFieldValue.Set(reflect.Append(attrFieldValue, reflect.ValueOf(i)))
				}
			case "[]uint64":
				for _, item := range items {
					i, e := strconv.ParseUint(item, 10, 64)
					if e != nil {
						return errors.Errorf("%s Parse uint64 separator data[%s] error: %v", attrConfig.Field, item, e)
					}
					attrFieldValue.Set(reflect.Append(attrFieldValue, reflect.ValueOf(i)))
				}
			case "[]float64":
				for _, item := range items {
					i, e := strconv.ParseFloat(item, 64)
					if e != nil {
						return errors.Errorf("%s Parse float64 separator data[%s] error: %v", attrConfig.Field, item, e)
					}
					attrFieldValue.Set(reflect.Append(attrFieldValue, reflect.ValueOf(i)))
				}
			default:
				return errors.Errorf("%s Unsupported separator data type: %v", attrConfig.Field, attrFieldValue.Type())
			}
		default:
			return errors.Errorf("%s Unsupported attr type: %s", attrConfig.Field, attrConfig.Type)
		}
	default:
		return errors.Errorf("%s Unsupported field kind: %v", attrConfig.Field, attrFieldValue.Kind())
	}
	return nil
}

// ValueAttr 将给定的attrStruct序列化为attributes格式的字符串
func ValueAttr(attr interface{}) (attrStr string, err error) {
	mattr, err := formatAttrToMattr(attr)
	attrStr = FormatMattr(mattr)
	return
}

// 将attrStruct格式化到map[string]string中
func formatAttrToMattr(attr interface{}) (mattr map[string]string, err error) {
	attrValue := reflect.ValueOf(attr)
	if attrValue.Kind() == reflect.Ptr {
		attrValue = attrValue.Elem()
	}
	if attrValue.Kind() != reflect.Struct {
		err = errors.Errorf("attr source is not a struct, it is %v", attrValue.Kind())
		return
	}
	mattr = make(map[string]string)
	numField := attrValue.NumField()
	var extraFieldName string // 不支持的字段的附加字段的名称（为空则不启用
	for i := 0; i < numField; i++ {
		// 获取字段的配置
		attrFieldStructType := attrValue.Type().Field(i)
		attrTag := attrFieldStructType.Tag.Get("attr")
		tagConfig, e := parseAttrTagConfig(attrTag)
		if e != nil {
			err = errors.Wrapf(e, "parse field %v Attr tag Config fail: %v", attrFieldStructType.Name, attrTag)
			return
		}
		if tagConfig.IsExtraMap {
			extraFieldName = attrFieldStructType.Name
			continue
		}
		if tagConfig.Field == "" { // 如果未配置，则使用默认值
			tagConfig.Field = attrFieldStructType.Name
		}
		// 根据tag中给出的字段配置获取字段的值
		attrFieldValue := attrValue.Field(i)
		// 判断是否有设置为空则忽略
		if tagConfig.OmitEmpty {
			if isAttrFieldZero(attrFieldValue) {
				// 如果为0值则不将其添加到map中
				continue
			}
		}
		value, e := formatAttrField(tagConfig, attrFieldValue)
		if e != nil {
			err = e
			return
		}
		mattr[tagConfig.Field] = value
	}
	// 如果附加字段不为空，则检查是否有未识别的字段，丢进附加字段里
	if extraFieldName != "" {
		attrField := attrValue.FieldByName(extraFieldName)
		if attrField.Kind() != reflect.Map {
			err = errors.Errorf("attr extra field[%s] type %v is not map", extraFieldName, attrField.Kind())
			return
		}
		// 作为Map遍历附加字段
		iter := attrField.MapRange()
		for iter.Next() {
			k := iter.Key().String()
			// 主字段不存在这个值，才替换
			if _, ok := mattr[k]; !ok {
				mattr[k] = iter.Value().String()
			}
		}
	}
	return
}

// 判断attr字段下是否为逻辑上的空值
func isAttrFieldZero(attrFieldValue reflect.Value) bool {
	if attrFieldValue.IsZero() { // 该反射字段如果本身就是空值，则直接返回
		return true
	}
	switch attrFieldValue.Type().String() {
	case "time.Time": // 时间类型
		if attrFieldValue.MethodByName("IsZero").Call(nil)[0].Bool() {
			return true
		}
		// 如果本身非Zero，则使用Unix方法也进行判断(如果unix为0也认为等于0)
		return attrFieldValue.MethodByName("Unix").Call(nil)[0].Int() == 0
	case "model.GeoPoint": // 地理坐标类型
		return attrFieldValue.FieldByName("Lat").Float() == 0 && attrFieldValue.FieldByName("Lng").Float() == 0
	}
	switch attrFieldValue.Kind() { // 根据attrField类型选择不同的判空策略
	case reflect.Struct: // 结构体需要尝试判断其下的每个值
		fieldNum := attrFieldValue.NumField()
		for i := 0; i < fieldNum; i++ {
			if !isAttrFieldZero(attrFieldValue.Field(i)) { // 递归向下判断其每个field
				// 只要有一行不为空，则认为非空
				return false
			}
		}
		// 所有行都为空，则真的为空
		return true
	case reflect.Map, reflect.Slice: // 切片与Map都可以通过Len判断长度
		if attrFieldValue.Len() == 0 {
			return true
		}
	}
	return false
}

// 根据attrTagConfig将struct中对应field序列化到配置的字符串value
func formatAttrField(attrConfig attrTagConfig, attrFieldValue reflect.Value) (value string, err error) {
	switch attrFieldValue.Type().String() { // 针对某些特殊类型预先处理
	case "time.Time": // 时间类型
		if attrFieldValue.MethodByName("IsZero").Call(nil)[0].Bool() {
			// 如果时间为0值则返回空字符串
			return
		}
		switch attrConfig.Type {
		case "", "datetime":
			res := attrFieldValue.MethodByName("Format").Call([]reflect.Value{
				// Param of time.Format
				reflect.ValueOf("2006-01-02 15:04:05"),
			})
			value = res[0].String()
		case "date":
			res := attrFieldValue.MethodByName("Format").Call([]reflect.Value{
				// Param of time.Format
				reflect.ValueOf("2006-01-02"),
			})
			value = res[0].String()
		case "timestamp":
			res := attrFieldValue.MethodByName("Unix").Call(nil)
			value = strconv.FormatInt(res[0].Int(), 10)
		default:
			err = errors.Errorf("Unsupported time.Time attr type: %s", attrConfig.Type)
		}
		return
	case "model.GeoPoint": // 地理坐标类型
		if attrFieldValue.FieldByName("Lat").Float() == 0 && attrFieldValue.FieldByName("Lng").Float() == 0 {
			// 经纬度同时为0，则认为坐标点未填入
			return
		}
		switch attrConfig.Type {
		case "geostring":
			lat, lng := attrFieldValue.FieldByName("Lat").Float(), attrFieldValue.FieldByName("Lng").Float()
			value = strconv.FormatFloat(lat, 'f', -1, 64) + "," + strconv.FormatFloat(lng, 'f', -1, 64)
		default:
			err = errors.Errorf("Unsupported model.GeoPoint attr type: %s", attrConfig.Type)
		}
		return
	case "[]model.GeoPoint": // 针对点坐标组
		switch attrConfig.Type {
		case "json": // 传统json格式存储
			j, e := json.Marshal(attrFieldValue.Interface())
			if e != nil {
				err = errors.WithStack(e)
				return
			}
			value = string(j)
		case "geoarray": // 格式化成如同[[113,22],[114,22]]这样的格式，经度在前
			// 先取出结构，再生成
			points := [][]float64{}
			for i := 0; i < attrFieldValue.Len(); i++ {
				points = append(points, []float64{
					attrFieldValue.Index(i).FieldByName("Lng").Float(),
					attrFieldValue.Index(i).FieldByName("Lat").Float(),
				})
			}
			j, e := json.Marshal(points)
			if e != nil {
				err = errors.WithStack(e)
				return
			}
			value = string(j)
		}
		return
	}
	// TODO: 指针怎么处理？
	switch attrFieldValue.Kind() { // 根据attrField类型选择不同的序列化策略
	case reflect.String: // 字符串直接赋值即可
		switch attrConfig.Type {
		case "":
			value = attrFieldValue.String()
		case "json":
			j, e := json.Marshal(attrFieldValue.String())
			if e != nil {
				err = errors.WithStack(e)
				return
			}
			value = string(j)
		default:
			err = errors.Errorf("Unsupported attr type: %s", attrConfig.Type)
			return
		}
	case reflect.Bool:
		switch attrConfig.Type {
		case "", "truefalse":
			if attrFieldValue.Bool() {
				value = "true"
			} else {
				value = "false"
			}
		case "int":
			if attrFieldValue.Bool() {
				value = "1"
			} else {
				value = "0"
			}
		case "yesno":
			if attrFieldValue.Bool() {
				value = "yes"
			} else {
				value = "no"
			}
		case "yesnone":
			if attrFieldValue.Bool() {
				value = "yes"
			} else {
				value = "none"
			}
		}
	case reflect.Int64:
		value = strconv.FormatInt(attrFieldValue.Int(), 10)
	case reflect.Float64:
		value = strconv.FormatFloat(attrFieldValue.Float(), 'f', -1, 64)
	case reflect.Struct, reflect.Map:
		// 根据Type需要做不同的反序列化策略
		switch attrConfig.Type {
		case "", "json": // 如果为空也默认为json
			j, e := json.Marshal(attrFieldValue.Interface())
			if e != nil {
				err = errors.WithStack(e)
				return
			}
			value = string(j)
		default:
			err = errors.Errorf("Unsupported attr type: %s", attrConfig.Type)
			return
		}
	case reflect.Slice:
		// 根据Type需要做不同的反序列化策略
		switch attrConfig.Type {
		case "", "json": // 如果为空也默认为json
			j, e := json.Marshal(attrFieldValue.Interface())
			if e != nil {
				err = errors.WithStack(e)
				return
			}
			value = string(j)
		case "comma", "space": // 分隔符分隔
			var items []string
			len := attrFieldValue.Len()
			// 检查元素的type
			switch attrFieldValue.Type().String() {
			case "[]string":
				for i := 0; i < len; i++ {
					item := attrFieldValue.Index(i).String()
					items = append(items, item)
				}
			case "[]int64":
				for i := 0; i < len; i++ {
					item := strconv.FormatInt(attrFieldValue.Index(i).Int(), 10)
					items = append(items, item)
				}
			case "[]uint64":
				for i := 0; i < len; i++ {
					item := strconv.FormatUint(attrFieldValue.Index(i).Uint(), 10)
					items = append(items, item)
				}
			case "[]float64":
				for i := 0; i < len; i++ {
					item := strconv.FormatFloat(attrFieldValue.Index(i).Float(), 'f', -1, 64)
					items = append(items, item)
				}
			default:
				err = errors.Errorf("Unsupported separator data type: %v", attrFieldValue.Type())
				return
			}
			var splitSign string // 判定分隔符
			switch attrConfig.Type {
			case "comma":
				splitSign = ","
			case "space":
				splitSign = " "
			}
			value = strings.Join(items, splitSign)
		default:
			err = errors.Errorf("Unsupported attr type: %s", attrConfig.Type)
			return
		}
	default:
		err = errors.Errorf("Unsupported field kind: %v", attrFieldValue.Kind())
		return
	}
	return
}
