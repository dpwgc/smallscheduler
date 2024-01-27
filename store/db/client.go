package db

import (
	"context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"smallscheduler/base"
	"smallscheduler/model"
	"time"
)

func newClient() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(base.Config().Db.Dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 &DBLogger{},
	})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&model.Metadata{}, &model.Task{})
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
	db.Model(&model.Metadata{}).Where("id = ?", 1).Count(&total)
	if total == 0 {
		return db.Table("metadata").Create(&model.Metadata{
			Id:              1,
			TaskEditVersion: 0,
		}).Error
	}
	return nil
}

type DBLogger struct{}

func (l *DBLogger) LogMode(logger.LogLevel) logger.Interface {
	return l
}

func (l *DBLogger) Info(ctx context.Context, text string, args ...interface{}) {
	base.Logger.Info(text, args)
}

func (l *DBLogger) Warn(ctx context.Context, text string, args ...interface{}) {
	base.Logger.Warn(text, args)
}

func (l *DBLogger) Error(ctx context.Context, text string, args ...interface{}) {
	base.Logger.Error(text, args)
}

func (l *DBLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {

}
