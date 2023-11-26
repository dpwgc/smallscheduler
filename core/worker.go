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
var taskCachePool sync.Map
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
	var version uint64 = 0
	taskEditVersion.Store(version)
	go func() {
		for {
			loadWorkers()
			time.Sleep(TaskListRefreshCycle)
		}
	}()
}

// 加载工作者列表
func loadWorkers() {
	latestVersion, err := service.GetTaskEditVersion()
	if err != nil {
		log.Println(base.LogErrorTag, err)
		return
	}
	if latestVersion == taskEditVersion.Load().(int64) {
		return
	}
	taskEditVersion.Store(latestVersion)
	// 获取当前系统中的所有任务的cron表达式
	cronList, err := service.ListStartedCron()
	if err != nil {
		log.Println(base.LogErrorTag, err)
		return
	}
	//为每个cron表达式生成一个工作者
	for _, cronStr := range cronList {
		loadWorker(cronStr)
	}
}

func loadWorker(cronStr string) {
	//判断是否已经存在工作者
	check, _ := workerFactory.Load(cronStr)
	if check != nil {
		return
	}
	//创建工作者（协程定时任务）
	worker := newWithSeconds()
	//装配函数
	_, err := worker.AddFunc(cronStr, func() {
		execute(cronStr)
	})
	if err != nil {
		log.Println(base.LogErrorTag, err)
		return
	}
	//获取该cron下的所有任务
	taskList, err := service.ListStartedTaskByCron(cronStr)
	if err != nil {
		log.Println(base.LogErrorTag, err)
		return
	}
	if len(taskList) == 0 {
		return
	}
	//将该工作者装入工作者列表
	taskCachePool.Store(cronStr, taskList)
	//启动工作者
	worker.Start()
	//将该工作者装入工作者列表
	workerFactory.Store(cronStr, worker)
}

// 返回一个支持至 秒 级别的 cron
func newWithSeconds() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}
