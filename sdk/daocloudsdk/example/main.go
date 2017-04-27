package main

import (
	"fmt"
	"github.com/yinhui87/go-component/config"
	"github.com/yinhui87/go-component/sdk/daocloudsdk"
)

func main() {
	config.Load(".env")
	sdk, err := daocloudsdk.NewDaocloudSDK(config.Env("DAOCLOUD_ACCESS_TOKEN"))
	if err != nil {
		panic(err)
	}
	flows, err := sdk.GetBuildFlowList()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Got %d buildflow\n", len(flows))
	if len(flows) > 0 {
		flow := flows[0]
		build, err := sdk.Build(flow, "master")
		if err != nil {
			fmt.Println("Build  " + flow.Name + " error: " + err.Error())
		} else {
			fmt.Printf("Build "+flow.Name+" created: %#v\n", build)
		}
	}
}
