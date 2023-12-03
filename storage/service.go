package storage

func NewService() (*Service, error) {
	repository, err := NewRepository()
	return &Service{
		repository: repository,
	}, err
}

type Service struct {
	repository *Repository
}

func (s *Service) ListTask(name string, cron string, status int, pageIndex int, pageSize int) ([]Task, int64, error) {
	return s.repository.ListTask(name, cron, status, pageIndex, pageSize)
}

func (s *Service) GetTask(id int64) (Task, error) {
	return s.repository.GetTask(id)
}

func (s *Service) ListStartedTaskByCron(cron string) ([]Task, error) {
	return s.repository.ListStartedTaskByCron(cron)
}

func (s *Service) ListStartedCron() ([]string, error) {
	return s.repository.ListStartedCron()
}

func (s *Service) TryExecuteTask(task Task) (int64, error) {
	return s.repository.TryExecuteTask(task)
}

func (s *Service) AddTask(task Task) (int64, error) {
	err := s.repository.ChangeTaskEditVersion()
	if err != nil {
		return 0, err
	}
	return s.repository.AddTask(task)
}

func (s *Service) EditTask(task Task) error {
	err := s.repository.ChangeTaskEditVersion()
	if err != nil {
		return err
	}
	return s.repository.EditTask(task)
}

func (s *Service) DeleteTask(id int64) error {
	return s.repository.DeleteTask(id)
}

func (s *Service) AddRecord(record Record) error {
	return s.repository.AddRecord(record)
}

func (s *Service) ListRecord(taskId int64, sharding string, pageIndex int, pageSize int) ([]Record, int64, error) {
	return s.repository.ListRecord(taskId, sharding, pageIndex, pageSize)
}

func (s *Service) GetTaskEditVersion() (int64, error) {
	return s.repository.GetTaskEditVersion()
}
