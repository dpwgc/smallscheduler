package main

import (
	"alisa-dispatch-center/src/api"
	"alisa-dispatch-center/src/base"
	"alisa-dispatch-center/src/core"
)

func main() {
	base.InitLog()
	base.InitConfig()
	core.InitWorkers()
	api.InitHttpRouter()
}
