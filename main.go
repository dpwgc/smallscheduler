package main

import (
	"fmt"
	"smallscheduler/api"
	"smallscheduler/base"
	"smallscheduler/core"
)

func main() {

	fmt.Print("   _____                 _ _    _____      _              _       _           \n  / ____|               | | |  / ____|    | |            | |     | |          \n | (___  _ __ ___   __ _| | | | (___   ___| |__   ___  __| |_   _| | ___ _ __ \n  \\___ \\| '_ ` _ \\ / _` | | |  \\___ \\ / __| '_ \\ / _ \\/ _` | | | | |/ _ \\ '__|\n  ____) | | | | | | (_| | | |  ____) | (__| | | |  __/ (_| | |_| | |  __/ |   \n |_____/|_| |_| |_|\\__,_|_|_| |_____/ \\___|_| |_|\\___|\\__,_|\\__,_|_|\\___|_|   \n ")
	fmt.Printf("\033[1;32;40m%s\033[0m\n", " === Small Scheduler (v1.0.0) === ")

	base.InitConfig()

	fmt.Println("[Config]:", base.ConfigJson())

	base.InitLog()
	core.InitWorkers()

	fmt.Println("[Console]:", fmt.Sprintf("http://localhost:%v%s/web/index.html", base.Config().Server.Port, base.Config().Server.ContextPath))

	api.InitHttpRouter()
}
