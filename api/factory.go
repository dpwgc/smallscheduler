package api

import (
	"encoding/json"
	"net/url"
	"smallscheduler/core"
	"smallscheduler/model"
)

func (c *Controller) buildTaskDTO(task model.Task) *model.TaskDTO {
	headerObj := map[string]string{}
	if len(task.Header) > 0 {
		_ = json.Unmarshal([]byte(task.Header), &headerObj)
	}
	return &model.TaskDTO{
		Id:         task.Id,
		Status:     task.Status,
		Name:       task.Name,
		Tag:        task.Tag,
		Spec:       task.Spec,
		RetryMax:   task.RetryMax,
		RetryCycle: task.RetryCycle,
		Url:        task.Url,
		BackupUrl:  task.BackupUrl,
		Method:     task.Method,
		Body:       task.Body,
		Header:     headerObj,
		Total:      task.Total,
		CreatedAt:  task.CreatedAt.UnixMilli(),
		UpdatedAt:  task.UpdatedAt.UnixMilli(),
	}
}

func (c *Controller) buildTaskPageDTO(list []model.Task, total int64) *model.PageDTO {
	var dtoList []model.TaskDTO
	if len(list) > 0 {
		for _, v := range list {
			dtoList = append(dtoList, *c.buildTaskDTO(v))
		}
	}
	return &model.PageDTO{
		Total: total,
		List:  dtoList,
	}
}

func (c *Controller) buildRecordPageDTO(list []model.Record, total int64) *model.PageDTO {
	var dtoList []model.RecordDTO
	if len(list) > 0 {
		for _, v := range list {
			dtoList = append(dtoList, model.RecordDTO{
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
	return &model.PageDTO{
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
	if len(command.Tag) == 0 {
		return "tag is empty"
	}
	if len(command.Name) == 0 {
		return "name is empty"
	}
	if len(command.Url) == 0 {
		return "url is empty"
	}
	if !isValidUrl(command.Url) {
		return "url format is incorrect"
	}
	if len(command.BackupUrl) > 0 && !isValidUrl(command.BackupUrl) {
		return "backup url format is incorrect"
	}
	if len(command.Spec) == 0 {
		return "spec is empty"
	}
	checkWorker := core.NewCronWorker()
	defer func() {
		checkWorker = nil
	}()
	_, err := checkWorker.AddFunc(command.Spec, func() {})
	if err != nil {
		return "spec error: " + err.Error()
	}
	if command.Method != "GET" && command.Method != "POST" && command.Method != "PUT" && command.Method != "PATCH" && command.Method != "DELETE" {
		return "method is not match"
	}
	return ""
}

func (c *Controller) checkEditTaskCommand(command model.TaskCommand) string {
	if len(command.Spec) > 0 {
		checkWorker := core.NewCronWorker()
		defer func() {
			checkWorker = nil
		}()
		_, err := checkWorker.AddFunc(command.Spec, func() {})
		if err != nil {
			return "spec error: " + err.Error()
		}
	}
	if len(command.Url) > 0 && !isValidUrl(command.Url) {
		return "url format is incorrect"
	}
	if len(command.BackupUrl) > 0 && command.BackupUrl != "nil" && !isValidUrl(command.BackupUrl) {
		return "backup url format is incorrect"
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
		Spec:       command.Spec,
		RetryMax:   command.RetryMax,
		RetryCycle: command.RetryCycle,
		Url:        command.Url,
		BackupUrl:  command.BackupUrl,
		Method:     command.Method,
		Body:       command.Body,
		Header:     headerJson,
	}
}

func isValidUrl(sourceUrl string) bool {
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
