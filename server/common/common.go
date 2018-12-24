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
	conf, _ := ioutil.ReadFile("gatherlog/server/conf/config.yaml")
	yaml.Unmarshal(conf,&this)
	return this
}

func GetPid(){
	c := Config{}
	config := c.ParseConfig()
	pwd ,_ := os.Getwd()
	fd, err := os.OpenFile(pwd + "/" + config.Pid,os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {logger.Error(err)}
	defer fd.Close()

	pid := os.Getpid()
	_, err = fd.Write([]byte(strconv.Itoa(pid)))
	fd.Sync()
}