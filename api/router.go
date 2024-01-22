package api

import (
	"encoding/json"
	"fmt"
	"github.com/dpwgc/easierweb"
	"log/slog"
	"net/http"
	"smallscheduler/base"
	"smallscheduler/model"
)

// InitHttpRouter HTTP路由配置
func InitHttpRouter() {

	controller, err := NewController()
	if err != nil {
		base.Logger.Error(err.Error())
		return
	}

	router := easierweb.New(easierweb.RouterOptions{
		ErrorHandle:    errorHandle(),
		ResponseHandle: responseHandle(),
		Logger:         base.Logger,
		RootPath:       base.Config().Server.ContextPath,
	}).Use(logMiddleware())

	router.EasyGET("/tasks", controller.ListTask)
	router.EasyGET("/task/:id", controller.GetTask)

	router.EasyGET("/tags", controller.ListTag)
	router.EasyGET("/crons", controller.ListCron)

	router.EasyPOST("/task", controller.AddTask)
	router.EasyPUT("/task/:id", controller.EditTask)
	router.EasyDELETE("/task/:id", controller.DeleteTask)

	router.EasyGET("/execute/:id", controller.ExecuteTask)

	router.EasyGET("/records", controller.ListRecord)

	router.EasyGET("/health", controller.Health)
	router.EasyGET("/shutdown", controller.Shutdown)

	if base.Config().Server.ConsoleEnable {
		router.Static("/web/*filepath", "web")
	}

	err = router.Run(fmt.Sprintf(":%v", base.Config().Server.Port))
	if err != nil {
		base.Logger.Error(err.Error())
		return
	}
}

func errorHandle() easierweb.ErrorHandle {
	return func(ctx *easierweb.Context, err any) {
		ctx.WriteJSON(http.StatusBadRequest, model.CommonDTO{
			Ok:  false,
			Msg: fmt.Sprintf("unexpected error: %s", err),
		})
	}
}

func responseHandle() easierweb.ResponseHandle {
	return func(ctx *easierweb.Context, result any, err error) {
		if err != nil {
			ctx.WriteJSON(http.StatusBadRequest, model.CommonDTO{
				Ok:  false,
				Msg: err.Error(),
			})
		} else {
			if ctx.Method() == "POST" {
				ctx.WriteJSON(http.StatusCreated, result)
			} else {
				ctx.WriteJSON(http.StatusOK, result)
			}
		}
	}
}

func logMiddleware() easierweb.Handle {
	return func(ctx *easierweb.Context) {

		ctx.Next()

		path := ""
		query := ""
		body := ""
		result := ""

		if len(ctx.Path) > 0 {
			marshal, err := json.Marshal(ctx.Path)
			if err != nil {
				path = err.Error()
			} else {
				path = string(marshal)
			}
		}
		if len(ctx.Query) > 0 {
			marshal, err := json.Marshal(ctx.Query)
			if err != nil {
				query = err.Error()
			} else {
				query = string(marshal)
			}
		}
		sizeLimit := 1024 * 1024
		if len(ctx.Body) > 0 {
			if len(ctx.Body) > sizeLimit {
				body = "body is too large"
			} else {
				body = string(ctx.Body)
			}
		}
		if len(ctx.Result) > 0 {
			if len(ctx.Result) > sizeLimit {
				result = "result is too large"
			} else {
				result = string(ctx.Result)
			}
		}

		ctx.Logger.Info(ctx.Proto(), slog.String("method", ctx.Request.Method),
			slog.String("url", ctx.Request.URL.String()),
			slog.String("client", ctx.Request.RemoteAddr),
			slog.String("path", path),
			slog.String("query", query),
			slog.String("body", body),
			slog.Int("code", ctx.Code),
			slog.String("result", result))
	}
}
