package base

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

const ConfigFilePath = "./config.yaml"

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
	configBytes, err := os.ReadFile(ConfigFilePath)
	if err != nil {
		log.Println(LogErrorTag, err)
		panic(err)
	}
	err = yaml.Unmarshal(configBytes, &Config)
	if err != nil {
		log.Println(LogErrorTag, err)
		panic(err)
	}
}
