package storage

import (
	"time"
)

type Base struct {
	Id        uint64    `gorm:"column:id;autoIncrement;primaryKey;comment:id"`
	Status    int8      `gorm:"column:status;comment:status = 1 or -1;default:1;"`
	CreatedAt time.Time `gorm:"column:created_at;comment:create time;default:CURRENT_TIMESTAMP;"`
	UpdatedAt time.Time `gorm:"column:updated_at;comment:update time;default:CURRENT_TIMESTAMP;"`
}

type Task struct {
	Base
	NamespaceId uint64 `gorm:"index:,;column:namespace_id;comment:namespace id;default:0;"`
	GroupId     uint64 `gorm:"index:,;column:group_id;comment:group id;default:0;"`
	AppId       uint64 `gorm:"index:,;column:app_id;comment:app id;default:0;"`
	Env         uint8  `gorm:"column:env;comment:app environment;default:0;"`
	Partition   uint16 `gorm:"column:partition;comment:task partition;default:0;"`
	Name        string `gorm:"column:name;type:varchar(127);comment:core name;default:'';"`
	Remark      string `gorm:"column:remark;type:text;comment:core remark;default:'';"`
	Cron        string `gorm:"index:,;column:cron;type:varchar(15);comment:cron;default:'';"`
	Url         string `gorm:"column:url;type:text;comment:request url;default:'';"`
	Method      string `gorm:"column:method;type:varchar(7);comment:request method;default:'';"`
	Body        string `gorm:"column:body;type:text;comment:request body;default:'';"`
	Header      string `gorm:"column:header;type:text;comment:request header;default:'';"`
}
