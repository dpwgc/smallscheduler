package core

import (
	"alisa-dispatch-center/src/base"
	"alisa-dispatch-center/src/storage"
	"github.com/robfig/cron/v3"
	"log"
	"sync"
	"time"
)

const (
	TaskListRefreshCycle = time.Duration(1) * time.Second
)

var workerMap sync.Map

var service *storage.Service

// InitWorkers 初始化工作者列表
func InitWorkers() {
	s, err := storage.NewService()
	if err != nil {
		log.Fatal(base.LogErrorTag, err)
		return
	}
	service = s
	go func() {
		for {
			loadWorkers()
			time.Sleep(TaskListRefreshCycle)
		}
	}()
}

// 加载工作者列表
func loadWorkers() {
	// 获取当前系统中的所有任务的cron表达式
	cronList, err := service.ListCron()
	if err != nil {
		log.Println(base.LogErrorTag, err)
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
			execute(cronStr)
		})
		if err != nil {
			log.Println(base.LogErrorTag, err)
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
