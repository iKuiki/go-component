package kugousdk

import (
	"testing"
)

func TestSearchMusic(t *testing.T) {
	page := uint64(1)
	size := uint64(10)
	keyword := "keyword"
	sdk, err := NewKugouSdk()
	if err != nil {
		t.Fatalf("KugouSdk Init Fail: %s\n", err.Error())
	}
	songs, err := sdk.Search(keyword, page, size)
	if err != nil {
		t.Fatalf("kugousdk Search Error: %s\n", err.Error())
	}
	for k, v := range songs {
		t.Logf("songs%d: %s - %s\n", k, v.Songname, v.Singername)
		d, err := sdk.Detail(v.Hash)
		if err != nil {
			t.Errorf("Get detail for %s error: %s\n", v.Hash, err.Error())
		} else {
			t.Logf("filesize: %d, url: %s\n", d.FileSize, d.URL)
		}
	}
}
