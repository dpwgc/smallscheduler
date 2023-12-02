package api

import (
	"encoding/json"
	"log"
	"net/http"
	"smallscheduler/base"
	"smallscheduler/core"
	"smallscheduler/storage"
)

const (
	OkCode                = 200
	CreatedCode           = 201
	NoContentCode         = 204
	ErrorCode             = 400
	IOErrorType           = 10
	UnmarshalErrorType    = 11
	PathParamErrorType    = 20
	QueryParamErrorType   = 21
	CommandParamErrorType = 22
	ServiceErrorType      = 30
)

type PageDTO struct {
	Total int64 `json:"total"`
	List  any   `json:"list"`
}

type TaskCommand struct {
	Status     int32             `json:"status"`
	Name       string            `json:"name"`
	Cron       string            `json:"cron"`
	Delay      int32             `json:"delay"`
	RetryMax   int32             `json:"retryMax"`
	RetryCycle int32             `json:"retryCycle"`
	Url        string            `json:"url"`
	Method     string            `json:"method"`
	Body       string            `json:"body"`
	Header     map[string]string `json:"header"`
}

type TaskDTO struct {
	Id         int64             `json:"id"`
	Status     int32             `json:"status"`
	Name       string            `json:"name"`
	Cron       string            `json:"cron"`
	Delay      int32             `json:"delay"`
	RetryMax   int32             `json:"retryMax"`
	RetryCycle int32             `json:"retryCycle"`
	Url        string            `json:"url"`
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
	RetryCount int32  `json:"retryCount"`
}

type CreatedDTO struct {
	Id int64 `json:"id"`
}

type ErrorDTO struct {
	Type int    `json:"type"`
	Msg  string `json:"msg"`
}

func (c *Controller) success(w http.ResponseWriter, code int, obj any) {
	resultBytes := []byte("")
	if obj != nil {
		resultBytes, _ = json.Marshal(obj)
	}
	c.write(w, code, resultBytes)
}

func (c *Controller) error(w http.ResponseWriter, eType int, msg string) {
	resultBytes, _ := json.Marshal(ErrorDTO{
		Type: eType,
		Msg:  msg,
	})
	c.write(w, ErrorCode, resultBytes)
}

func (c *Controller) write(w http.ResponseWriter, code int, body []byte) {
	w.WriteHeader(code)
	_, err := w.Write(body)
	if err != nil {
		log.Println(base.LogErrorTag, err)
		return
	}
}

func (c *Controller) buildTaskDTO(task storage.Task) TaskDTO {
	headerObj := map[string]string{}
	if len(task.Header) > 0 {
		_ = json.Unmarshal([]byte(task.Header), &headerObj)
	}
	return TaskDTO{
		Id:         task.Id,
		Status:     task.Status,
		Name:       task.Name,
		Cron:       task.Cron,
		Delay:      task.Delay,
		RetryMax:   task.RetryMax,
		RetryCycle: task.RetryCycle,
		Url:        task.Url,
		Method:     task.Method,
		Body:       task.Body,
		Header:     headerObj,
		Total:      task.Total,
		CreatedAt:  task.CreatedAt.UnixMilli(),
		UpdatedAt:  task.UpdatedAt.UnixMilli(),
	}
}

func (c *Controller) buildTaskPageDTO(list []storage.Task, total int64) PageDTO {
	var dtoList []TaskDTO
	if len(list) > 0 {
		for _, v := range list {
			dtoList = append(dtoList, c.buildTaskDTO(v))
		}
	}
	return PageDTO{
		Total: total,
		List:  dtoList,
	}
}

func (c *Controller) buildRecordPageDTO(list []storage.Record, total int64) PageDTO {
	var dtoList []RecordDTO
	if len(list) > 0 {
		for _, v := range list {
			dtoList = append(dtoList, RecordDTO{
				Id:         v.Id,
				TaskId:     v.TaskId,
				Result:     v.Result,
				Code:       v.Code,
				TimeCost:   v.TimeCost,
				RetryCount: v.RetryCount,
				ExecutedAt: v.ExecutedAt.UnixMilli(),
			})
		}
	}
	return PageDTO{
		Total: total,
		List:  dtoList,
	}
}

func (c *Controller) checkAddTaskCommand(command TaskCommand) string {
	if len(command.Name) == 0 {
		return "name is empty"
	}
	if len(command.Url) == 0 {
		return "url is empty"
	}
	if len(command.Cron) == 0 {
		return "cron is empty"
	}
	checkWorker := core.NewCronWorker()
	defer func() {
		checkWorker = nil
	}()
	_, err := checkWorker.AddFunc(command.Cron, func() {})
	if err != nil {
		return err.Error()
	}
	if command.Method != core.Get && command.Method != core.Post && command.Method != core.Put && command.Method != core.Patch && command.Method != core.Delete {
		return "method is not match"
	}
	return ""
}

func (c *Controller) checkEditTaskCommand(command TaskCommand) string {
	if len(command.Cron) > 0 {
		checkWorker := core.NewCronWorker()
		defer func() {
			checkWorker = nil
		}()
		_, err := checkWorker.AddFunc(command.Cron, func() {})
		if err != nil {
			return err.Error()
		}
	}
	if len(command.Method) > 0 && command.Method != core.Get && command.Method != core.Post && command.Method != core.Put && command.Method != core.Patch && command.Method != core.Delete {
		return "method is not match"
	}
	return ""
}

func (c *Controller) buildTask(id int64, command TaskCommand) storage.Task {
	headerJson := ""
	if len(command.Header) > 0 {
		headerBytes, err := json.Marshal(command.Header)
		if err == nil {
			headerJson = string(headerBytes)
		}
	}
	return storage.Task{
		Id:         id,
		Status:     command.Status,
		Name:       command.Name,
		Cron:       command.Cron,
		Delay:      command.Delay,
		RetryMax:   command.RetryMax,
		RetryCycle: command.RetryCycle,
		Url:        command.Url,
		Method:     command.Method,
		Body:       command.Body,
		Header:     headerJson,
	}
}
