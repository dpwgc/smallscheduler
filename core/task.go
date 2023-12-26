package core

import (
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"io"
	"net/http"
	"smallscheduler/base"
	"smallscheduler/model"
	"strings"
	"time"
)

var Shutdown = false

// 调度任务
func scheduled(cronStr string) {
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
		go func(task model.Task) {
			if Shutdown {
				return
			}
			if time.Now().UnixMilli() <= task.TimeLock {
				return
			}
			yes, err := service.TryExecuteTask(task)
			if err != nil {
				base.Logger.Error(err.Error())
				return
			}
			if yes == 0 {
				return
			}
			// 使用主url发起请求
			if Execute(task, task.Url, 0) {
				return
			}
			// 如果主url请求失败，且有备用url，使用备用url发起请求
			if len(task.BackupUrl) > 0 {
				Execute(task, task.BackupUrl, 1)
			}
		}(task)
	}
}

// Execute 执行任务
func Execute(task model.Task, url string, isBackup int32) bool {
	for i := 0; i <= int(task.RetryMax); i++ {
		record := model.Record{
			TaskId:     task.Id,
			IsBackup:   isBackup,
			RetryCount: int32(i),
			ExecutedAt: time.Now(),
		}
		code, timeCost, result := httpSend(task.Method, url, task.Body, task.Header)
		record.Result = result
		record.Code = int32(code)
		record.TimeCost = int32(timeCost)
		err := service.AddRecord(record)
		if err != nil {
			base.Logger.Error(err.Error())
		}
		if record.Code >= 200 && record.Code < 300 {
			return true
		}
		if task.RetryCycle > 0 {
			time.Sleep(time.Duration(task.RetryCycle) * time.Millisecond)
		}
	}
	return false
}

// 发送http请求
func httpSend(method, url, body, header string) (int, int64, string) {
	if method != "POST" && method != "GET" && method != "PUT" && method != "PATCH" && method != "DELETE" {
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
