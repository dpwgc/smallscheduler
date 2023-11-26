package storage

import (
	"time"
)

type Metadata struct {
	Id              int64 `gorm:"column:id;not null;autoIncrement;primaryKey;"`
	TaskEditVersion int64 `gorm:"column:task_edit_version;not null;default:0;"`
}

type Task struct {
	Id        int64     `gorm:"column:id;not null;autoIncrement;primaryKey;"`
	Status    int32     `gorm:"column:status;not null;default:1;"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP;"`
	Name      string    `gorm:"column:name;index:uni_name,unique;not null;type:varchar(200);default:'';"`
	Cron      string    `gorm:"column:cron;index:idx_cron;not null;type:varchar(50);default:'';"`
	Url       string    `gorm:"column:url;not null;type:text;default:'';"`
	Method    string    `gorm:"column:method;not null;type:varchar(5);default:'';"`
	Body      string    `gorm:"column:body;not null;type:text;default:'';"`
	Header    string    `gorm:"column:header;not null;type:text;default:'';"`
	Total     int64     `gorm:"column:total;not null;default:0;"`
}

type Record struct {
	Id         int64     `gorm:"column:id;not null;autoIncrement;primaryKey;"`
	TaskId     int64     `gorm:"column:task_id;index:idx_task_id;not null;"`
	ExecutedAt time.Time `gorm:"column:executed_at;index:idx_executed_at;not null;default:CURRENT_TIMESTAMP;"`
	Result     string    `gorm:"column:result;not null;type:text;default:'';"`
	Status     int32     `gorm:"column:status;not null;default:1;"`
}
