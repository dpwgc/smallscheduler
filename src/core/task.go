package core

import (
	"alisa-dispatch-center/src/base"
	"alisa-dispatch-center/src/storage/rdb"
	"encoding/json"
	"errors"
	"github.com/robfig/cron/v3"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	Post = "POST"
	Get  = "GET"
)

// 批量执行任务
func execute(cronStr string) {
	//获取该cron的所有任务
	taskList, err := service.ListTaskToServer(cronStr)
	if err != nil {
		log.Println(base.LogErrorTag, err)
		return
	}
	//如果任务列表长度为0，则删除该工作者
	if len(taskList) == 0 {
		worker, _ := workerMap.Load(cronStr)
		worker.(*cron.Cron).Stop()
		workerMap.Delete(cronStr)
	}
	//循环请求
	for _, task := range taskList {
		go func(task rdb.Task) {
			i, err := service.DoTask(task)
			if err != nil {
				log.Println(base.LogErrorTag, err)
				return
			}
			if i == 0 {
				return
			}
			resp, err := request(task.Method, task.Url, task.Body, task.Header)
			if err != nil {
				log.Println(base.LogErrorTag, err)
			}
			err = service.SaveRecord(rdb.Record{
				TaskId: task.Id,
				Result: string(resp),
			})
			if err != nil {
				log.Println(base.LogErrorTag, err)
			}
		}(task)
	}
}

func request(method, url, body, header string) ([]byte, error) {
	if method != Post && method != Get {
		return nil, errors.New("method error")
	}
	payload := strings.NewReader(body)
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	if method == Post {
		req.Header.Add("Content-Type", "application/json")
	}
	if len(header) > 0 {
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
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println(base.LogErrorTag, err)
		}
	}(response.Body)
	return io.ReadAll(response.Body)
}
