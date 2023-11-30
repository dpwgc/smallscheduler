package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"smallscheduler/storage"
	"strconv"
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
		c.fail(w, err.Error())
		return
	}
	c.success(w, c.buildTaskPageDTO(list, total))
}

func (c *Controller) GetTask(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, _ := strconv.ParseInt(p.ByName("id"), 10, 64)
	task, err := c.service.GetTask(id)
	if err != nil {
		c.fail(w, err.Error())
		return
	}
	c.success(w, c.buildTaskDTO(task))
}

func (c *Controller) AddTask(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cmd := TaskCommand{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.fail(w, err.Error())
		return
	}
	err = json.Unmarshal(body, &cmd)
	if err != nil {
		c.fail(w, err.Error())
		return
	}
	err = c.service.SaveTask(c.buildTask(0, cmd))
	if err != nil {
		c.fail(w, err.Error())
		return
	}
	c.success(w, nil)
}

func (c *Controller) EditTask(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, _ := strconv.ParseInt(p.ByName("id"), 10, 64)
	cmd := TaskCommand{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.fail(w, err.Error())
		return
	}
	err = json.Unmarshal(body, &cmd)
	if err != nil {
		c.fail(w, err.Error())
		return
	}
	err = c.service.SaveTask(c.buildTask(id, cmd))
	if err != nil {
		c.fail(w, err.Error())
		return
	}
	c.success(w, nil)
}

func (c *Controller) DeleteTask(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, _ := strconv.ParseInt(p.ByName("id"), 10, 64)
	err := c.service.DeleteTask(id)
	if err != nil {
		c.fail(w, err.Error())
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
		c.fail(w, err.Error())
		return
	}
	c.success(w, c.buildRecordPageDTO(list, total))
}
