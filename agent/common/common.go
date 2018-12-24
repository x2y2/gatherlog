package common

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"time"
	"os"
	"github.com/wonderivan/logger"
	"strconv"
)

type Config struct {
	Logpath string
	Chkpoint string
	Pid string
	Host string
	Buffer int
	Interval time.Duration

}

func (this *Config)ParseConfig() *Config{
	pwd,_ := os.Getwd()
	configfile := pwd + "/gatherlog/agent/conf/config.yaml"
	conf, _ := ioutil.ReadFile(configfile)
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