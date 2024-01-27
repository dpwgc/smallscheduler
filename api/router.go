package api

import (
	"crypto/tls"
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

	adapter, err := newAdapter()
	if err != nil {
		base.Logger.Error(err.Error())
		return
	}

	base.Logger.Info("start http server")

	router := easierweb.New(easierweb.RouterOptions{
		ErrorHandle:    errorHandle(),
		ResponseHandle: responseHandle(),
		Logger:         base.Logger,
		RootPath:       base.Config().Server.ContextPath,
	}).Use(logMiddleware())

	router.EasyGET("/tasks", adapter.ListTask)
	router.EasyGET("/task/:id", adapter.GetTask)

	router.EasyGET("/tags", adapter.ListTag)
	router.EasyGET("/specs", adapter.ListSpec)

	router.EasyPOST("/task", adapter.AddTask)
	router.EasyPUT("/task/:id", adapter.EditTask)
	router.EasyDELETE("/task/:id", adapter.DeleteTask)

	router.EasyGET("/execute/:id", adapter.ExecuteTask)

	router.EasyGET("/records", adapter.ListRecord)

	router.EasyGET("/health", adapter.Health)
	router.EasyGET("/shutdown", adapter.Shutdown)

	if base.Config().Server.ConsoleEnable {
		router.Static("/web/*filepath", "web")
	}

	host := fmt.Sprintf("%s:%v", base.Config().Server.Addr, base.Config().Server.Port)
	if base.Config().Server.TLS {
		err = router.RunTLS(host, base.Config().Server.CertFile, base.Config().Server.KeyFile, &tls.Config{})
	} else {
		err = router.Run(host)
	}
	base.Logger.Error(err.Error())
}

func errorHandle() easierweb.ErrorHandle {
	return func(ctx *easierweb.Context, err any) {
		ctx.WriteJSON(http.StatusBadRequest, model.CommonDTO{
			Msg: fmt.Sprintf("unexpected error: %s", err),
		})
	}
}

func responseHandle() easierweb.ResponseHandle {
	return func(ctx *easierweb.Context, result any, err error) {
		if err != nil {
			ctx.WriteJSON(http.StatusBadRequest, model.CommonDTO{
				Msg: err.Error(),
			})
		} else {
			if result == nil {
				ctx.NoContent(http.StatusNoContent)
				return
			}
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

		if ctx.Proto() == "/health" {
			return
		}

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
