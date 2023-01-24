package constant

import "time"

const (
	HttpRequestSuccessCode = 10000
	HttpRequestFailCode    = 20001
	HttpUriPrefix          = "/ui/v2"

	ConfigFilePath = "./config.yaml"

	LogErrorTag          = "[ERROR]"
	LogWarnTag           = "[WARN]"
	LogInfoTag           = "[INFO]"
	LogFilePath          = "./logs/runtime.%Y-%m-%d.log"
	LogFileRotationTime  = time.Duration(24) * time.Hour
	TaskListRefreshCycle = time.Duration(1) * time.Second
)
