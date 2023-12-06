package model

type PageDTO struct {
	Total int64 `json:"total"`
	List  any   `json:"list"`
}

type TaskCommand struct {
	Status     int32             `json:"status"`
	Name       string            `json:"name"`
	Tag        string            `json:"tag"`
	Cron       string            `json:"cron"`
	RetryMax   int32             `json:"retryMax"`
	RetryCycle int32             `json:"retryCycle"`
	Url        string            `json:"url"`
	BackupUrl  string            `json:"backupUrl"`
	Method     string            `json:"method"`
	Body       string            `json:"body"`
	Header     map[string]string `json:"header"`
}

type TaskDTO struct {
	Id         int64             `json:"id"`
	Status     int32             `json:"status"`
	Name       string            `json:"name"`
	Tag        string            `json:"tag"`
	Cron       string            `json:"cron"`
	RetryMax   int32             `json:"retryMax"`
	RetryCycle int32             `json:"retryCycle"`
	Url        string            `json:"url"`
	BackupUrl  string            `json:"backupUrl"`
	Method     string            `json:"method"`
	Body       string            `json:"body"`
	Header     map[string]string `json:"header"`
	Total      int64             `json:"total"`
	CreatedAt  int64             `json:"createdAt"`
	UpdatedAt  int64             `json:"updatedAt"`
}

type RecordDTO struct {
	Id         int64  `json:"id"`
	TaskId     int64  `json:"taskId"`
	ExecutedAt int64  `json:"executedAt"`
	Result     string `json:"result"`
	TimeCost   int32  `json:"timeCost"`
	Code       int32  `json:"code"`
	IsBackup   int32  `json:"isBackup"`
	RetryCount int32  `json:"retryCount"`
}

type CreatedDTO struct {
	Id int64 `json:"id"`
}

type ErrorDTO struct {
	Type int    `json:"type"`
	Msg  string `json:"msg"`
}
