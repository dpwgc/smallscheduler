package core

import (
	"encoding/json"
	"errors"
	"github.com/robfig/cron/v3"
	"io"
	"log"
	"net/http"
	"smallscheduler/base"
	"smallscheduler/storage"
	"strings"
)

const (
	Post = "POST"
	Get  = "GET"
)

// 批量执行任务
func execute(cronStr string) {
	//获取该cron的所有任务
	taskCache, _ := taskCachePool.Load(cronStr)
	taskList := taskCache.([]storage.Task)
	//如果任务列表长度为0，则删除该工作者
	if len(taskList) == 0 {
		worker, _ := workerFactory.Load(cronStr)
		worker.(*cron.Cron).Stop()
		workerFactory.Delete(cronStr)
	}
	//循环请求
	for _, task := range taskList {
		go func(task storage.Task) {
			i, err := service.ExecuteTask(task)
			if err != nil {
				log.Println(base.LogErrorTag, err)
				return
			}
			if i == 0 {
				return
			}
			record := storage.Record{
				TaskId: task.Id,
			}
			response, err := request(task.Method, task.Url, task.Body, task.Header)
			if err != nil {
				record.Status = 2
				record.Result = err.Error()
				log.Println(base.LogErrorTag, err)
			} else {
				record.Status = 1
				record.Result = string(response)
			}
			err = service.SaveRecord(record)
			if err != nil {
				log.Println(base.LogErrorTag, err)
			}
		}(task)
	}
}

func request(method, url, body, header string) ([]byte, error) {
	if method != Post && method != Get {
		return nil, errors.New("method is not match")
	}
	payload := strings.NewReader(body)
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	if method == Post {
		req.Header.Add("Content-Type", "application/json")
	}
	if len(header) > 2 {
		var headerMap map[string]string
		err = json.Unmarshal([]byte(header), &headerMap)
		if err != nil {
			return nil, err
		}
		if len(headerMap) > 0 {
			for k, v := range headerMap {
				req.Header.Add(k, v)
			}
		}
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(body io.ReadCloser) {
		err = body.Close()
		if err != nil {
			log.Println(base.LogErrorTag, err)
		}
	}(response.Body)
	return io.ReadAll(response.Body)
}
