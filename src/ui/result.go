package ui

import (
	"alisa-dispatch-center/src/common"
	"alisa-dispatch-center/src/constant"
	"encoding/json"
	"net/http"
)

type resultSuccessDTO struct {
	Code int16 `json:"code"`
	Data any   `json:"data"`
}

type resultFailDTO struct {
	Code int16  `json:"code"`
	Msg  string `json:"msg"`
}

// 统一返回模版
func result(w http.ResponseWriter, code int16, msg string, data any) {
	//响应成功
	if code == constant.HttpRequestSuccessCode {
		result := resultSuccessDTO{
			Code: code,
			Data: data,
		}
		resultBytes, err := json.Marshal(result)
		if err != nil {
			common.Log.Println(constant.LogErrorTag, err)
		}
		_, err = w.Write(resultBytes)
		if err != nil {
			return
		}
	} else {
		result := resultFailDTO{
			Code: code,
			Msg:  msg,
		}
		resultBytes, err := json.Marshal(result)
		if err != nil {
			common.Log.Println(constant.LogErrorTag, err)
		}
		_, err = w.Write(resultBytes)
		if err != nil {
			return
		}
	}
}
