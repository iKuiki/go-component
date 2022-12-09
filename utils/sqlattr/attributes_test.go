package sqlattr_test

import (
	"database/sql/driver"
	"testing"

	"github.com/iKuiki/go-component/utils/sqlattr"
	"github.com/stretchr/testify/assert"
)

type testAttr struct {
	DataA string `attr:"field:dataA;omitempty:true"`
	// 附加字段，保存未识别的字段
	Extra map[string]string `attr:"extra:true"`
}

// Scan 扫描
func (attr *testAttr) Scan(value interface{}) error {
	return sqlattr.ScanAttr(attr, string(value.([]byte)))
}

// Value 序列化
func (attr testAttr) Value() (driver.Value, error) {
	attrStr, e := sqlattr.ValueAttr(attr)
	return []byte(attrStr), e
}

// 测试attr的序列化与解析
func TestAttributesMarshal(t *testing.T) {
	a := testAttr{
		DataA: "abc",
		Extra: map[string]string{
			"dataA": "def",
		},
	}
	val, err := a.Value()
	assert.NoError(t, err)
	var a2 testAttr
	err = a2.Scan(val)
	assert.NoError(t, err)
	assert.Equal(t, a.DataA, a2.DataA)
}
