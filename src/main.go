package main

import (
	"alisa-dispatch-center/src/common"
	"alisa-dispatch-center/src/core"
	"alisa-dispatch-center/src/ui"
)

func main() {
	common.InitLog()
	common.InitConfig()
	core.InitWorkers()
	ui.InitHttpRouter()
}
