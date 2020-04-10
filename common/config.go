package common

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var Config *config

type server struct {
	Url  string `yaml:"url" json:"-"`
	Mode string `yaml:"mode" json:"-"`
}

type supervisor struct {
	Name string `yaml:"name" json:"name"`
	Url  string `yaml:"url" json:"-"`
}

type config struct {
	Server         *server       `yaml:"server" json:"-"`
	SupervisorList []*supervisor `yaml:"supervisorList" json:"supervisorList"`
}

func LoadConfig(configPath string) (err error) {
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalln(err)
	}
	err = yaml.Unmarshal(configData, &Config)
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}

func InitConfig() {
	configPathDefault, _ := os.Getwd()
	configPathDefault = path.Join(configPathDefault, "example.yml")
	configPath := flag.String("c", configPathDefault, "config file")
	flag.Parse()
	if err := LoadConfig(*configPath); err != nil {
		log.Fatalln(err)
	}
}
