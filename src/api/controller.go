package api

import (
	"alisa-dispatch-center/src/base"
	"alisa-dispatch-center/src/storage"
	"alisa-dispatch-center/src/storage/rdb"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
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

	list, total, err := c.service.ListTaskToUser(name, status, pageIndex, pageSize)
	if err != nil {
		c.error(w, ServiceErrorCode, err.Error())
		return
	}
	c.success(w, c.buildTaskPageDTO(list, total))
}

func (c *Controller) SaveTask(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cmd := TaskCommand{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.error(w, ParameterErrorCode, err.Error())
		return
	}
	err = json.Unmarshal(body, &cmd)
	if err != nil {
		c.error(w, ParameterErrorCode, err.Error())
		return
	}
	err = c.service.SaveTask(c.buildTask(cmd))
	if err != nil {
		c.error(w, ServiceErrorCode, err.Error())
		return
	}
	c.success(w, nil)
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

func (c *Controller) error(w http.ResponseWriter, code int16, msg string) {
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

func (c *Controller) buildTaskPageDTO(list []rdb.Task, total int64) TaskPageDTO {
	var dtoList []TaskDTO
	if len(list) > 0 {
		for _, v := range list {
			dtoList = append(dtoList, TaskDTO{
				Id:        v.Id,
				Status:    v.Status,
				Name:      v.Name,
				Cron:      v.Cron,
				Url:       v.Url,
				Method:    v.Method,
				Body:      v.Body,
				Header:    v.Header,
				Total:     v.Total,
				CreatedAt: v.CreatedAt.UnixMilli(),
				UpdatedAt: v.UpdatedAt.UnixMilli(),
			})
		}
	}
	return TaskPageDTO{
		Total: total,
		List:  dtoList,
	}
}

func (c *Controller) buildTask(command TaskCommand) rdb.Task {
	return rdb.Task{
		Id:     command.Id,
		Status: command.Status,
		Name:   command.Name,
		Cron:   command.Cron,
		Url:    command.Url,
		Method: command.Method,
		Body:   command.Body,
		Header: command.Header,
	}
}

type TaskPageDTO struct {
	Total int64     `json:"total"`
	List  []TaskDTO `json:"list"`
}

type ResultDTO struct {
	Code int16  `json:"code"`
	Data any    `json:"data,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

type TaskCommand struct {
	Id     uint64 `json:"id"`
	Status int32  `json:"status"`
	Name   string `json:"name"`
	Cron   string `json:"cron"`
	Url    string `json:"url"`
	Method string `json:"method"`
	Body   string `json:"body"`
	Header string `json:"header"`
}

type TaskDTO struct {
	Id        uint64 `json:"id"`
	Status    int32  `json:"status"`
	Name      string `json:"name"`
	Cron      string `json:"cron"`
	Url       string `json:"url"`
	Method    string `json:"method"`
	Body      string `json:"body"`
	Header    string `json:"header"`
	Total     uint64 `json:"total"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
