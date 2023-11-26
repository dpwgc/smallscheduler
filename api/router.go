package api

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"smallscheduler/base"
)

// InitHttpRouter HTTP路由配置
func InitHttpRouter() {

	controller, err := NewController()
	if err != nil {
		log.Fatal(base.LogErrorTag, err)
		return
	}

	router := httprouter.New()

	router.GET(fmt.Sprintf("%s%s", base.Config().Server.ContextPath, "/tasks"), middleware(controller.ListTask))
	router.PUT(fmt.Sprintf("%s%s", base.Config().Server.ContextPath, "/task"), middleware(controller.SaveTask))

	port := base.Config().Server.Port
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
