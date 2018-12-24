package main

import (
	"gatherlog/agent/common"
	"gatherlog/agent/grpc"
	"github.com/wonderivan/logger"
)

func main(){
	c := common.Config{}
	config := c.ParseConfig()

	err := g.GRpcStart(config.Host)
	common.GetPid()

	if err != nil{
		logger.Error(err)
	}
	logger.Info("Connect to %s: ",config.Host)
}