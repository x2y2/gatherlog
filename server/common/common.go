package common

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Logpath string
	Addr string

}

func (this *Config)ParseConfig() *Config{
	conf, _ := ioutil.ReadFile("gatherlog/server/conf/config.yaml")
	yaml.Unmarshal(conf,&this)
	return this
}
