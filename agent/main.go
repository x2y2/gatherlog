package main

import (
	"gatherlog/agent/common"
	"gatherlog/agent/rpc"
)

func main(){
	c := common.Config{}
	config := c.ParseConfig()

	rs := rpc.RpcServer{config.Host,"tcp"}
	rs.RpcStart()
}