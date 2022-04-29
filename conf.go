package conf

import (
	"io/ioutil"
	"sync"

	"github.com/xf005/logger"
	"gopkg.in/yaml.v3"
)

type Conf struct {
	Server Server `yaml:"server"`
	Redis  Redis  `yaml:"redis"`
}

type Server struct {
	Port string `yaml:"port"`
}

type Redis struct {
	Host string `yaml:"host"`
	User string
	Pass string
	Db   int
}

var (
	once sync.Once
	conf *Conf
)

func NewConf() *Conf {
	once.Do(func() {
		logger.Info("conf init...")
		file, err := ioutil.ReadFile("./conf.yml")
		if err != nil {
			logger.Error(err.Error())
		}
		var c Conf
		err = yaml.Unmarshal(file, &c)
		if err != nil {
			logger.Error(err.Error())
		}
		conf = &c
	})
	return conf
}
