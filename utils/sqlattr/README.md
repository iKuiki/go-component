# Attributes使用说明

为了方便Golang对数据库中attributes的访问，特地写了基础方法提供sql driver层的对attribute的序列化与反序列化。

基本原理：

go的sql driver层支持自定义数据类型的扫码与取值方法，只要实现了driver层这个接口的数据类型都可以被sql直接写入到对应的数据列中。

## 定义attributes

我们对一个新的model要添加attributes字段的支持的话，首先需要确定这个model的attribute字段都有哪些field，blackcard-api项目下已经提供了一个扫描全表的attributes字段的小工具scanAttr。

### scanAttr

scanAttr小工具会遍历全表，取出所有attributes并解析所有可能的field，对所有的field会取一个非空值来参考

对于枚举值（遍历后该field取值可能数少于10的），scanAttr工具会列出其所有可能的枚举值。

小工具使用方式如下：

``` bash
cd blackcard-api
# 以shop表为例
go run tools/scanAttr/scanAttr.go -t shop
# 对于某些过大的表(如user\order等几十上百万行的)，取出所有列可能不太现实，我们只需要遍历一部分数据就好了,则可以加上-total参数来限制需要遍历的总数据
go run tools/scanAttr/scanAttr.go -t user -total 10000 -c 1000 # 总共取10000条数据，-c是分页大小，每页1000条数据
```

执行完成后，会返回如下数据：

``` log
# 已精简无用信息
tableName: user, pageSize: 1000 # 传入配置，要扫描的表名、分页大小
total count:  2369290 # 表总count
total need:  10000 # 需要扫码的行数

# 以下为每行的具体数据
# 本处只显示了user表attributes的部分field
access_token: count(values)[766] value authbseB38653e40abe84bb88a59330e1836dE62 # 数据列，10000行中发现了766中数据类型，详细detail values不予显示
adminGroupId: count(values)[4] value 4 # 也是数据列，但是因为数据比较少，所以会显示所有可能的values
values:  [2 3 1 4]
adminShopId: count(values)[5] value 62
values:  [2 3 4 1 62]
channel: count(values)[3] value alipay # 枚举列，这样的列在struct中应当定义为枚举类型
values:  [mini scan_pay alipay]
gender: count(values)[3] value 男 # 枚举列，这样的列在struct中应当定义为枚举类型
values:  [男 未知 女]
payAppid: count(values)[2] value wx2d5696dbf79087f9 # 半数据半枚举列，这样的列在struct中应当根据上下文代码决定是否需要定义为枚举，payAppid就属于不应当定义为枚举的列
values:  [wx2d5696dbf79087f9 wx2332b095db4081ab]

# 以下为所有扫描发现的field
# attributes应当实现所有field的对应字段才能保证能够正常读取
# mysql的attr中有attrStruct中不存在的数据列不会导致报错，但是会导致无法在代码中读取该数据
# attrStruct中如果有数据列的定义无法从mysql的结果中解析，则一定会导致报错
[access_token address adminGroupId adminShopId allinpayUserid avatar bdMobile bdName channel city country currentCity currentCityId gender joinedZhuanpanEvent next_shop_ids old_mobile_beifen openId originWxAvatar payAppid province token wxAvatar]
```

在扫描attr完成后，就是构造数据列

### attrStruct的tag

为了能够更好的实现struct与mysql数据列的序列化与反序列化，我们需要在struct中定义数据需要如何被反序列化，以及数据的字段对应关系。

目前attrStruct的tag可读取2个信息

- field struct字段对应到mysql中attributes的field的对应关系
- type 从attributes的field中读取、写入数据时使用的序列化、反序列化方式，只对attrStruct的数据列为struct与map时有效，当前可选的type为json
- default 默认值，如果字段为空，则给定一个默认值

例如，setting的attrStruct定义如下：

``` golang
// SettingAttributes 设置的attributes
type SettingAttributes struct {
    HistoryVersion map[string]SettingAttributesHistoryVersionItem `attr:"field:historyVersion;type:json"` // 历史版本,json格式
    Force          int                                            `attr:"field:force"`                    // 是否强制更新，0:否，1:是
    Title          string                                         `attr:"field:title"`                    // 标题
    AdminName      string                                         `attr:"field:adminName"`                // 管理员的名字
    GenerateTime   time.Time                                      `attr:"field:generateTime"`             // 声称时间
    LastModified   time.Time                                      `attr:"field:lastModified"`             // 最后编辑时间
}

// SettingAttributesHistoryVersionItem 设置的attributes中历史版本的item
// historyVersion是一个php的关联数组，在go中就是map[string]SettingAttributesHistoryVersionItem
type SettingAttributesHistoryVersionItem struct {
    AdminID int64  `json:"admin_id"`
    Updated string `json:"updated"`
    Value   string `json:"value"`
    Version int64  `json:"version"`
}
```

上面SettingAttributes的HistoryVersion字段为php中的关联数组，此类型对应到go中的map，如果时一个object的json则应当对应到go中的一个struct

关于default，因为php的处理逻辑中，对于值不存在的情况，有些字段有设置默认值，所以我们也必须针对此设置默认值。
对于php的默认值，我们可以去php代码对应的Model检查__get这个魔术方法，例如User中的isNewCustomer字段定义如下

``` php
            case 'isNewCustomer':
                return $this->getDbAttribute('isNewCustomer', 1);
```

从中我们可以看出，php的这个isNewCustomer字段如果值不存在，则应给默认值1

### attrStruct的类型

根据attr扫描的结果，我们为user构建Attributes字段，应当将该字段定义为一个struct，struct的每一个字段对应mysql中attributes列的每一个field，根据struct中字段的type，反序列化时会使用不同的机制。

#### time.Time

对于time.Time,默认会按照2006-01-02 15:04:05格式来识别为时间，根据Tag中attr的type字段，可以识别不同格式的时间。目前支持的type如下

- datetime 也是默认格式，如果不传该参数则使用此格式，格式为2006-01-02 15:04:05
- date 只有日期的格式，格式为2006-01-02
- timestamp Unix时间戳，期待格式为精确到秒的Unix时间戳(一般为10位)

#### model.GeoPoint

对于model.GeoPoint，为了在es里能将其解析到geo类型，将其格式化为lat,lng这样的形式，以方便dts模块直接作为字符串同步。

- geostring geo字符串类型，格式为lat,lng，注意：纬度必须在前，经度在后

#### []model.GeoPoint

对于gps坐标数组，为了便于存储为一些标准格式，增加了以下解析与格式化

- json 直接储存为json，格式为```[{"Lat":22.544171,"Lng":113.941987},{"Lat":22.543126,"Lng":113.942003}]```这样的格式
- geoarray 按地图软件标准，储存为经度在前的二维数组，内层数组为gps点，外层为点的数组，例如```[[113.941987,22.544171],[113.942003,22.543126]]```

#### String

因为attributes在数据库中的储存就是一个map[string]string的格式，其数据格式原生就是string，所以string数据格式是无需处理，直接赋值的。

### Bool

当struct的type为bool时，会直接判断字符串为"true"或"false"，根据Tag中attr的type字段，可以识别不同的bool格式。目前支持的type如下

- truefalse 默认格式，"true"=>true,"false"=>false
- int 数字格式，"1"=>true,"0"=>false
- yesno yes或no的格式，"yes"=>true,"no"=>false
- yesnone yes或none的格式，"yes"=>true,"none"=>false

### Int64

当struct的type为int64时，会尝试将字符串以strconv.ParseInt方法解析为整型数值，目前只支持int64，不对int、int32做支持

### Float64

当struct的type为float64时，会尝试将字符串以strconv.ParseFloat方法解析为浮点数，目前只支持float64，暂不考虑兼容float32

### Struct\Map\Slice

当struct字段的type为struct或map时，则会根据字段tag中定义的type来将字符串解析到struct\map中，默认以json作为解析方法

### Slice

当struct字段的slice或map时，则会根据字段tag中定义的type来将字符串解析到slice中，默认以json作为解析方法，支持如下方法

- json 将value作为json string反序列化到切片中
- comma 将value作为逗号分隔的元素反序列化到切片中，目前只支持int64、float64、string三种形式的逗号分隔数组
- space 将value作为空格分隔的元素反序列化到切片中，目前只支持int64、float64、string三种形式的空格分隔数组

## 实现sql的scan与value方法

通过实现sql driver的scan于value方法，可以让go的sql驱动直接将数据从mysql中读取、写入。
在attributes_marshal中实现了通过反射将attributes映射到attrStruct中的方法，所以只需要在每个struct内实现Scan于Value方法并调用attributes_marshal中的通过用scanAttr、valueAttr方法即可。
示例如下：

``` golang
// Scan 扫描
func (attr *SettingAttributes) Scan(value interface{}) error {
    return scanAttr(attr, string(value.([]byte)))
}

// Value ShopDealPicture序列化
func (attr SettingAttributes) Value() (driver.Value, error) {
    attrStr, e := valueAttr(attr)
    return []byte(attrStr), e
}
```

## 测试

当实现好attr后，还应当提供相关的test方法来保证其正常工作，我们至少应该提供一个读一个写的方法来保证其能正常读写数据。

Setting模型的测试方法如下

``` golang


// 测试读取数据，以保证数据库数据都能被正常解析
func TestReadSetting(t *testing.T) {
    var count uint64
    err := mainDB.Model(model.Setting{}).Count(&count).Error
    assert.NoError(t, err)
    t.Log(count)

    pageSize := int64(5)
    for offset := int64(0); offset < int64(count); offset += pageSize {
        var settings []model.Setting
        err = mainDB.Limit(pageSize).Offset(offset).Find(&settings).Error
        assert.NoErrorf(t, err, "当前offset: %d", offset)
    }
}

// 随机读取一个Setting，然后用Create方法创建，测试attr写入是否正常
func TestWriteSetting(t *testing.T) {
    var setting model.Setting
    err := mainDB.First(&setting).Error
    assert.NoError(t, err)
    setting2 := model.Setting{
        Key:        testtool.CreateTestSetting().Key,
        Title:      setting.Title + " for rope test",
        Value:      setting.Value,
        AdminID:    setting.AdminID,
        Version:    setting.Version,
        Status:     setting.Status,
        Attributes: setting.Attributes,
    }
    err = mainDB.Create(&setting2).Error
    assert.NoError(t, err)
    if setting2.ID > 0 {
        defer mainDB.Where("id = ?", setting2.ID).Delete(&setting2)
    }
    t.Logf("setting2 saved: %d", setting2.ID)
}
```
