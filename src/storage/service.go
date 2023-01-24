package storage

func InitService() Service {
	return &ServiceImpl{
		repository: initRepository(),
	}
}

type Service interface {
	ListTaskToUser(appId uint64, env uint8, name string) ([]Task, error)
	ListTaskToServer(cron string) ([]Task, error)
	ListCron() ([]string, error)
	SaveTask(task Task) error
}

type ServiceImpl struct {
	repository Repository
}

func (s *ServiceImpl) ListTaskToUser(appId uint64, env uint8, name string) ([]Task, error) {
	return s.repository.SelectTaskToUser(appId, env, name)
}

func (s *ServiceImpl) ListTaskToServer(cron string) ([]Task, error) {
	return s.repository.SelectTaskByServer(cron)
}

func (s *ServiceImpl) ListCron() ([]string, error) {
	return s.repository.SelectCron()
}

func (s *ServiceImpl) SaveTask(task Task) error {
	// update
	if task.Id > 0 {
		return s.repository.UpdateTask(task)
	}
	// insert
	return s.repository.InsertTask(task)
}
