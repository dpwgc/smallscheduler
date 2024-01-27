package main

import (
	"fmt"
	"smallscheduler/api"
	"smallscheduler/base"
	"smallscheduler/core"
	"strings"
)

func main() {

	fmt.Print("   _____                 _ _    _____      _              _       _           \n  / ____|               | | |  / ____|    | |            | |     | |          \n | (___  _ __ ___   __ _| | | | (___   ___| |__   ___  __| |_   _| | ___ _ __ \n  \\___ \\| '_ ` _ \\ / _` | | |  \\___ \\ / __| '_ \\ / _ \\/ _` | | | | |/ _ \\ '__|\n  ____) | | | | | | (_| | | |  ____) | (__| | | |  __/ (_| | |_| | |  __/ |   \n |_____/|_| |_| |_|\\__,_|_|_| |_____/ \\___|_| |_|\\___|\\__,_|\\__,_|_|\\___|_|   \n ")
	fmt.Printf("\033[1;32;40m%s\033[0m\n", " === Small Scheduler (v1.0.1) === ")

	base.InitConfig()

	fmt.Println("[Server]", strings.ReplaceAll(base.ServerYaml(), "\n", " "))
	fmt.Println("[Database]", strings.ReplaceAll(base.DbYaml(), "\n", " "))
	fmt.Println("[Log]", strings.ReplaceAll(base.LogYaml(), "\n", " "))

	base.InitLog()
	core.InitWorkers()

	fmt.Println("[Console]", fmt.Sprintf("http://localhost:%v%s/web/index.html", base.Config().Server.Port, base.Config().Server.ContextPath))

	api.InitHttpRouter()

	fmt.Printf("\033[1;31;40m%s\033[0m\n", fmt.Sprintf(" <shutdown> log file: %s ", base.Config().Log.Path+"/small-scheduler.log"))
}
