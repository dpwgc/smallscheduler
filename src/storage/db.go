package storage

import (
	"alisa-dispatch-center/src/common"
	"alisa-dispatch-center/src/constant"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(common.Config.Db.Dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	})
	if err != nil {
		common.Log.Println(constant.LogErrorTag, err.Error())
		panic(err)
	}
	err = db.AutoMigrate(&Task{})
	if err != nil {
		common.Log.Println(constant.LogErrorTag, err.Error())
		panic(err)
	}
	return db
}
