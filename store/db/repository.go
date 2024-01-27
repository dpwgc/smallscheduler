package db

import (
	"fmt"
	"gorm.io/gorm"
	"smallscheduler/base"
	"smallscheduler/model"
	"time"
)

func NewRepository() (*Repository, error) {
	db, err := newClient()
	return &Repository{
		DB: db,
	}, err
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) ListTask(name string, tag string, spec string, status int, pageIndex int, pageSize int) ([]model.Task, int64, error) {
	var taskList []model.Task
	var total int64
	sql := r.DB.Model(&model.Task{})
	if status > 0 {
		sql = sql.Where("status = ?", status)
	}
	if len(name) > 0 {
		sql = sql.Where(fmt.Sprintf("name like %q", "%"+name+"%"))
	}
	if len(tag) > 0 {
		sql = sql.Where("tag = ?", tag)
	}
	if len(spec) > 0 {
		sql = sql.Where("spec = ?", spec)
	}
	sql.Count(&total)
	sql = sql.Order("id desc").Limit(pageSize).Offset((pageIndex - 1) * pageSize)
	err := sql.Find(&taskList).Error
	return taskList, total, err
}

func (r *Repository) GetTask(id int64) (model.Task, error) {
	var task model.Task
	err := r.DB.Model(&model.Task{}).Where("id = ?", id).First(&task).Error
	if err != nil {
		return model.Task{}, err
	}
	return task, nil
}

func (r *Repository) ListTagCount(status int) ([]model.TagCount, error) {
	var tagCountList []model.TagCount
	sql := r.DB.Table("task").Select("tag", "count(*) as total")
	if status > 0 {
		sql = sql.Where("status = ?", status)
	}
	err := sql.Group("tag").Find(&tagCountList).Error
	if err != nil {
		return nil, err
	}
	return tagCountList, err
}

func (r *Repository) ListSpecCount(status int) ([]model.SpecCount, error) {
	var specCountList []model.SpecCount
	sql := r.DB.Table("task").Select("spec", "count(*) as number")
	if status > 0 {
		sql = sql.Where("status = ?", status)
	}
	err := sql.Group("spec").Find(&specCountList).Error
	if err != nil {
		return nil, err
	}
	return specCountList, err
}

func (r *Repository) ListStartedTaskBySpec(spec string) ([]model.Task, error) {
	var taskList []model.Task
	err := r.DB.Model(&model.Task{}).Where("spec = ? and status = ?", spec, 1).Find(&taskList).Error
	return taskList, err
}

func (r *Repository) ListStartedSpec() ([]string, error) {
	var specList []string
	specCountList, err := r.ListSpecCount(1)
	if err != nil {
		return nil, err
	}
	for _, v := range specCountList {
		specList = append(specList, v.Spec)
	}
	return specList, err
}

func (r *Repository) TryExecuteTask(task model.Task) (int64, error) {
	sql := r.DB.Table("task").Where("id = ? and total = ?", task.Id, task.Total).Updates(map[string]interface{}{
		"total":     gorm.Expr("total + ?", 1),
		"time_lock": time.Now().UnixMilli() + base.Config().Db.ExecutedLockTime,
	})
	if sql.Error != nil {
		return 0, sql.Error
	}
	return sql.RowsAffected, nil
}

func (r *Repository) AddTask(task model.Task) (int64, error) {
	task.UpdatedAt = time.Now()
	task.CreatedAt = task.UpdatedAt
	return task.Id, r.DB.Table("task").Create(&task).Error
}

func (r *Repository) EditTask(task model.Task) error {
	task.UpdatedAt = time.Now()
	if task.BackupUrl == "nil" {
		err := r.DB.Table("task").Where("id = ?", task.Id).Updates(map[string]any{
			"backup_url": "",
			"updated_at": task.UpdatedAt,
		}).Error
		if err != nil {
			return err
		}
		task.BackupUrl = ""
	}
	return r.DB.Table("task").Where("id = ?", task.Id).Updates(&task).Error
}

func (r *Repository) DeleteTask(id int64) error {
	return r.DB.Table("task").Delete(model.Task{}, id).Error
}

func (r *Repository) AddRecord(record model.Record) error {
	err := r.DB.AutoMigrate(&model.Record{})
	if err != nil {
		return err
	}
	return r.DB.Model(&model.Record{}).Create(&record).Error
}

func (r *Repository) ListRecord(shard string, taskId int64, code int, startTime string, endTime string, pageIndex int, pageSize int) ([]model.Record, int64, error) {
	var recordList []model.Record
	var total int64
	sql := r.DB.Table(fmt.Sprintf("record_%s", shard)).Where("task_id = ?", taskId)
	if code != 0 {
		sql.Where("code = ?", code)
	}
	if len(startTime) > 0 {
		sql.Where("executed_at >= ?", startTime)
	}
	if len(endTime) > 0 {
		sql.Where("executed_at <= ?", endTime)
	}
	sql.Count(&total)
	sql = sql.Order("id desc").Limit(pageSize).Offset((pageIndex - 1) * pageSize)
	err := sql.Find(&recordList).Error
	return recordList, total, err
}

func (r *Repository) ChangeTaskEditVersion() error {
	return r.DB.Table("metadata").Where("id = ?", 1).UpdateColumn("task_edit_version", gorm.Expr("task_edit_version + ?", 1)).Error
}

func (r *Repository) GetTaskEditVersion() (int64, error) {
	var metadata model.Metadata
	err := r.DB.Model(&model.Metadata{}).Where("id = ?", 1).First(&metadata).Error
	if err != nil {
		return 0, err
	}
	return metadata.TaskEditVersion, nil
}
