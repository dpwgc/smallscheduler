package api

import (
	"encoding/json"
	"net/http"
	"smallscheduler/base"
	"smallscheduler/core"
	"smallscheduler/model"
	"strings"
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

func (c *Controller) success(w http.ResponseWriter, code int, obj any) {
	resultBytes := []byte("")
	if obj != nil {
		resultBytes, _ = json.Marshal(obj)
	}
	c.write(w, code, resultBytes)
}

func (c *Controller) error(w http.ResponseWriter, eType int, msg string) {
	resultBytes, _ := json.Marshal(model.ErrorDTO{
		Type: eType,
		Msg:  msg,
	})
	c.write(w, ErrorCode, resultBytes)
}

func (c *Controller) write(w http.ResponseWriter, code int, body []byte) {
	w.WriteHeader(code)
	_, err := w.Write(body)
	if err != nil {
		base.Logger.Error(err.Error())
		return
	}
}

func (c *Controller) buildTaskDTO(task model.Task) model.TaskDTO {
	headerObj := map[string]string{}
	if len(task.Header) > 0 {
		_ = json.Unmarshal([]byte(task.Header), &headerObj)
	}
	return model.TaskDTO{
		Id:         task.Id,
		Status:     task.Status,
		Name:       task.Name,
		Tag:        task.Tag,
		Cron:       task.Cron,
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

func (c *Controller) buildTaskPageDTO(list []model.Task, total int64) model.PageDTO {
	var dtoList []model.TaskDTO
	if len(list) > 0 {
		for _, v := range list {
			dtoList = append(dtoList, c.buildTaskDTO(v))
		}
	}
	return model.PageDTO{
		Total: total,
		List:  dtoList,
	}
}

func (c *Controller) buildRecordPageDTO(list []model.Record, total int64) model.PageDTO {
	var dtoList []model.RecordDTO
	if len(list) > 0 {
		for _, v := range list {
			dtoList = append(dtoList, model.RecordDTO{
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
	return model.PageDTO{
		Total: total,
		List:  dtoList,
	}
}

func (c *Controller) checkPageQueryParams(pageIndex int, pageSize int) string {
	if pageIndex <= 0 {
		return "page index must be greater than 0"
	}
	if pageSize < 0 {
		return "page size must be greater than or equal to 0"
	}
	return ""
}

func (c *Controller) checkAddTaskCommand(command model.TaskCommand) string {
	if len(command.Name) == 0 {
		return "name is empty"
	}
	if len(command.Url) == 0 {
		return "url is empty"
	}
	if !strings.HasPrefix(command.Url, "http") && !strings.HasPrefix(command.Url, "HTTP") && !strings.HasPrefix(command.Url, "Http") {
		return "url must start with http"
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
	if command.Method != "GET" && command.Method != "POST" && command.Method != "PUT" && command.Method != "PATCH" && command.Method != "DELETE" {
		return "method is not match"
	}
	return ""
}

func (c *Controller) checkEditTaskCommand(command model.TaskCommand) string {
	if len(command.Cron) > 0 {
		checkWorker := core.NewCronWorker()
		defer func() {
			checkWorker = nil
		}()
		_, err := checkWorker.AddFunc(command.Cron, func() {})
		if err != nil {
			return "cron spec error: " + err.Error()
		}
	}
	if len(command.Url) > 0 && !strings.HasPrefix(command.Url, "http") && !strings.HasPrefix(command.Url, "HTTP") && !strings.HasPrefix(command.Url, "Http") {
		return "url must start with http"
	}
	if len(command.Method) > 0 && command.Method != "GET" && command.Method != "POST" && command.Method != "PUT" && command.Method != "PATCH" && command.Method != "DELETE" {
		return "method is not match"
	}
	return ""
}

func (c *Controller) buildTask(id int64, command model.TaskCommand) model.Task {
	headerJson := ""
	if len(command.Header) > 0 {
		headerBytes, err := json.Marshal(command.Header)
		if err == nil {
			headerJson = string(headerBytes)
		}
	}
	return model.Task{
		Id:         id,
		Status:     command.Status,
		Name:       command.Name,
		Tag:        command.Tag,
		Cron:       command.Cron,
		RetryMax:   command.RetryMax,
		RetryCycle: command.RetryCycle,
		Url:        command.Url,
		Method:     command.Method,
		Body:       command.Body,
		Header:     headerJson,
	}
}
