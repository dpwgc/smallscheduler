package api

import (
	"errors"
	"fmt"
	"github.com/dpwgc/easierweb"
	"os"
	"smallscheduler/core"
	"smallscheduler/model"
	"smallscheduler/storage"
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

func (c *Controller) ListTask(ctx *easierweb.Context, query model.TaskQuery) (*model.PageDTO, error) {

	query.Name = strings.TrimSpace(query.Name)
	query.Tag = strings.TrimSpace(query.Tag)
	query.Cron = strings.TrimSpace(query.Cron)

	tip := c.checkPageQueryParams(query.PageIndex, query.PageSize)
	if len(tip) > 0 {
		return nil, errors.New(tip)
	}

	list, total, err := c.service.ListTask(query.Name, query.Tag, query.Cron, query.Status, query.PageIndex, query.PageSize)
	if err != nil {
		return nil, err
	}
	return c.buildTaskPageDTO(list, total), nil
}

func (c *Controller) GetTask(ctx *easierweb.Context) (*model.TaskDTO, error) {
	id := ctx.Path.GetInt64("id")
	task, err := c.service.GetTask(id)
	if err != nil {
		return nil, err
	}
	return c.buildTaskDTO(task), nil
}

func (c *Controller) ListTag(ctx *easierweb.Context) (*[]model.TagCount, error) {
	status := ctx.Query.GetInt("status")
	list, err := c.service.ListTagCount(status)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func (c *Controller) ListCron(ctx *easierweb.Context) (*[]model.CronCount, error) {
	status := ctx.Query.GetInt("status")
	list, err := c.service.ListCronCount(status)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func (c *Controller) AddTask(ctx *easierweb.Context, cmd model.TaskCommand) (*model.CreatedDTO, error) {

	cmd.Cron = strings.TrimSpace(cmd.Cron)
	cmd.Tag = strings.TrimSpace(cmd.Tag)
	cmd.Name = strings.TrimSpace(cmd.Name)
	cmd.Url = strings.TrimSpace(cmd.Url)

	if len(cmd.Tag) == 0 {
		cmd.Tag = "default"
	}

	tip := c.checkAddTaskCommand(cmd)
	if len(tip) > 0 {
		return nil, errors.New(tip)
	}

	id, err := c.service.AddTask(c.buildTask(0, cmd))
	if err != nil {
		return nil, err
	}

	return &model.CreatedDTO{
		Id: id,
	}, nil
}

func (c *Controller) EditTask(ctx *easierweb.Context, cmd model.TaskCommand) error {

	id := ctx.Path.GetInt64("id")

	cmd.Cron = strings.TrimSpace(cmd.Cron)
	cmd.Tag = strings.TrimSpace(cmd.Tag)
	cmd.Name = strings.TrimSpace(cmd.Name)
	cmd.Url = strings.TrimSpace(cmd.Url)

	tip := c.checkEditTaskCommand(cmd)
	if len(tip) > 0 {
		return errors.New(tip)
	}
	return c.service.EditTask(c.buildTask(id, cmd))
}

func (c *Controller) DeleteTask(ctx *easierweb.Context) error {

	id := ctx.Path.GetInt64("id")

	return c.service.DeleteTask(id)
}

func (c *Controller) ExecuteTask(ctx *easierweb.Context) error {

	id := ctx.Path.GetInt64("id")

	task, err := c.service.GetTask(id)
	if err != nil {
		return err
	}
	go func() {
		// 使用主url发起请求
		if core.Execute(task, task.Url, 0) {
			return
		}
		// 如果主url请求失败，且有备用url，使用备用url发起请求
		if len(task.BackupUrl) > 0 {
			core.Execute(task, task.BackupUrl, 1)
		}
	}()
	return nil
}

func (c *Controller) ListRecord(ctx *easierweb.Context, query model.RecordQuery) (*model.PageDTO, error) {

	query.StartTime = strings.TrimSpace(query.StartTime)
	query.EndTime = strings.TrimSpace(query.EndTime)
	query.Sharding = strings.TrimSpace(query.Sharding)

	if len(query.Sharding) < 7 {
		dateStr := time.Now().Format("2006-01-02")
		dateArr := strings.Split(dateStr, "-")
		query.Sharding = fmt.Sprintf("%s_%s", dateArr[0], dateArr[1])
	}

	tip := c.checkPageQueryParams(query.PageIndex, query.PageSize)
	if len(tip) > 0 {
		return nil, errors.New(tip)
	}

	list, total, err := c.service.ListRecord(query.Sharding, query.TaskId, query.Code, query.StartTime, query.EndTime, query.PageIndex, query.PageSize)
	if err != nil {
		return nil, err
	}
	return c.buildRecordPageDTO(list, total), nil
}

func (c *Controller) Health(ctx *easierweb.Context) (*model.CommonDTO, error) {
	if core.Shutdown {
		return nil, errors.New("closed")
	}
	return &model.CommonDTO{
		Ok:  true,
		Msg: "running",
	}, nil
}

func (c *Controller) Shutdown(ctx *easierweb.Context) (*model.CommonDTO, error) {
	wait := ctx.Query.GetInt("wait")
	if wait <= 0 {
		return nil, errors.New("wait must be greater than 0")
	}
	if core.Shutdown {
		return nil, errors.New("shutdown command has been triggered")
	}
	if strings.Contains(ctx.Host(), "localhost") || strings.Contains(ctx.Host(), "127.0.0.1") || strings.Contains(ctx.Host(), "0.0.0.0") {
		core.Shutdown = true
		go func() {
			time.Sleep(time.Duration(wait) * time.Second)
			os.Exit(1)
		}()
		return &model.CommonDTO{
			Ok:  true,
			Msg: fmt.Sprintf("shutdown after %v seconds", wait),
		}, nil
	}
	return nil, errors.New("only local shutdown requests are accepted")
}
