package core

import (
	"alisa-dispatch-center/src/common"
	"alisa-dispatch-center/src/constant"
	"alisa-dispatch-center/src/storage"
	"fmt"
	"github.com/robfig/cron/v3"
)

// 批量执行任务
func performTasks(cronStr string) {
	//获取该cron的所有任务
	taskList, err := service.ListTaskToServer(cronStr)
	if err != nil {
		common.Log.Println(constant.LogErrorTag, err)
		return
	}
	//如果任务列表长度为0，则删除该工作者
	if len(taskList) == 0 {
		worker, _ := workerMap.Load(cronStr)
		worker.(*cron.Cron).Stop()
		workerMap.Delete(cronStr)
	}
	//循环请求
	for _, task := range taskList {
		request(task)
	}
}

// TODO
func request(task storage.Task) {
	go func() {
		fmt.Println(task)
	}()
}
