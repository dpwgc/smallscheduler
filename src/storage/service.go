package storage

func InitService() Service {
	return &ServiceImpl{
		repository: initRepository(),
	}
}

type Service interface {
	ListTaskToUser(appId uint64, env uint8, name string) ([]Task, error)
	ListTaskToServer(partition string) ([]Task, error)
	SaveTask(task Task) error
}

type ServiceImpl struct {
	repository Repository
}

func (s *ServiceImpl) ListTaskToUser(appId uint64, env uint8, name string) ([]Task, error) {
	return s.repository.SelectTaskByApp(appId, env, name)
}

func (s *ServiceImpl) ListTaskToServer(partition string) ([]Task, error) {
	return s.repository.SelectTaskByPartition(partition)
}

func (s *ServiceImpl) SaveTask(task Task) error {
	// update
	if task.Id > 0 {
		return s.repository.UpdateTask(task)
	}
	// insert
	return s.repository.InsertTask(task)
}
