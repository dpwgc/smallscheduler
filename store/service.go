package store

import (
	"smallscheduler/model"
	"smallscheduler/store/db"
)

func NewService() (*Service, error) {
	repository, err := db.NewRepository()
	return &Service{
		repository: repository,
	}, err
}

type Service struct {
	repository *db.Repository
}

func (s *Service) ListTask(name string, tag string, spec string, status int, pageIndex int, pageSize int) ([]model.Task, int64, error) {
	return s.repository.ListTask(name, tag, spec, status, pageIndex, pageSize)
}

func (s *Service) GetTask(id int64) (model.Task, error) {
	return s.repository.GetTask(id)
}

func (s *Service) ListTagCount(status int) ([]model.TagCount, error) {
	return s.repository.ListTagCount(status)
}

func (s *Service) ListSpecCount(status int) ([]model.SpecCount, error) {
	return s.repository.ListSpecCount(status)
}

func (s *Service) ListStartedTaskBySpec(spec string) ([]model.Task, error) {
	return s.repository.ListStartedTaskBySpec(spec)
}

func (s *Service) ListStartedSpec() ([]string, error) {
	return s.repository.ListStartedSpec()
}

func (s *Service) TryExecuteTask(task model.Task) (int64, error) {
	return s.repository.TryExecuteTask(task)
}

func (s *Service) AddTask(task model.Task) (int64, error) {
	err := s.repository.ChangeTaskEditVersion()
	if err != nil {
		return 0, err
	}
	return s.repository.AddTask(task)
}

func (s *Service) EditTask(task model.Task) error {
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

func (s *Service) ListRecord(shard string, taskId int64, code int, startTime string, endTime string, pageIndex int, pageSize int) ([]model.Record, int64, error) {
	return s.repository.ListRecord(shard, taskId, code, startTime, endTime, pageIndex, pageSize)
}

func (s *Service) GetTaskEditVersion() (int64, error) {
	return s.repository.GetTaskEditVersion()
}
