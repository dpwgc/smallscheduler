package base

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

var config ConfigModel

type ConfigModel struct {
	Server struct {
		Port             int    `yaml:"port"`
		ContextPath      string `yaml:"context-path"`
		ExecutedLockTime int64  `yaml:"executed-lock-time"`
	} `yaml:"server"`
	Db struct {
		Dsn string `yaml:"dsn"`
	} `yaml:"db"`
	Log struct {
		MaxSize    int `yaml:"max-size"`
		MaxAge     int `yaml:"max-age"`
		MaxBackups int `yaml:"max-backups"`
	} `yaml:"log"`
}

func Config() ConfigModel {
	return config
}

func InitConfig() {
	//加载客户端配置
	configBytes, err := os.ReadFile("./config.yaml")
	if err != nil {
		fmt.Println("config error:", err)
		time.Sleep(3 * time.Second)
		panic(err)
	}
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		fmt.Println("config error:", err)
		time.Sleep(3 * time.Second)
		panic(err)
	}
}
