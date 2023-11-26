package rdb

import (
	"alisa-dispatch-center/src/base"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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

func (r *Repository) ListTaskToUser(name string, status int, pageIndex int, pageSize int) ([]Task, int64, error) {
	var taskList []Task
	var total int64
	var sql *gorm.DB
	if status != 0 {
		sql = r.DB.Model(&Task{}).Where("status = ?", status)
	}
	if len(name) > 0 {
		sql = sql.Where(fmt.Sprintf("name like %q", "%"+name+"%"))
	}
	sql.Count(&total)
	sql = sql.Order("updated_at desc").Limit(pageSize).Offset((pageIndex - 1) * pageSize)
	err := sql.Find(&taskList).Error
	return taskList, total, err
}

func (r *Repository) ListTaskToServer(cron string) ([]Task, error) {
	var taskList []Task
	err := r.DB.Model(&Task{}).Where("cron = ? and status = ?", cron, 1).Find(&taskList).Error
	return taskList, err
}

func (r *Repository) ListCron() ([]string, error) {
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

func (r *Repository) DoTask(task Task) (int64, error) {
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
		return r.DB.Table("task").Updates(&task).Error
	} else {
		task.CreatedAt = task.UpdatedAt
		task.Status = 1
		return r.DB.Table("task").Create(&task).Error
	}
}

func (r *Repository) SaveRecord(record Record) error {
	record.ExecutedAt = time.Now()
	return r.DB.Table("record").Create(&record).Error
}

func loadDB() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(base.Config.Db.Dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&Task{})
	err = db.AutoMigrate(&Record{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
