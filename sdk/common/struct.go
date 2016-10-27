package common

import (
	"fmt"
)

type SdkError struct {
	Code    int    `url:"error" jsonp:"error" json:"ret"`
	Message string `url:"error_description" jsonp:"error_description" json:"msg"`
}

func (this *SdkError) GetCode() int {
	return this.Code
}

func (this *SdkError) GetMsg() string {
	return this.Message
}

func (this *SdkError) Error() string {
	return fmt.Sprintf("错误码为：%v，错误描述为：%v", this.Code, this.Message)
}
