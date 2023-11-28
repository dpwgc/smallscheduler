package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
	"smallscheduler/base"
	"smallscheduler/storage"
	"strconv"
)

const (
	SuccessCode        = 200
	ServiceErrorCode   = 400
	ParameterErrorCode = 401
)

func NewController() (*Controller, error) {
	service, err := storage.NewService()
	return &Controller{
		service: service,
	}, err
}

type Controller struct {
	service *storage.Service
}

func (c *Controller) ListTask(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	values := r.URL.Query()

	name := values.Get("name")
	status, _ := strconv.Atoi(values.Get("status"))
	pageIndex, _ := strconv.Atoi(values.Get("pageIndex"))
	pageSize, _ := strconv.Atoi(values.Get("pageSize"))

	list, total, err := c.service.ListTask(name, status, pageIndex, pageSize)
	if err != nil {
		c.fail(w, ServiceErrorCode, err.Error())
		return
	}
	c.success(w, c.buildTaskPageDTO(list, total))
}

func (c *Controller) GetTask(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, _ := strconv.ParseInt(p.ByName("id"), 10, 64)
	task, err := c.service.GetTask(id)
	if err != nil {
		c.fail(w, ServiceErrorCode, err.Error())
		return
	}
	c.success(w, c.buildTaskDTO(task))
}

func (c *Controller) AddTask(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cmd := TaskCommand{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.fail(w, ParameterErrorCode, err.Error())
		return
	}
	err = json.Unmarshal(body, &cmd)
	if err != nil {
		c.fail(w, ParameterErrorCode, err.Error())
		return
	}
	err = c.service.SaveTask(c.buildTask(0, cmd))
	if err != nil {
		c.fail(w, ServiceErrorCode, err.Error())
		return
	}
	c.success(w, nil)
}

func (c *Controller) EditTask(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, _ := strconv.ParseInt(p.ByName("id"), 10, 64)
	cmd := TaskCommand{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.fail(w, ParameterErrorCode, err.Error())
		return
	}
	err = json.Unmarshal(body, &cmd)
	if err != nil {
		c.fail(w, ParameterErrorCode, err.Error())
		return
	}
	err = c.service.SaveTask(c.buildTask(id, cmd))
	if err != nil {
		c.fail(w, ServiceErrorCode, err.Error())
		return
	}
	c.success(w, nil)
}

func (c *Controller) DeleteTask(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, _ := strconv.ParseInt(p.ByName("id"), 10, 64)
	err := c.service.DeleteTask(id)
	if err != nil {
		c.fail(w, ServiceErrorCode, err.Error())
		return
	}
	c.success(w, nil)
}

func (c *Controller) ListRecord(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	values := r.URL.Query()

	taskId, _ := strconv.ParseInt(values.Get("taskId"), 10, 64)
	startTime := values.Get("startTime")
	endTime := values.Get("endTime")
	pageIndex, _ := strconv.Atoi(values.Get("pageIndex"))
	pageSize, _ := strconv.Atoi(values.Get("pageSize"))

	list, total, err := c.service.ListRecord(taskId, startTime, endTime, pageIndex, pageSize)
	if err != nil {
		c.fail(w, ServiceErrorCode, err.Error())
		return
	}
	c.success(w, c.buildRecordPageDTO(list, total))
}

func (c *Controller) success(w http.ResponseWriter, data any) {
	//响应成功
	result := ResultDTO{
		Code: SuccessCode,
		Data: data,
	}
	resultBytes, _ := json.Marshal(result)
	_, err := w.Write(resultBytes)
	if err != nil {
		log.Println(base.LogErrorTag, err)
		return
	}
}

func (c *Controller) fail(w http.ResponseWriter, code int16, msg string) {
	result := ResultDTO{
		Code: code,
		Msg:  msg,
	}
	resultBytes, _ := json.Marshal(result)
	_, err := w.Write(resultBytes)
	if err != nil {
		log.Println(base.LogErrorTag, err)
		return
	}
}

func (c *Controller) buildTaskDTO(task storage.Task) TaskDTO {
	headerObj := map[string]string{}
	if len(task.Header) > 0 {
		err := json.Unmarshal([]byte(task.Header), &headerObj)
		if err != nil {
			headerObj = nil
		}
	}
	return TaskDTO{
		Id:        task.Id,
		Status:    task.Status,
		Name:      task.Name,
		Cron:      task.Cron,
		Url:       task.Url,
		Method:    task.Method,
		Body:      task.Body,
		Header:    headerObj,
		Total:     task.Total,
		CreatedAt: task.CreatedAt.UnixMilli(),
		UpdatedAt: task.UpdatedAt.UnixMilli(),
	}
}

func (c *Controller) buildTaskPageDTO(list []storage.Task, total int64) PageDTO {
	var dtoList []TaskDTO
	if len(list) > 0 {
		for _, v := range list {
			headerObj := map[string]string{}
			if len(v.Header) > 0 {
				err := json.Unmarshal([]byte(v.Header), &headerObj)
				if err != nil {
					headerObj = nil
				}
			}
			dtoList = append(dtoList, TaskDTO{
				Id:        v.Id,
				Status:    v.Status,
				Name:      v.Name,
				Cron:      v.Cron,
				Url:       v.Url,
				Method:    v.Method,
				Body:      v.Body,
				Header:    headerObj,
				Total:     v.Total,
				CreatedAt: v.CreatedAt.UnixMilli(),
				UpdatedAt: v.UpdatedAt.UnixMilli(),
			})
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
				ExecutedAt: v.ExecutedAt.UnixMilli(),
			})
		}
	}
	return PageDTO{
		Total: total,
		List:  dtoList,
	}
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
		Id:     id,
		Status: command.Status,
		Name:   command.Name,
		Cron:   command.Cron,
		Url:    command.Url,
		Method: command.Method,
		Body:   command.Body,
		Header: headerJson,
	}
}

type PageDTO struct {
	Total int64 `json:"total"`
	List  any   `json:"list"`
}

type ResultDTO struct {
	Code int16  `json:"code"`
	Data any    `json:"data,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

type TaskCommand struct {
	Status int32             `json:"status"`
	Name   string            `json:"name"`
	Cron   string            `json:"cron"`
	Url    string            `json:"url"`
	Method string            `json:"method"`
	Body   string            `json:"body"`
	Header map[string]string `json:"header"`
}

type TaskDTO struct {
	Id        int64             `json:"id"`
	Status    int32             `json:"status"`
	Name      string            `json:"name"`
	Cron      string            `json:"cron"`
	Url       string            `json:"url"`
	Method    string            `json:"method"`
	Body      string            `json:"body"`
	Header    map[string]string `json:"header"`
	Total     int64             `json:"total"`
	CreatedAt int64             `json:"createdAt"`
	UpdatedAt int64             `json:"updatedAt"`
}

type RecordDTO struct {
	Id         int64  `json:"id"`
	TaskId     int64  `json:"taskId"`
	ExecutedAt int64  `json:"executedAt"`
	Result     string `json:"result"`
	Code       int32  `json:"code"`
}
