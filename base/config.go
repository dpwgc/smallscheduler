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
		ConsoleEnable    bool   `yaml:"console-enable"`
		ExecutedLockTime int64  `yaml:"executed-lock-time"`
	} `yaml:"server"`
	Db struct {
		Dsn string `yaml:"dsn"`
	} `yaml:"db"`
	Log struct {
		Path       string `yaml:"path"`
		MaxSize    int    `yaml:"max-size"`
		MaxAge     int    `yaml:"max-age"`
		MaxBackups int    `yaml:"max-backups"`
	} `yaml:"log"`
}

func ServerYaml() string {
	marshal, err := yaml.Marshal(config.Server)
	if err != nil {
		return ""
	}
	return string(marshal)
}

func DbYaml() string {
	marshal, err := yaml.Marshal(config.Db)
	if err != nil {
		return ""
	}
	return string(marshal)
}

func LogYaml() string {
	marshal, err := yaml.Marshal(config.Log)
	if err != nil {
		return ""
	}
	return string(marshal)
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
