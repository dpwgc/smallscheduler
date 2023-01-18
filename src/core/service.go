package core

import (
	"alisa-dispatch-center/src/storage"
)

func InitService() Service {
	repository = initRepository()
	return &defaultService{}
}

type Service interface {
	ListTaskToUser(appId uint64, env uint8, name string) ([]storage.Task, error)
	ListTaskToServer(cron string) ([]storage.Task, error)
	SaveTask(task storage.Task) error
}

type defaultService struct{}

var repository Repository

func (s *defaultService) ListTaskToUser(appId uint64, env uint8, name string) ([]storage.Task, error) {
	return repository.SelectTaskByApp(appId, env, name)
}

func (s *defaultService) ListTaskToServer(cron string) ([]storage.Task, error) {
	return repository.SelectTaskByCron(cron)
}

func (s *defaultService) SaveTask(task storage.Task) error {
	// update
	if task.Id > 0 {
		return repository.UpdateTask(task)
	}
	// insert
	return repository.InsertTask(task)
}
