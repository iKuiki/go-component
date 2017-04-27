package daocloudsdk

import (
	"github.com/astaxie/beego/httplib"
)

func (this *DaocloudSDK) fillHeader(req *httplib.BeegoHTTPRequest) {
	req.Header("Authorization", "token "+this.accessToken)
}
