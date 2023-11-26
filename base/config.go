package base

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

const ConfigFilePath = "./config.yaml"

var config ConfigModel

type ConfigModel struct {
	Server struct {
		Port        int    `yaml:"port"`
		ContextPath string `yaml:"context-path"`
	} `yaml:"server"`
	Db struct {
		Dsn string `yaml:"dsn"`
	} `yaml:"db"`
	Log struct {
		FileMaxAge int `yaml:"file-max-age"`
	} `yaml:"log"`
}

func Config() ConfigModel {
	return config
}

func InitConfig() {
	//加载客户端配置
	configBytes, err := os.ReadFile(ConfigFilePath)
	if err != nil {
		log.Println(LogErrorTag, err)
		panic(err)
	}
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		log.Println(LogErrorTag, err)
		panic(err)
	}
}
