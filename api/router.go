package api

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log/slog"
	"net/http"
	"smallscheduler/base"
)

// InitHttpRouter HTTP路由配置
func InitHttpRouter() {

	controller, err := NewController()
	if err != nil {
		base.Logger.Error(err.Error())
		return
	}

	base.Logger.Info("start http router")

	router := httprouter.New()

	router.GET(fmt.Sprintf("%s%s", base.Config().Server.ContextPath, "/tasks"), middleware(controller.ListTask))
	router.GET(fmt.Sprintf("%s%s", base.Config().Server.ContextPath, "/task/:id"), middleware(controller.GetTask))

	router.POST(fmt.Sprintf("%s%s", base.Config().Server.ContextPath, "/task"), middleware(controller.AddTask))
	router.PUT(fmt.Sprintf("%s%s", base.Config().Server.ContextPath, "/task/:id"), middleware(controller.EditTask))
	router.DELETE(fmt.Sprintf("%s%s", base.Config().Server.ContextPath, "/task/:id"), middleware(controller.DeleteTask))

	router.GET(fmt.Sprintf("%s%s", base.Config().Server.ContextPath, "/records"), middleware(controller.ListRecord))

	router.ServeFiles(fmt.Sprintf("%s%s", base.Config().Server.ContextPath, "/web/*filepath"), http.Dir("web"))

	port := base.Config().Server.Port
	err = http.ListenAndServe(fmt.Sprintf(":%v", port), router)
	if err != nil {
		base.Logger.Error(err.Error())
		panic(err)
		return
	}
}

// 中间件
func middleware(h ...httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		base.Logger.Info("api", slog.String("uri", r.RequestURI), slog.String("method", r.Method), slog.String("remoteAddr", r.RemoteAddr))
		for _, handler := range h {
			handler(w, r, p)
		}
	}
}
