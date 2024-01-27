package job

import (
	"errors"
	"github.com/robfig/cron/v3"
	"smallscheduler/base"
	"smallscheduler/store"
	"sync"
	"sync/atomic"
	"time"
)

const (
	TaskListRefreshCycle = time.Duration(1) * time.Second
)

var workerFactory sync.Map
var service *store.Service
var taskEditVersion atomic.Value

// InitWorkers 初始化工作者列表
func InitWorkers() {
	s, err := store.NewService()
	if err != nil {
		base.Logger.Error(err.Error())
		return
	}
	service = s
	var version int64 = -1
	taskEditVersion.Store(version)
	go func() {
		for {
			time.Sleep(TaskListRefreshCycle)
			latestVersion, err := service.GetTaskEditVersion()
			if err != nil {
				base.Logger.Error(err.Error())
				continue
			}
			if latestVersion == taskEditVersion.Load().(int64) {
				continue
			}

			// 加载两次，避免被覆盖
			loadWorkers()
			time.Sleep(TaskListRefreshCycle)
			loadWorkers()

			// 更新版本号
			taskEditVersion.Store(latestVersion)

			base.Logger.Info("scheduler workers data is loaded successfully")
		}
	}()
	base.Logger.Info("scheduler background program is loaded successfully")
}

// 加载工作者列表
func loadWorkers() {
	// 获取当前系统中的所有任务的cron表达式
	specList, err := service.ListStartedSpec()
	if err != nil {
		base.Logger.Error(err.Error())
	}
	//为每个cron表达式生成一个工作者
	for _, spec := range specList {
		loadWorker(spec)
	}
}

func loadWorker(spec string) {
	//判断是否已经存在工作者
	oldWorker, _ := workerFactory.Load(spec)
	if oldWorker != nil {
		return
	}
	//创建工作者（协程定时任务）
	worker := NewCronWorker()
	//装配函数
	_, err := worker.AddFunc(spec, func() {
		scheduled(spec)
	})
	if err != nil {
		base.Logger.Error(err.Error())
		return
	}
	//启动工作者
	worker.Start()
	//将该工作者装入工作者列表
	workerFactory.Store(spec, worker)

	base.Logger.Info("a new worker has been created")
}

// NewCronWorker 返回一个支持至 秒 级别的 cron
func NewCronWorker() *cron.Cron {
	return cron.New(cron.WithSeconds(), cron.WithChain())
}

func VerifySpec(spec string) error {
	checkWorker := NewCronWorker()
	defer func() {
		checkWorker = nil
	}()
	_, err := checkWorker.AddFunc(spec, func() {})
	if err != nil {
		return errors.New("spec error: " + err.Error())
	}
	return nil
}
