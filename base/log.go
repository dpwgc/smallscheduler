package base

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
	"time"
)

var Logger *slog.Logger

func InitLog() {
	r := &lumberjack.Logger{
		Filename:   Config().Log.Path,
		LocalTime:  true,
		MaxSize:    Config().Log.MaxSize,
		MaxAge:     Config().Log.MaxAge,
		MaxBackups: Config().Log.MaxBackups,
		Compress:   false,
	}
	Logger = slog.New(slog.NewJSONHandler(r, &slog.HandlerOptions{
		AddSource: true, // 输出日志语句的位置信息
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey { // 格式化 key 为 "time" 的属性值
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format(time.DateTime))
				}
			}
			return a
		},
	}))
	Logger.Info("log module is loaded successfully")
}
