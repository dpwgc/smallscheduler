package core

import (
	"alisa-dispatch-center/src/storage"
	"fmt"
	"time"
)

func initRepository() Repository {
	return &defaultRepository{}
}

type Repository interface {
	SelectTaskByApp(appId uint64, env uint8, name string) ([]storage.Task, error)
	SelectTaskByCron(cron string) ([]storage.Task, error)
	InsertTask(task storage.Task) error
	UpdateTask(task storage.Task) error
}

type defaultRepository struct{}

func (r *defaultRepository) SelectTaskByApp(appId uint64, env uint8, name string) ([]storage.Task, error) {
	var taskList []storage.Task
	sql := storage.DBClient().Where("app_id = ? and status = ?", appId, 1)
	if env != 0 {
		sql = sql.Where("env = ?", env)
	}
	if len(name) > 0 {
		sql = sql.Where(fmt.Sprintf(" name like %q ", "%"+name+"%"))
	}
	err := sql.Find(&taskList).Error
	return taskList, err
}

func (r *defaultRepository) SelectTaskByCron(cron string) ([]storage.Task, error) {
	var taskList []storage.Task
	err := storage.DBClient().Where("cron = ? and status = ?", cron, 1).Find(&taskList).Error
	return taskList, err
}

func (r *defaultRepository) InsertTask(task storage.Task) error {
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	task.Status = 1
	return storage.DBClient().Create(&task).Error
}

func (r *defaultRepository) UpdateTask(task storage.Task) error {
	task.UpdatedAt = time.Now()
	return storage.DBClient().Updates(&task).Error
}
