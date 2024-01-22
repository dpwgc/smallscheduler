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

	controller, err := NewController()
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

	var adminGroup *easierweb.Group
	if base.Config().Server.Auth {
		adminGroup = router.Group("", basicAuthMiddleware())
	} else {
		adminGroup = router.Group("")
	}

	adminGroup.EasyGET("/tasks", controller.ListTask)
	adminGroup.EasyGET("/task/:id", controller.GetTask)

	adminGroup.EasyGET("/tags", controller.ListTag)
	adminGroup.EasyGET("/crons", controller.ListCron)

	adminGroup.EasyPOST("/task", controller.AddTask)
	adminGroup.EasyPUT("/task/:id", controller.EditTask)
	adminGroup.EasyDELETE("/task/:id", controller.DeleteTask)

	adminGroup.EasyGET("/execute/:id", controller.ExecuteTask)

	adminGroup.EasyGET("/records", controller.ListRecord)

	serverGroup := router.Group("")
	serverGroup.EasyGET("/health", controller.Health)
	serverGroup.EasyGET("/shutdown", controller.Shutdown)

	if base.Config().Server.ConsoleEnable {
		router.Static("/web/*filepath", "web")
	}

	host := fmt.Sprintf("%s:%v", base.Config().Server.Addr, base.Config().Server.Port)
	if base.Config().Server.TLS {
		err = router.RunTLS(host, base.Config().Server.CertFile, base.Config().Server.KeyFile, &tls.Config{})
	} else {
		err = router.Run(host)
	}

	if err != nil {
		base.Logger.Error(err.Error())
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

func basicAuthMiddleware() easierweb.Handle {
	return func(ctx *easierweb.Context) {
		user, pass, ok := ctx.Request.BasicAuth()
		if !ok {
			ctx.WriteJSON(http.StatusUnauthorized, model.CommonDTO{
				Ok:  false,
				Msg: "basic authentication failed",
			})
			return
		}
		if base.Config().Server.Username != user || base.Config().Server.Password != pass {
			ctx.WriteJSON(http.StatusUnauthorized, model.CommonDTO{
				Ok:  false,
				Msg: "username or password error",
			})
			return
		}
		// 真正需要处理的业务
		ctx.Next()
	}
}
