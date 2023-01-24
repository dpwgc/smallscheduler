package storage

import (
	"alisa-dispatch-center/src/common"
	"fmt"
	"gorm.io/gorm"
	"time"
)

func initRepository() Repository {
	return &repositoryImpl{
		DB: initDB(),
	}
}

type Repository interface {
	SelectTaskByApp(appId uint64, env uint8, name string) ([]Task, error)
	SelectTaskByCron(cron string) ([]Task, error)
	SelectCron() ([]string, error)
	InsertTask(task Task) error
	UpdateTask(task Task) error
}

type repositoryImpl struct {
	DB *gorm.DB
}

func (r *repositoryImpl) SelectTaskByApp(appId uint64, env uint8, name string) ([]Task, error) {
	var taskList []Task
	sql := r.DB.Model(&Task{}).Where("app_id = ? and status = ?", appId, 1)
	if env != 0 {
		sql = sql.Where("env = ?", env)
	}
	if len(name) > 0 {
		sql = sql.Where(fmt.Sprintf(" name like %q ", "%"+name+"%"))
	}
	err := sql.Find(&taskList).Error
	return taskList, err
}

func (r *repositoryImpl) SelectTaskByCron(cron string) ([]Task, error) {
	var taskList []Task
	err := r.DB.Model(&Task{}).Where("cron = ? and partition = ? and status = ?", cron, common.Config.Server.Partition, 1).Find(&taskList).Error
	return taskList, err
}

func (r *repositoryImpl) SelectCron() ([]string, error) {
	var taskList []Task
	var cronList []string
	err := r.DB.Model(&Task{}).Where("status = ?", 1).Find(&taskList).Error
	if err != nil {
		return nil, err
	}
	for _, task := range taskList {
		cronList = append(cronList, task.Cron)
	}
	return cronList, err
}

func (r *repositoryImpl) InsertTask(task Task) error {
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	task.Status = 1
	return r.DB.Create(&task).Error
}

func (r *repositoryImpl) UpdateTask(task Task) error {
	task.UpdatedAt = time.Now()
	return r.DB.Updates(&task).Error
}
