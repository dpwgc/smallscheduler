package ui

import (
	"alisa-dispatch-center/src/common"
	"alisa-dispatch-center/src/constant"
	"alisa-dispatch-center/src/storage"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"strconv"
)

func initApi() Api {
	return &apiImpl{
		service: storage.InitService(),
	}
}

type Api interface {
	ListTask(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	SaveTask(w http.ResponseWriter, r *http.Request, p httprouter.Params)
}

type apiImpl struct {
	service storage.Service
}

func (a *apiImpl) ListTask(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	values := r.URL.Query()

	appIdStr := values.Get("appId")
	appId, err := strconv.Atoi(appIdStr)
	if err != nil {
		common.Log.Println(constant.LogErrorTag, err)
		result(w, constant.HttpRequestFailCode, "", nil)
		return
	}

	envStr := values.Get("env")
	env, err := strconv.Atoi(envStr)
	if err != nil {
		common.Log.Println(constant.LogErrorTag, err)
		result(w, constant.HttpRequestFailCode, "", nil)
		return
	}

	name := values.Get("name")
	list, err := a.service.ListTaskToUser(uint64(appId), uint8(env), name)
	if err != nil {
		common.Log.Println(constant.LogErrorTag, err)
		result(w, constant.HttpRequestFailCode, "", nil)
		return
	}
	result(w, constant.HttpRequestSuccessCode, "", list)
}

func (a *apiImpl) SaveTask(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	task := storage.Task{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		common.Log.Println(constant.LogErrorTag, err)
		result(w, constant.HttpRequestFailCode, fmt.Sprintf("parameter error: %s", err.Error()), nil)
		return
	}
	err = json.Unmarshal(body, &task)
	if err != nil {
		common.Log.Println(constant.LogErrorTag, err)
		result(w, constant.HttpRequestFailCode, fmt.Sprintf("parameter error: %s", err.Error()), nil)
		return
	}
	err = a.service.SaveTask(task)
	if err != nil {
		common.Log.Println(constant.LogErrorTag, err)
		result(w, constant.HttpRequestFailCode, fmt.Sprintf("parameter error: %s", err.Error()), nil)
		return
	}
	result(w, constant.HttpRequestSuccessCode, "", nil)
}
