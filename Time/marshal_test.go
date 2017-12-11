package Time

import (
	"encoding/json"
	"testing"
	"time"
)

type TestStruct struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt Time   `json:"created_at"`
}

func TestPrint(t *testing.T) {
	nowTime := Time(time.Now())
	t.Log(nowTime)
}

func TestMarshal(t *testing.T) {
	test := TestStruct{
		ID:        1,
		Name:      "test",
		CreatedAt: Time(time.Now()),
	}
	js, _ := json.Marshal(test)
	t.Log(string(js))
}

func TestUnmarshal(t *testing.T) {
	var test TestStruct
	js := `{"id":1,"name":"test","created_at":"2017-12-11 17:57:56"}`
	err := json.Unmarshal([]byte(js), &test)
	if err != nil {
		t.Fatalf("Unmarshal error: %#v", err)
	}
	t.Log(test)
}
