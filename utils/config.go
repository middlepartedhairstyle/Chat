package utils

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	App struct {
		AppName string `yaml:"appName"`
		AppMode string `yaml:"appMode"`
		AppHost string `yaml:"appHost"`
		AppPort string `yaml:"appPort"`
	} `yaml:"app"`
	Email struct {
		From     string `yaml:"from"`
		To       string `yaml:"to"`
		Port     int    `yaml:"port"`
		Host     string `yaml:"host"`
		Password string `yaml:"password"`
	} `yaml:"email"`
	MySQL struct {
		Driver   string `yaml:"driver"`
		User     string `yaml:"user"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
		DBName   string `yaml:"db_name"`
	} `yaml:"mysql"`
	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		Db       int    `yaml:"db"`
	}
}

var Cfg = Config{}

func ReadConfig(path string) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &Cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println(Cfg)
}
