package model

import (
	"fmt"
	"strings"
	"time"
)

type Metadata struct {
	Id              int64 `gorm:"column:id;not null;type:bigint(20);autoIncrement;primaryKey;"`
	TaskEditVersion int64 `gorm:"column:task_edit_version;not null;type:bigint(20);default:0;"`
}

type Task struct {
	Id         int64     `gorm:"column:id;not null;type:bigint(20);autoIncrement;primaryKey;"`
	Status     int32     `gorm:"column:status;not null;type:int(11);default:2;"`
	CreatedAt  time.Time `gorm:"column:created_at;not null;type:datetime;default:CURRENT_TIMESTAMP;"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null;type:datetime;default:CURRENT_TIMESTAMP;"`
	Name       string    `gorm:"column:name;not null;type:varchar(80);default:'';"`
	Cron       string    `gorm:"column:cron;index:idx_cron;not null;type:varchar(40);default:'';"`
	Tag        string    `gorm:"column:tag;index:idx_tag;not null;type:varchar(40);default:'';"`
	RetryMax   int32     `gorm:"column:retry_max;not null;type:int(11);default:0;"`
	RetryCycle int32     `gorm:"column:retry_cycle;not null;type:int(11);default:0;"`
	Url        string    `gorm:"column:url;type:text;"`
	BackupUrl  string    `gorm:"column:backup_url;type:text;"`
	Method     string    `gorm:"column:method;not null;type:varchar(6);default:'GET';"`
	Body       string    `gorm:"column:body;type:text;"`
	Header     string    `gorm:"column:header;type:text;"`
	Total      int64     `gorm:"column:total;not null;type:bigint(20);default:0;"`
}

type CronCount struct {
	Cron   string `gorm:"column:cron;" json:"cron"`
	Number int64  `gorm:"column:number;" json:"number"`
}

type TagCount struct {
	Tag    string `gorm:"column:tag;" json:"tag"`
	Number int64  `gorm:"column:number;" json:"number"`
}

type Record struct {
	Id         int64     `gorm:"column:id;not null;type:bigint(20);autoIncrement;primaryKey;"`
	TaskId     int64     `gorm:"column:task_id;index:idx_task_id;not null;type:bigint(20);default:0;"`
	ExecutedAt time.Time `gorm:"column:executed_at;index:idx_executed_at;not null;type:datetime;default:CURRENT_TIMESTAMP;"`
	Result     string    `gorm:"column:result;type:text;"`
	Code       int32     `gorm:"column:code;not null;type:int(11);default:0;"`
	TimeCost   int32     `gorm:"column:time_cost;not null;type:int(11);default:0;"`
	RetryCount int32     `gorm:"column:retry_count;not null;type:int(11);default:0;"`
	IsBackup   int32     `gorm:"column:is_backup;not null;type:int(11);default:0;"`
}

func (Record) TableName() string {
	dateStr := time.Now().Format("2006-01-02")
	dateArr := strings.Split(dateStr, "-")
	return fmt.Sprintf("record_%s_%s", dateArr[0], dateArr[1])
}
