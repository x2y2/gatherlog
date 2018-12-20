package common

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"time"
)

type Config struct {
	Logpath string
	Chkpoint string
	Host string
	Buffer int
	Interval time.Duration

}

func (this *Config)ParseConfig() *Config{
	conf, _ := ioutil.ReadFile("gatherlog/agent/conf/config.yaml")
	yaml.Unmarshal(conf,&this)
	return this
}
