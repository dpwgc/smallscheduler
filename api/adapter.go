package api

import (
	"errors"
	"fmt"
	"github.com/dpwgc/easierweb"
	"os"
	"smallscheduler/job"
	"smallscheduler/model"
	"smallscheduler/store"
	"strings"
	"time"
)

func newAdapter() (*Adapter, error) {
	service, err := store.NewService()
	return &Adapter{
		service: service,
	}, err
}

type Adapter struct {
	service *store.Service
}

func (a *Adapter) ListTask(ctx *easierweb.Context, query model.TaskQuery) (*model.PageDTO, error) {
	err := query.ConversionAndVerify()
	if err != nil {
		return nil, err
	}
	list, total, err := a.service.ListTask(query.Name, query.Tag, query.Spec, query.Status, query.PageIndex, query.PageSize)
	if err != nil {
		return nil, err
	}
	return model.NewPageDTO().BuildWithTask(list, total), nil
}

func (a *Adapter) GetTask(ctx *easierweb.Context) (*model.TaskDTO, error) {
	task, err := a.service.GetTask(ctx.Path.GetInt64("id"))
	if err != nil {
		return nil, err
	}
	return model.NewTaskDTO().Build(task), nil
}

func (a *Adapter) ListTag(ctx *easierweb.Context) (*[]model.TagCount, error) {
	list, err := a.service.ListTagCount(ctx.Query.GetInt("status"))
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func (a *Adapter) ListSpec(ctx *easierweb.Context) (*[]model.SpecCount, error) {
	list, err := a.service.ListSpecCount(ctx.Query.GetInt("status"))
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func (a *Adapter) AddTask(ctx *easierweb.Context, cmd model.TaskCommand) (*model.CreatedDTO, error) {
	err := cmd.ConversionAndVerifyWithAdd()
	if err != nil {
		return nil, err
	}
	err = job.VerifySpec(cmd.Spec)
	if err != nil {
		return nil, err
	}
	id, err := a.service.AddTask(*model.NewTask().Build(0, cmd))
	if err != nil {
		return nil, err
	}
	return &model.CreatedDTO{Id: id}, nil
}

func (a *Adapter) EditTask(ctx *easierweb.Context, cmd model.TaskCommand) error {
	err := cmd.ConversionAndVerifyWithEdit()
	if err != nil {
		return err
	}
	if len(cmd.Spec) > 0 {
		err = job.VerifySpec(cmd.Spec)
		if err != nil {
			return err
		}
	}
	return a.service.EditTask(*model.NewTask().Build(ctx.Path.GetInt64("id"), cmd))
}

func (a *Adapter) DeleteTask(ctx *easierweb.Context) error {
	return a.service.DeleteTask(ctx.Path.GetInt64("id"))
}

func (a *Adapter) ExecuteTask(ctx *easierweb.Context) error {
	task, err := a.service.GetTask(ctx.Path.GetInt64("id"))
	if err != nil {
		return err
	}
	go func() {
		// 使用主url发起请求
		if job.Execute(task, task.Url, 0) {
			return
		}
		// 如果主url请求失败，且有备用url，使用备用url发起请求
		if len(task.BackupUrl) > 0 {
			job.Execute(task, task.BackupUrl, 1)
		}
	}()
	return nil
}

func (a *Adapter) ListRecord(ctx *easierweb.Context, query model.RecordQuery) (*model.PageDTO, error) {
	err := query.ConversionAndVerify()
	if err != nil {
		return nil, err
	}
	list, total, err := a.service.ListRecord(query.Shard, query.TaskId, query.Code, query.StartTime, query.EndTime, query.PageIndex, query.PageSize)
	if err != nil {
		return nil, err
	}
	return model.NewPageDTO().BuildWithRecord(list, total), nil
}

func (a *Adapter) Health(ctx *easierweb.Context) (*model.CommonDTO, error) {
	if job.Shutdown {
		return nil, errors.New("closed")
	}
	return &model.CommonDTO{Msg: "running"}, nil
}

func (a *Adapter) Shutdown(ctx *easierweb.Context) (*model.CommonDTO, error) {
	wait := ctx.Query.GetInt("wait")
	if wait <= 0 {
		return nil, errors.New("wait must be greater than 0")
	}
	if job.Shutdown {
		return nil, errors.New("shutdown request has been triggered")
	}
	if strings.Contains(ctx.Host(), "localhost") || strings.Contains(ctx.Host(), "127.0.0.1") || strings.Contains(ctx.Host(), "0.0.0.0") {
		job.Shutdown = true
		go func() {
			time.Sleep(time.Duration(wait) * time.Second)
			os.Exit(1)
		}()
		return &model.CommonDTO{Msg: fmt.Sprintf("shutdown after %v seconds", wait)}, nil
	}
	return nil, errors.New("only local shutdown request are accepted")
}
