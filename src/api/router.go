package api

import (
	"alisa-dispatch-center/src/base"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

const HttpUriPrefix = "/v1"

// InitHttpRouter HTTP路由配置
func InitHttpRouter() {

	controller, err := NewController()
	if err != nil {
		log.Fatal(base.LogErrorTag, err)
		return
	}

	router := httprouter.New()

	router.GET(fmt.Sprintf("%s%s", HttpUriPrefix, "/tasks"), middleware(controller.ListTask))
	router.PUT(fmt.Sprintf("%s%s", HttpUriPrefix, "/task"), middleware(controller.SaveTask))

	port := base.Config.Server.Port
	err = http.ListenAndServe(fmt.Sprintf(":%v", port), router)
	if err != nil {
		log.Fatal(base.LogErrorTag, err)
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
