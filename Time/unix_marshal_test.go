package Time

import (
	"encoding/json"
	"testing"
	"time"
)

type UnixTestStruct struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	CreatedAt UnixTime `json:"created_at"`
}

func TestUnixPrint(t *testing.T) {
	nowTime := UnixTime(time.Now())
	t.Log(nowTime)
}

func TestUnixMarshal(t *testing.T) {
	test := UnixTestStruct{
		ID:        1,
		Name:      "test",
		CreatedAt: UnixTime(time.Now()),
	}
	js, _ := json.Marshal(test)
	t.Log(string(js))
}

func TestUnixUnmarshal(t *testing.T) {
	var test UnixTestStruct
	js := `{"id":1,"name":"test","created_at":1512985068}`
	err := json.Unmarshal([]byte(js), &test)
	if err != nil {
		t.Fatalf("Unmarshal error: %#v", err)
	}
	t.Log(test)
}
