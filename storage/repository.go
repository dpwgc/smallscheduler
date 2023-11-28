package storage

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"smallscheduler/base"
	"time"
)

func NewRepository() (*Repository, error) {
	db, err := loadDB()
	return &Repository{
		DB: db,
	}, err
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) ListTask(name string, status int, pageIndex int, pageSize int) ([]Task, int64, error) {
	var taskList []Task
	var total int64
	sql := r.DB.Model(&Task{})
	if status != 0 {
		sql = sql.Where("status = ?", status)
	}
	if len(name) > 0 {
		sql = sql.Where(fmt.Sprintf("name like %q", "%"+name+"%"))
	}
	sql.Count(&total)
	sql = sql.Order("id desc").Limit(pageSize).Offset((pageIndex - 1) * pageSize)
	err := sql.Find(&taskList).Error
	return taskList, total, err
}

func (r *Repository) GetTask(id int64) (Task, error) {
	var task Task
	err := r.DB.Model(&Task{}).Where("id = ?", id).First(&task).Error
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

func (r *Repository) ListStartedTaskByCron(cron string) ([]Task, error) {
	var taskList []Task
	err := r.DB.Model(&Task{}).Where("cron = ? and status = ?", cron, 1).Find(&taskList).Error
	return taskList, err
}

func (r *Repository) ListStartedCron() ([]string, error) {
	var taskList []Task
	var cronList []string
	err := r.DB.Model(&Task{}).Select("cron").Where("status = ?", 1).Group("cron").Find(&taskList).Error
	if err != nil {
		return nil, err
	}
	for _, task := range taskList {
		cronList = append(cronList, task.Cron)
	}
	return cronList, err
}

func (r *Repository) ExecuteTask(id int64) (int64, error) {
	task := Task{}
	err := r.DB.Model(&Task{}).Select("id", "total").Where("id = ?", id).First(&task).Error
	if err != nil {
		return 0, err
	}
	task.UpdatedAt = time.Now()
	sql := r.DB.Table("task").Where("id = ? and total = ?", task.Id, task.Total).UpdateColumn("total", gorm.Expr("total + ?", 1))
	if sql.Error != nil {
		return 0, sql.Error
	}
	return sql.RowsAffected, nil
}

func (r *Repository) SaveTask(task Task) error {
	task.UpdatedAt = time.Now()
	if task.Id > 0 {
		return r.DB.Table("task").Where("id = ?", task.Id).Updates(&task).Error
	} else {
		task.CreatedAt = task.UpdatedAt
		task.Status = 1
		return r.DB.Table("task").Create(&task).Error
	}
}

func (r *Repository) DeleteTask(id int64) error {
	if id <= 0 {
		return errors.New("id is abnormal")
	}
	return r.DB.Table("task").Delete(Task{}, id).Error
}

func (r *Repository) SaveRecord(record Record) error {
	record.ExecutedAt = time.Now()
	return r.DB.Table("record").Create(&record).Error
}

func (r *Repository) ListRecord(taskId int64, startTime string, endTime string, pageIndex int, pageSize int) ([]Record, int64, error) {
	var recordList []Record
	var total int64
	sql := r.DB.Model(&Record{}).Where("task_id = ?", taskId)
	if len(startTime) > 0 {
		sql = sql.Where("executed_at >= ?", startTime)
	}
	if len(endTime) > 0 {
		sql = sql.Where("executed_at <= ?", endTime)
	}
	sql.Count(&total)
	sql = sql.Order("id desc").Limit(pageSize).Offset((pageIndex - 1) * pageSize)
	err := sql.Find(&recordList).Error
	return recordList, total, err
}

func (r *Repository) AddTaskEditVersion() error {
	return r.DB.Table("metadata").Where("id = ?", 1).UpdateColumn("task_edit_version", gorm.Expr("task_edit_version + ?", 1)).Error
}

func (r *Repository) GetTaskEditVersion() (int64, error) {
	var metadata Metadata
	err := r.DB.Model(&Metadata{}).Where("id = ?", 1).First(&metadata).Error
	if err != nil {
		return 0, err
	}
	return metadata.TaskEditVersion, nil
}

func loadDB() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(base.Config().Db.Dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&Metadata{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&Task{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&Record{})
	if err != nil {
		return nil, err
	}
	err = loadMetadata(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func loadMetadata(db *gorm.DB) error {
	var total int64
	db.Model(&Metadata{}).Where("id = ?", 1).Count(&total)
	if total == 0 {
		return db.Table("metadata").Create(&Metadata{
			Id:              1,
			TaskEditVersion: 0,
		}).Error
	}
	return nil
}
