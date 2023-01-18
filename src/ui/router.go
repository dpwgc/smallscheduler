package ui

import (
	"alisa-dispatch-center/src/common"
	"alisa-dispatch-center/src/constant"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

// InitHttpRouter HTTP路由配置
func InitHttpRouter() {

	api := initApi()

	r := httprouter.New()

	r.GET(fmt.Sprintf("%s%s", constant.HttpUriPrefix, "/tasks"), middleware(api.ListTask))
	r.PUT(fmt.Sprintf("%s%s", constant.HttpUriPrefix, "/task"), middleware(api.SaveTask))

	//加载端口号
	port := common.Config.Server.Port
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), r)
	if err != nil {
		log.Fatal(constant.LogInfoTag, err)
		return
	}
}

// 中间件
func middleware(h ...httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// TODO
		for _, handler := range h {
			handler(w, r, p)
		}
	}
}
