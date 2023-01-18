package common

import (
	"alisa-dispatch-center/src/constant"
	"github.com/lestrrat-go/file-rotatelogs"
	"log"
	"time"
)

var Log *log.Logger

func InitLog() {
	writer, _ := rotatelogs.New(
		constant.LogFilePath,
		rotatelogs.WithMaxAge(time.Duration(Config.Log.FileMaxAge*24)*time.Hour),
		rotatelogs.WithRotationTime(constant.LogFileRotationTime),
	)
	log.SetOutput(writer)
	Log = log.Default()
}
