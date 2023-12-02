package core

import (
	"github.com/robfig/cron/v3"
	"log"
	"smallscheduler/base"
	"smallscheduler/storage"
	"sync"
	"sync/atomic"
	"time"
)

const (
	TaskListRefreshCycle = time.Duration(1) * time.Second
)

var workerFactory sync.Map
var service *storage.Service
var taskEditVersion atomic.Value

// InitWorkers 初始化工作者列表
func InitWorkers() {
	s, err := storage.NewService()
	if err != nil {
		log.Fatal(base.LogErrorTag, err)
		return
	}
	service = s
	var version int64 = -1
	taskEditVersion.Store(version)
	go func() {
		for {
			latestVersion, err := service.GetTaskEditVersion()
			if err != nil {
				log.Println(base.LogErrorTag, err)
				continue
			}
			if latestVersion == taskEditVersion.Load().(int64) {
				continue
			}
			// 加载两次，避免被覆盖
			loadWorkers()
			time.Sleep(TaskListRefreshCycle)
			loadWorkers()
			taskEditVersion.Store(latestVersion)
			time.Sleep(TaskListRefreshCycle)
		}
	}()
}

// 加载工作者列表
func loadWorkers() {
	// 获取当前系统中的所有任务的cron表达式
	cronList, err := service.ListStartedCron()
	if err != nil {
		log.Println(base.LogErrorTag, err)
	}
	//为每个cron表达式生成一个工作者
	for _, cronStr := range cronList {
		loadWorker(cronStr)
	}
}

func loadWorker(cronStr string) {
	//判断是否已经存在工作者
	oldWorker, _ := workerFactory.Load(cronStr)
	if oldWorker != nil {
		return
	}
	//创建工作者（协程定时任务）
	worker := NewCronWorker()
	//装配函数
	_, err := worker.AddFunc(cronStr, func() {
		execute(cronStr)
	})
	if err != nil {
		log.Println(base.LogErrorTag, err)
		return
	}
	//启动工作者
	worker.Start()
	//将该工作者装入工作者列表
	workerFactory.Store(cronStr, worker)
}

// NewCronWorker 返回一个支持至 秒 级别的 cron
func NewCronWorker() *cron.Cron {
	return cron.New(cron.WithSeconds(), cron.WithChain())
}
