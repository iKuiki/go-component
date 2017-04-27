package daocloudsdk

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/httplib"
	"github.com/panthesingh/goson"
)

type DaocloudSDK struct {
	accessToken string
}

func NewDaocloudSDK(accessToken string) (sdk *DaocloudSDK, err error) {
	sdk = &DaocloudSDK{
		accessToken: accessToken,
	}
	return sdk, nil
}

func (this *DaocloudSDK) GetBuildFlowList() (flows []BuildFlow, err error) {
	flows = make([]BuildFlow, 0)
	req := httplib.Get("https://openapi.daocloud.io/v1/build-flows")
	this.fillHeader(req)
	body, e := req.Bytes()
	if err != nil {
		return flows, errors.New("Request error: " + e.Error())
	}
	gson, e := goson.Parse(body)
	if e != nil {
		return flows, errors.New("Parse response error: " + e.Error())
	}
	if e := json.Unmarshal([]byte(gson.Get("build_flows").String()), &flows); e != nil {
		return flows, errors.New("Unmarshal response error: " + e.Error())
	}
	return flows, nil
}

func (this *DaocloudSDK) Build(flow BuildFlow, branch string) (build Build, err error) {
	build = Build{}
	req := httplib.Post("https://openapi.daocloud.io/v1/build-flows/" + flow.ID + "/builds")
	this.fillHeader(req)
	req.Header("Content-Type", "application/json")
	param := BuildParam{
		Branch: branch,
	}
	data, e := json.Marshal(param)
	if e != nil {
		return build, errors.New("Marshal request param error: " + e.Error())
	}
	req.Body(data)
	body, e := req.Bytes()
	if err != nil {
		return build, errors.New("Request error: " + e.Error())
	}
	if e := json.Unmarshal(body, &build); e != nil {
		return build, errors.New("Unmarshal response error: " + e.Error())
	}
	return build, nil
}
