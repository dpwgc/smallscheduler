package api

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"os"
	"smallscheduler/base"
	"smallscheduler/core"
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

	name := strings.TrimSpace(values.Get("name"))
	tag := strings.TrimSpace(values.Get("tag"))
	cron := strings.TrimSpace(values.Get("cron"))
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
	c.ok(w, c.buildTaskPageDTO(list, total))
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
	if task.Id <= 0 {
		c.error(w, ServiceErrorType, "task is empty")
		return
	}
	c.ok(w, c.buildTaskDTO(task))
}

func (c *Controller) ListTag(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	values := r.URL.Query()
	status, _ := strconv.Atoi(values.Get("status"))

	list, err := c.service.ListTagCount(status)
	if err != nil {
		c.error(w, ServiceErrorType, err.Error())
		return
	}
	c.ok(w, list)
}

func (c *Controller) ListCron(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	values := r.URL.Query()
	status, _ := strconv.Atoi(values.Get("status"))

	list, err := c.service.ListCronCount(status)
	if err != nil {
		c.error(w, ServiceErrorType, err.Error())
		return
	}
	c.ok(w, list)
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

	cmd.Cron = strings.TrimSpace(cmd.Cron)
	cmd.Tag = strings.TrimSpace(cmd.Tag)
	cmd.Name = strings.TrimSpace(cmd.Name)
	cmd.Url = strings.TrimSpace(cmd.Url)

	if len(cmd.Tag) == 0 {
		cmd.Tag = "default"
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
	c.created(w, model.CreatedDTO{
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

	cmd.Cron = strings.TrimSpace(cmd.Cron)
	cmd.Tag = strings.TrimSpace(cmd.Tag)
	cmd.Name = strings.TrimSpace(cmd.Name)
	cmd.Url = strings.TrimSpace(cmd.Url)

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
	c.noContent(w)
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
	c.noContent(w)
}

func (c *Controller) ExecuteTask(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
	if task.Id <= 0 {
		c.error(w, ServiceErrorType, "task is empty")
		return
	}
	go func() {
		// 使用主url发起请求
		if core.Handle(task, task.Url, 0) {
			return
		}
		// 如果主url请求失败，且有备用url，使用备用url发起请求
		if len(task.BackupUrl) > 0 {
			core.Handle(task, task.BackupUrl, 1)
		}
	}()
	c.noContent(w)
}

func (c *Controller) ListRecord(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	values := r.URL.Query()

	taskId, err := strconv.ParseInt(values.Get("taskId"), 10, 64)
	if err != nil {
		c.error(w, QueryParamErrorType, err.Error())
		return
	}
	startTime := strings.TrimSpace(values.Get("startTime"))
	endTime := strings.TrimSpace(values.Get("endTime"))
	sharding := strings.TrimSpace(values.Get("sharding"))
	code, _ := strconv.Atoi(values.Get("code"))
	pageIndex, _ := strconv.Atoi(values.Get("pageIndex"))
	pageSize, _ := strconv.Atoi(values.Get("pageSize"))

	if len(sharding) < 7 {
		dateStr := time.Now().Format("2006-01-02")
		dateArr := strings.Split(dateStr, "-")
		sharding = fmt.Sprintf("%s_%s", dateArr[0], dateArr[1])
	}

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
	c.ok(w, c.buildRecordPageDTO(list, total))
}

func (c *Controller) Health(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if core.Shutdown {
		w.WriteHeader(ErrorCode)
	} else {
		w.WriteHeader(OkCode)
	}
	_, err := w.Write([]byte("1"))
	if err != nil {
		base.Logger.Error(err.Error())
	}
}

func (c *Controller) Shutdown(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	values := r.URL.Query()
	wait, _ := strconv.Atoi(values.Get("wait"))
	if !core.Shutdown && wait > 0 {
		if strings.Contains(r.Host, "localhost") || strings.Contains(r.URL.Host, "127.0.0.1") || strings.Contains(r.URL.Host, "0.0.0.0") {
			core.Shutdown = true
			base.Logger.Warn(fmt.Sprintf("shutdown after %v seconds", wait))
			go func() {
				time.Sleep(time.Duration(wait) * time.Second)
				os.Exit(1)
			}()
		} else {
			base.Logger.Warn("only local shutdown requests are accepted")
		}
	}
	w.WriteHeader(OkCode)
	_, err := w.Write([]byte("1"))
	if err != nil {
		base.Logger.Error(err.Error())
	}
}
