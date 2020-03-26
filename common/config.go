package common

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var Config *config

type server struct {
	Url string `yaml:"url" json:"-"`
}

type supervisor struct {
	Name string `yaml:"name" json:"name"`
	Url  string `yaml:"url" json:"-"`
}

type config struct {
	Server         *server       `yaml:"server" json:"server"`
	SupervisorList []*supervisor `yaml:"supervisorList" json:"supervisorList"`
}

func LoadConfig(configPath string) (err error) {
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("config file read failed: %s", err)
		os.Exit(-1)
	}
	err = yaml.Unmarshal(configData, &Config)
	if err != nil {
		fmt.Printf("config parse failed: %s", err)
		os.Exit(-1)
	}
	return nil
}
