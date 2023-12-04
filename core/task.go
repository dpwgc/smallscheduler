package core

import (
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"io"
	"net/http"
	"smallscheduler/base"
	"smallscheduler/storage"
	"strings"
	"time"
)

const (
	Post   = "POST"
	Get    = "GET"
	Put    = "PUT"
	Patch  = "PATCH"
	Delete = "DELETE"
)

// 批量执行任务
func execute(cronStr string) {
	//获取该cron下的所有任务
	taskList, err := service.ListStartedTaskByCron(cronStr)
	if err != nil {
		base.Logger.Error(err.Error())
	}
	//如果任务列表长度为0，则删除该工作者
	if len(taskList) == 0 {
		worker, _ := workerFactory.Load(cronStr)
		worker.(*cron.Cron).Stop()
		workerFactory.Delete(cronStr)
		worker = nil
		base.Logger.Info("a invalid worker is deleted")
	}
	//循环请求
	for _, task := range taskList {
		go func(task storage.Task) {
			yes, err := service.TryExecuteTask(task)
			if err != nil {
				base.Logger.Error(err.Error())
				return
			}
			if yes == 0 {
				return
			}
			for i := 0; i <= int(task.RetryMax); i++ {
				record := storage.Record{
					TaskId:     task.Id,
					RetryCount: int32(i),
					ExecutedAt: time.Now(),
				}
				code, timeCost, result := request(task.Method, task.Url, task.Body, task.Header)
				record.Result = result
				record.Code = int32(code)
				record.TimeCost = int32(timeCost)
				err = service.AddRecord(record)
				if err != nil {
					base.Logger.Error(err.Error())
				}
				if record.Code >= 200 && record.Code < 300 {
					break
				}
				if task.RetryCycle > 0 {
					time.Sleep(time.Duration(task.RetryCycle) * time.Millisecond)
				}
			}
		}(task)
	}
}

func request(method, url, body, header string) (int, int64, string) {
	if method != Post && method != Get && method != Put && method != Patch && method != Delete {
		return -1, 0, "http method is not match"
	}
	if len(url) == 0 {
		return -1, 0, "http url is empty"
	}
	payload := strings.NewReader(body)
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return -1, 0, fmt.Sprintf("http build request error: %s", err.Error())
	}
	if len(header) > 2 {
		var headerMap map[string]string
		err = json.Unmarshal([]byte(header), &headerMap)
		if err != nil {
			return -1, 0, fmt.Sprintf("http header error: %s", err.Error())
		}
		if len(headerMap) > 0 {
			for k, v := range headerMap {
				req.Header.Add(k, v)
			}
		}
	}
	startTime := time.Now().UnixMilli()
	response, err := http.DefaultClient.Do(req)
	endTime := time.Now().UnixMilli()
	timeCost := endTime - startTime
	if err != nil {
		return -1, timeCost, fmt.Sprintf("http send error: %s", err.Error())
	}
	defer func(body io.ReadCloser) {
		err = body.Close()
		if err != nil {
			base.Logger.Error(err.Error())
		}
	}(response.Body)
	resultBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return -1, timeCost, fmt.Sprintf("http build response error: %s", err.Error())
	}
	return response.StatusCode, timeCost, string(resultBytes)
}
