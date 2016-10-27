package youdaosdk

import (
	"github.com/astaxie/beego/httplib"
)

func fillHeader(req *httplib.BeegoHTTPRequest, ydSdk *YoudaoSdk) {
	req.Param("keyfrom", ydSdk.appKey)
	req.Param("key", ydSdk.appSecret)
	req.Param("type", "data")
	req.Param("doctype", "json")
	req.Param("version", "1.1")
}
