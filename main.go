package main

import (
	"smallscheduler/api"
	"smallscheduler/base"
	"smallscheduler/core"
)

func main() {
	base.InitConfig()
	base.InitLog()
	core.InitWorkers()
	api.InitHttpRouter()
}
