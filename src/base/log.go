package base

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"log"
	"time"
)

const (
	LogErrorTag         = "[ERROR]"
	LogWarnTag          = "[WARN]"
	LogInfoTag          = "[INFO]"
	LogFilePath         = "./logs/runtime.%Y-%m-%d.log"
	LogFileRotationTime = time.Duration(24) * time.Hour
)

func InitLog() {
	writer, _ := rotatelogs.New(
		LogFilePath,
		rotatelogs.WithMaxAge(time.Duration(Config.Log.FileMaxAge*24)*time.Hour),
		rotatelogs.WithRotationTime(LogFileRotationTime),
	)
	log.SetFlags(log.Lshortfile)
	log.SetOutput(writer)
}
