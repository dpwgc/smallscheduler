package main

import (
	"smallscheduler/api"
	"smallscheduler/base"
	"smallscheduler/core"
)

func main() {
	base.InitLog()
	base.InitConfig()
	core.InitWorkers()
	api.InitHttpRouter()
}
