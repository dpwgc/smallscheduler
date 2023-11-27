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

func (s *Service) ListTask(name string, status int, pageIndex int, pageSize int) ([]Task, int64, error) {
	return s.repository.ListTask(name, status, pageIndex, pageSize)
}

func (s *Service) ListStartedTaskByCron(cron string) ([]Task, error) {
	return s.repository.ListStartedTaskByCron(cron)
}

func (s *Service) ListStartedCron() ([]string, error) {
	return s.repository.ListStartedCron()
}

func (s *Service) ExecuteTask(id int64) (int64, error) {
	return s.repository.ExecuteTask(id)
}

func (s *Service) SaveTask(task Task) error {
	err := s.repository.AddTaskEditVersion()
	if err != nil {
		return err
	}
	return s.repository.SaveTask(task)
}

func (s *Service) RemoveTask(id int64) error {
	return s.repository.RemoveTask(id)
}

func (s *Service) SaveRecord(record Record) error {
	return s.repository.SaveRecord(record)
}

func (s *Service) ListRecord(taskId int64, status int, startTime string, endTime string, pageIndex int, pageSize int) ([]Record, int64, error) {
	return s.repository.ListRecord(taskId, status, startTime, endTime, pageIndex, pageSize)
}

func (s *Service) GetTaskEditVersion() (int64, error) {
	return s.repository.GetTaskEditVersion()
}
