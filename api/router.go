package api

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log/slog"
	"net/http"
	"smallscheduler/base"
	"smallscheduler/core"
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

	contextPath := base.Config().Server.ContextPath

	router.GET(fmt.Sprintf("%s%s", contextPath, "/tasks"), middleware(controller.ListTask))
	router.GET(fmt.Sprintf("%s%s", contextPath, "/task/:id"), middleware(controller.GetTask))

	router.GET(fmt.Sprintf("%s%s", contextPath, "/tags"), middleware(controller.ListTag))
	router.GET(fmt.Sprintf("%s%s", contextPath, "/crons"), middleware(controller.ListCron))

	router.POST(fmt.Sprintf("%s%s", contextPath, "/task"), middleware(controller.AddTask))
	router.PUT(fmt.Sprintf("%s%s", contextPath, "/task/:id"), middleware(controller.EditTask))
	router.DELETE(fmt.Sprintf("%s%s", contextPath, "/task/:id"), middleware(controller.DeleteTask))

	router.GET(fmt.Sprintf("%s%s", contextPath, "/execute/:id"), middleware(controller.ExecuteTask))

	router.GET(fmt.Sprintf("%s%s", contextPath, "/records"), middleware(controller.ListRecord))

	router.GET(fmt.Sprintf("%s%s", contextPath, "/health"), controller.Health)
	router.GET(fmt.Sprintf("%s%s", contextPath, "/shutdown"), controller.Shutdown)

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
		if core.Shutdown {
			w.WriteHeader(ErrorCode)
			_, err := w.Write([]byte(""))
			if err != nil {
				base.Logger.Error(err.Error())
			}
			return
		}
		base.Logger.Info(fmt.Sprintf("[%s]%s", r.Method, r.RequestURI), slog.String("remoteAddr", r.RemoteAddr), slog.Int64("contentLength", r.ContentLength))
		for _, handler := range h {
			handler(w, r, p)
		}
	}
}
