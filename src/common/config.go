package common

import (
	"alisa-dispatch-center/src/constant"
	"gopkg.in/yaml.v3"
)

var Config configModel

type configModel struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Db struct {
		Dsn string `yaml:"dsn"`
	} `yaml:"db"`
	Log struct {
		FileMaxAge int `yaml:"file-max-age"`
	} `yaml:"log"`
}

func InitConfig() {
	//加载客户端配置
	configBytes, err := ReadFile(constant.ConfigFilePath)
	if err != nil {
		Log.Println(constant.LogErrorTag, err)
		panic(err)
	}
	err = yaml.Unmarshal(configBytes, &Config)
	if err != nil {
		Log.Println(constant.LogErrorTag, err)
		panic(err)
	}
}
