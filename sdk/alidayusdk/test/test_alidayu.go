package main

import (
	"fmt"
	"github.com/yinhui87/go-component/sdk/alidayusdk"
)

func main() {
	alidayuSdk, err := alidayusdk.NewAlidayu("appKey", "appSecret",
		alidayusdk.AlidayuConfig{
			SignName: "SignName",
		})
	if err != nil {
		panic(err)
	}
	err = alidayuSdk.SendMsg("13000000000", "TpleId", alidayusdk.AlidayuSmsCommonTemplateParam{
		Code:    "Code",
		Product: "Product",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Send Msg success")
}
