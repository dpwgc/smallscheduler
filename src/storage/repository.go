package storage

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

func initRepository() Repository {
	return &repositoryImpl{
		DB: DBClient(),
	}
}

type Repository interface {
	SelectTaskByApp(appId uint64, env uint8, name string) ([]Task, error)
	SelectTaskByPartition(partition string) ([]Task, error)
	InsertTask(task Task) error
	UpdateTask(task Task) error
}

type repositoryImpl struct {
	DB *gorm.DB
}

func (r *repositoryImpl) SelectTaskByApp(appId uint64, env uint8, name string) ([]Task, error) {
	var taskList []Task
	sql := r.DB.Where("app_id = ? and status = ?", appId, 1)
	if env != 0 {
		sql = sql.Where("env = ?", env)
	}
	if len(name) > 0 {
		sql = sql.Where(fmt.Sprintf(" name like %q ", "%"+name+"%"))
	}
	err := sql.Find(&taskList).Error
	return taskList, err
}

func (r *repositoryImpl) SelectTaskByPartition(partition string) ([]Task, error) {
	var taskList []Task
	err := r.DB.Where("partition = ? and status = ?", partition, 1).Find(&taskList).Error
	return taskList, err
}

func (r *repositoryImpl) InsertTask(task Task) error {
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	task.Status = 1
	return DBClient().Create(&task).Error
}

func (r *repositoryImpl) UpdateTask(task Task) error {
	task.UpdatedAt = time.Now()
	return r.DB.Updates(&task).Error
}
