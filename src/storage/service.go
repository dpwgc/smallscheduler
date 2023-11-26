package storage

import (
	"alisa-dispatch-center/src/storage/rdb"
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

func (s *Service) ListTaskToUser(name string, status int, pageIndex int, pageSize int) ([]rdb.Task, int64, error) {
	return s.repository.ListTaskToUser(name, status, pageIndex, pageSize)
}

func (s *Service) ListTaskToServer(cron string) ([]rdb.Task, error) {
	return s.repository.ListTaskToServer(cron)
}

func (s *Service) ListCron() ([]string, error) {
	return s.repository.ListCron()
}

func (s *Service) DoTask(task rdb.Task) (int64, error) {
	return s.repository.DoTask(task)
}

func (s *Service) SaveTask(task rdb.Task) error {
	return s.repository.SaveTask(task)
}

func (s *Service) SaveRecord(record rdb.Record) error {
	return s.repository.SaveRecord(record)
}
