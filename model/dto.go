package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type PageDTO struct {
	Total int64 `json:"total"`
	List  any   `json:"list"`
}

func NewPageDTO() *PageDTO {
	return &PageDTO{}
}

type TaskCommand struct {
	Status     int32             `json:"status"`
	Name       string            `json:"name"`
	Tag        string            `json:"tag"`
	Spec       string            `json:"spec"`
	RetryMax   int32             `json:"retryMax"`
	RetryCycle int32             `json:"retryCycle"`
	Url        string            `json:"url"`
	BackupUrl  string            `json:"backupUrl"`
	Method     string            `json:"method"`
	Body       string            `json:"body"`
	Header     map[string]string `json:"header"`
}

func (cmd *TaskCommand) ConversionAndVerifyWithAdd() error {
	cmd.Spec = strings.TrimSpace(cmd.Spec)
	cmd.Tag = strings.TrimSpace(cmd.Tag)
	cmd.Name = strings.TrimSpace(cmd.Name)
	cmd.Url = strings.TrimSpace(cmd.Url)
	if len(cmd.Tag) == 0 {
		cmd.Tag = "default"
	}
	if len(cmd.Name) == 0 {
		return errors.New("name is empty")
	}
	if len(cmd.Url) == 0 {
		return errors.New("url is empty")
	}
	if !verifyUrl(cmd.Url) {
		return errors.New("url format is incorrect")
	}
	if len(cmd.BackupUrl) > 0 && !verifyUrl(cmd.BackupUrl) {
		return errors.New("backup url format is incorrect")
	}
	if len(cmd.Spec) == 0 {
		return errors.New("spec is empty")
	}
	if cmd.Method != "GET" && cmd.Method != "POST" && cmd.Method != "PUT" && cmd.Method != "PATCH" && cmd.Method != "DELETE" {
		return errors.New("method is not match")
	}
	return nil
}

func (cmd *TaskCommand) ConversionAndVerifyWithEdit() error {
	cmd.Spec = strings.TrimSpace(cmd.Spec)
	cmd.Tag = strings.TrimSpace(cmd.Tag)
	cmd.Name = strings.TrimSpace(cmd.Name)
	cmd.Url = strings.TrimSpace(cmd.Url)
	if len(cmd.Url) > 0 && !verifyUrl(cmd.Url) {
		return errors.New("url format is incorrect")
	}
	if len(cmd.BackupUrl) > 0 && cmd.BackupUrl != "nil" && !verifyUrl(cmd.BackupUrl) {
		return errors.New("backup url format is incorrect")
	}
	if len(cmd.Method) > 0 && cmd.Method != "GET" && cmd.Method != "POST" && cmd.Method != "PUT" && cmd.Method != "PATCH" && cmd.Method != "DELETE" {
		return errors.New("method is not match")
	}
	return nil
}

type TaskQuery struct {
	Status    int    `json:"status"`
	Name      string `json:"name"`
	Tag       string `json:"tag"`
	Spec      string `json:"spec"`
	PageIndex int    `json:"pageIndex"`
	PageSize  int    `json:"pageSize"`
}

func (query *TaskQuery) ConversionAndVerify() error {
	query.Name = strings.TrimSpace(query.Name)
	query.Tag = strings.TrimSpace(query.Tag)
	query.Spec = strings.TrimSpace(query.Spec)
	return verifyPageQueryParams(query.PageIndex, query.PageSize)
}

type TaskDTO struct {
	Id         int64             `json:"id"`
	Status     int32             `json:"status"`
	Name       string            `json:"name"`
	Tag        string            `json:"tag"`
	Spec       string            `json:"spec"`
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

func NewTaskDTO() *TaskDTO {
	return &TaskDTO{}
}

func (dto *TaskDTO) Build(task Task) *TaskDTO {
	headerObj := map[string]string{}
	if len(task.Header) > 0 {
		_ = json.Unmarshal([]byte(task.Header), &headerObj)
	}
	dto.Id = task.Id
	dto.Status = task.Status
	dto.Name = task.Name
	dto.Tag = task.Tag
	dto.Spec = task.Spec
	dto.RetryMax = task.RetryMax
	dto.RetryCycle = task.RetryCycle
	dto.Url = task.Url
	dto.BackupUrl = task.BackupUrl
	dto.Method = task.Method
	dto.Body = task.Body
	dto.Header = headerObj
	dto.Total = task.Total
	dto.CreatedAt = task.CreatedAt.UnixMilli()
	dto.UpdatedAt = task.UpdatedAt.UnixMilli()
	return dto
}

func (dto *PageDTO) BuildWithTask(list []Task, total int64) *PageDTO {
	var dtoList []TaskDTO
	if len(list) > 0 {
		for _, v := range list {
			dtoItem := TaskDTO{}
			dtoList = append(dtoList, *dtoItem.Build(v))
		}
	}
	dto.Total = total
	dto.List = dtoList
	return dto
}

type RecordQuery struct {
	TaskId    int64  `json:"taskId"`
	Code      int    `json:"code"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Shard     string `json:"shard"`
	PageIndex int    `json:"pageIndex"`
	PageSize  int    `json:"pageSize"`
}

func (query *RecordQuery) ConversionAndVerify() error {
	query.StartTime = strings.TrimSpace(query.StartTime)
	query.EndTime = strings.TrimSpace(query.EndTime)
	query.Shard = strings.TrimSpace(query.Shard)
	if len(query.Shard) < 7 {
		dateStr := time.Now().Format("2006-01-02")
		dateArr := strings.Split(dateStr, "-")
		query.Shard = fmt.Sprintf("%s_%s", dateArr[0], dateArr[1])
	}
	return verifyPageQueryParams(query.PageIndex, query.PageSize)
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

func (dto *PageDTO) BuildWithRecord(list []Record, total int64) *PageDTO {
	var dtoList []RecordDTO
	if len(list) > 0 {
		for _, v := range list {
			dtoList = append(dtoList, RecordDTO{
				Id:         v.Id,
				TaskId:     v.TaskId,
				Result:     v.Result,
				Code:       v.Code,
				IsBackup:   v.IsBackup,
				TimeCost:   v.TimeCost,
				RetryCount: v.RetryCount,
				ExecutedAt: v.ExecutedAt.UnixMilli(),
			})
		}
	}
	dto.Total = total
	dto.List = dtoList
	return dto
}

type CreatedDTO struct {
	Id int64 `json:"id"`
}

type CommonDTO struct {
	Msg string `json:"msg"`
}

func verifyPageQueryParams(pageIndex int, pageSize int) error {
	if pageIndex <= 0 {
		return errors.New("page index must be greater than 0")
	}
	if pageSize < 0 {
		return errors.New("page size must be greater than or equal to 0")
	}
	return nil
}

func verifyUrl(sourceUrl string) bool {
	if len(sourceUrl) < 6 {
		return false
	}
	_, err := url.ParseRequestURI(sourceUrl)
	if err != nil {
		return false
	}
	u, err := url.Parse(sourceUrl)
	if err != nil || len(u.Scheme) == 0 || len(u.Host) == 0 {
		return false
	}
	// Check if the URL has a valid scheme (http or https)
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return true
}
