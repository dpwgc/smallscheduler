package api

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"smallscheduler/model"
	"smallscheduler/storage"
	"strconv"
	"strings"
	"time"
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
	tag := values.Get("tag")
	cron := values.Get("cron")
	status, _ := strconv.Atoi(values.Get("status"))
	pageIndex, _ := strconv.Atoi(values.Get("pageIndex"))
	pageSize, _ := strconv.Atoi(values.Get("pageSize"))

	tip := c.checkPageQueryParams(pageIndex, pageSize)
	if len(tip) > 0 {
		c.error(w, QueryParamErrorType, tip)
		return
	}

	list, total, err := c.service.ListTask(name, tag, cron, status, pageIndex, pageSize)
	if err != nil {
		c.error(w, ServiceErrorType, err.Error())
		return
	}
	c.success(w, OkCode, c.buildTaskPageDTO(list, total))
}

func (c *Controller) GetTask(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil {
		c.error(w, PathParamErrorType, err.Error())
		return
	}
	task, err := c.service.GetTask(id)
	if err != nil {
		c.error(w, ServiceErrorType, err.Error())
		return
	}
	c.success(w, OkCode, c.buildTaskDTO(task))
}

func (c *Controller) ListTag(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	values := r.URL.Query()
	status, _ := strconv.Atoi(values.Get("status"))

	list, err := c.service.ListTagCount(status)
	if err != nil {
		c.error(w, ServiceErrorType, err.Error())
		return
	}
	c.success(w, OkCode, list)
}

func (c *Controller) ListCron(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	values := r.URL.Query()
	status, _ := strconv.Atoi(values.Get("status"))

	list, err := c.service.ListCronCount(status)
	if err != nil {
		c.error(w, ServiceErrorType, err.Error())
		return
	}
	c.success(w, OkCode, list)
}

func (c *Controller) AddTask(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cmd := model.TaskCommand{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.error(w, IOErrorType, err.Error())
		return
	}
	err = json.Unmarshal(body, &cmd)
	if err != nil {
		c.error(w, UnmarshalErrorType, err.Error())
		return
	}
	tip := c.checkAddTaskCommand(cmd)
	if len(tip) > 0 {
		c.error(w, CommandParamErrorType, tip)
		return
	}
	id, err := c.service.AddTask(c.buildTask(0, cmd))
	if err != nil {
		c.error(w, ServiceErrorType, err.Error())
		return
	}
	c.success(w, CreatedCode, model.CreatedDTO{
		Id: id,
	})
}

func (c *Controller) EditTask(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil {
		c.error(w, PathParamErrorType, err.Error())
		return
	}
	cmd := model.TaskCommand{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.error(w, IOErrorType, err.Error())
		return
	}
	err = json.Unmarshal(body, &cmd)
	if err != nil {
		c.error(w, UnmarshalErrorType, err.Error())
		return
	}
	tip := c.checkEditTaskCommand(cmd)
	if len(tip) > 0 {
		c.error(w, CommandParamErrorType, tip)
		return
	}
	err = c.service.EditTask(c.buildTask(id, cmd))
	if err != nil {
		c.error(w, ServiceErrorType, err.Error())
		return
	}
	c.success(w, NoContentCode, nil)
}

func (c *Controller) DeleteTask(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil {
		c.error(w, PathParamErrorType, err.Error())
		return
	}
	err = c.service.DeleteTask(id)
	if err != nil {
		c.error(w, ServiceErrorType, err.Error())
		return
	}
	c.success(w, NoContentCode, nil)
}

func (c *Controller) ListRecord(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	values := r.URL.Query()

	taskId, err := strconv.ParseInt(values.Get("taskId"), 10, 64)
	if err != nil {
		c.error(w, QueryParamErrorType, err.Error())
		return
	}
	startTime := values.Get("startTime")
	endTime := values.Get("endTime")
	code, _ := strconv.Atoi(values.Get("code"))
	sharding := values.Get("sharding")
	if len(sharding) < 7 {
		dateStr := time.Now().Format("2006-01-02")
		dateArr := strings.Split(dateStr, "-")
		sharding = fmt.Sprintf("%s_%s", dateArr[0], dateArr[1])
	}
	pageIndex, _ := strconv.Atoi(values.Get("pageIndex"))
	pageSize, _ := strconv.Atoi(values.Get("pageSize"))

	tip := c.checkPageQueryParams(pageIndex, pageSize)
	if len(tip) > 0 {
		c.error(w, QueryParamErrorType, tip)
		return
	}

	list, total, err := c.service.ListRecord(sharding, taskId, code, startTime, endTime, pageIndex, pageSize)
	if err != nil {
		c.error(w, ServiceErrorType, err.Error())
		return
	}
	c.success(w, OkCode, c.buildRecordPageDTO(list, total))
}
