package rdb

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"smallscheduler/base"
	"smallscheduler/model"
)

func newClient() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(base.Config().Db.Dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
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
