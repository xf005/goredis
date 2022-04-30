package goredis

import (
	"sync"

	"github.com/xf005/logger"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Conf struct {
	Redis struct {
		Host string
		Pass string
		Db   int
	}
}

var (
	syncOnce sync.Once
	conf     *Conf
)

func NewConf() *Conf {
	syncOnce.Do(func() {
		logger.Info("redis init...")
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
