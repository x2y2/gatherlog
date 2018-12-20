package chkpoint

import (
	"os"
	"gopkg.in/ini.v1"
	"strconv"
	"log"
	"gatherlog/agent/common"
)

func Chkpt(chkpt,logfile,key string) int64{
	//chkpoint 文件不存在，设置offsieze为0，表示从头开始读日志文件
	if !ChkptExist(chkpt){
		return int64(0)
	}
	//从chkpoint点读日志文件
	offset := GetChkpt(chkpt,logfile,key)
	return offset
}

func ChkptExist(chkpt string) bool{
	c := common.Config{}
	config := c.ParseConfig()
	_, err := os.Stat(config.Chkpoint)
	if err != nil {
		os.Create(config.Chkpoint)
		return false
	}
	return  true
}

func GetChkpt(chkpt,section,key string) int64{
	cfg,err := ini.Load(chkpt)
	if err != nil{log.Fatal(err)}

	offset := cfg.Section(section).Key(key).String()
	if offset != ""{
		offset,_ := strconv.ParseInt(offset,10,64)
		return offset
	}
	return int64(0)
}

func SetChkpt(chkpt,section,key,value string){
	cfg,err := ini.Load(chkpt)
	if err != nil{log.Fatal(err)}

	cfg.Section(section).Key(key).SetValue(value)
	cfg.SaveTo(chkpt)
}

