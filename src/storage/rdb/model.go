package rdb

import (
	"time"
)

type Task struct {
	Id        uint64    `gorm:"column:id;not null;autoIncrement;primaryKey;"`
	Status    int32     `gorm:"column:status;not null;default:1;"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP;"`
	Name      string    `gorm:"column:name;index:uni_name,unique;not null;type:varchar(200);default:'';"`
	Cron      string    `gorm:"column:cron;index:idx_cron;not null;type:varchar(50);default:'';"`
	Url       string    `gorm:"column:url;not null;type:text;default:'';"`
	Method    string    `gorm:"column:method;not null;type:varchar(5);default:'';"`
	Body      string    `gorm:"column:body;not null;type:text;default:'';"`
	Header    string    `gorm:"column:header;not null;type:text;default:'';"`
	Total     uint64    `gorm:"column:total;not null;default:0;"`
}

type Record struct {
	Id         uint64    `gorm:"column:id;not null;autoIncrement;primaryKey;"`
	TaskId     uint64    `gorm:"column:task_id;index:idx_task_id;not null;"`
	ExecutedAt time.Time `gorm:"column:executed_at;not null;default:CURRENT_TIMESTAMP;"`
	Result     string    `gorm:"column:result;not null;type:text;default:'';"`
}
