package common

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"os"
	"github.com/wonderivan/logger"
	"strconv"
)

type Config struct {
	Logpath string
	Addr string
	Pid string

}

func (this *Config)ParseConfig() *Config{
	conf, _ := ioutil.ReadFile("/Users/wangpei/Documents/GitHub/go/src/gatherlog/server/conf/config.yaml")
	yaml.Unmarshal(conf,&this)
	return this
}

func GetPid(){
	c := Config{}
	config := c.ParseConfig()
	
	fd, err := os.OpenFile("/Users/wangpei/Documents/GitHub/go/src/gatherlog/server/run/" + config.Pid,os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {logger.Error(err)}
	defer fd.Close()

	pid := os.Getpid()
	_, err = fd.Write([]byte(strconv.Itoa(pid)))
	fd.Sync()
}

func Log(){
	logger.SetLogger("gatherlog/server/conf/log.json")
}