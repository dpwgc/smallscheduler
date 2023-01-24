package core

import (
	"alisa-dispatch-center/src/common"
	"alisa-dispatch-center/src/constant"
	"alisa-dispatch-center/src/storage"
	"github.com/robfig/cron/v3"
	"sync"
	"time"
)

var workerMap sync.Map

var service storage.Service

// InitWorkers 初始化工作者列表
func InitWorkers() {
	service = storage.InitService()
	go func() {
		for {
			loadWorkers()
			time.Sleep(constant.TaskListRefreshCycle)
		}
	}()
}

// 加载工作者列表
func loadWorkers() {
	// 获取当前系统中的所有任务的cron表达式
	cronList, err := service.ListCron()
	if err != nil {
		common.Log.Println(constant.LogErrorTag, err)
		return
	}
	//为每个cron表达式生成一个工作者
	for _, cronStr := range cronList {
		//判断是否已经存在工作者
		check, _ := workerMap.Load(cronStr)
		if check != nil {
			continue
		}
		//创建工作者（协程定时任务）
		worker := newWithSeconds()
		//装配函数
		_, err = worker.AddFunc(cronStr, func() {
			performTasks(cronStr)
		})
		if err != nil {
			common.Log.Println(constant.LogErrorTag, err)
			continue
		}
		//启动工作者
		worker.Start()
		//将该工作者装入工作者列表
		workerMap.Store(cronStr, worker)
	}
}

// 返回一个支持至 秒 级别的 cron
func newWithSeconds() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}
