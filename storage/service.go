package storage

import (
	"smallscheduler/model"
	"smallscheduler/storage/rdb"
	"strings"
)

func NewService() (*Service, error) {
	repository, err := rdb.NewRepository()
	return &Service{
		repository: repository,
	}, err
}

type Service struct {
	repository *rdb.Repository
}

func (s *Service) ListTask(name string, tag string, cron string, status int, pageIndex int, pageSize int) ([]model.Task, int64, error) {
	return s.repository.ListTask(name, tag, cron, status, pageIndex, pageSize)
}

func (s *Service) GetTask(id int64) (model.Task, error) {
	return s.repository.GetTask(id)
}

func (s *Service) ListTagCount(status int) ([]model.TagCount, error) {
	return s.repository.ListTagCount(status)
}

func (s *Service) ListCronCount(status int) ([]model.CronCount, error) {
	return s.repository.ListCronCount(status)
}

func (s *Service) ListStartedTaskByCron(cron string) ([]model.Task, error) {
	return s.repository.ListStartedTaskByCron(cron)
}

func (s *Service) ListStartedCron() ([]string, error) {
	return s.repository.ListStartedCron()
}

func (s *Service) TryExecuteTask(task model.Task) (int64, error) {
	return s.repository.TryExecuteTask(task)
}

func (s *Service) AddTask(task model.Task) (int64, error) {
	task.Cron = strings.TrimSpace(task.Cron)
	task.Tag = strings.TrimSpace(task.Tag)
	task.Name = strings.TrimSpace(task.Name)
	task.Url = strings.TrimSpace(task.Url)
	err := s.repository.ChangeTaskEditVersion()
	if err != nil {
		return 0, err
	}
	return s.repository.AddTask(task)
}

func (s *Service) EditTask(task model.Task) error {
	task.Cron = strings.TrimSpace(task.Cron)
	task.Tag = strings.TrimSpace(task.Tag)
	task.Name = strings.TrimSpace(task.Name)
	task.Url = strings.TrimSpace(task.Url)
	err := s.repository.ChangeTaskEditVersion()
	if err != nil {
		return err
	}
	return s.repository.EditTask(task)
}

func (s *Service) DeleteTask(id int64) error {
	return s.repository.DeleteTask(id)
}

func (s *Service) AddRecord(record model.Record) error {
	return s.repository.AddRecord(record)
}

func (s *Service) ListRecord(sharding string, taskId int64, code int, startTime string, endTime string, pageIndex int, pageSize int) ([]model.Record, int64, error) {
	return s.repository.ListRecord(sharding, taskId, code, startTime, endTime, pageIndex, pageSize)
}

func (s *Service) GetTaskEditVersion() (int64, error) {
	return s.repository.GetTaskEditVersion()
}
