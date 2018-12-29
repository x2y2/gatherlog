package g

import (
	"net"
	"gatherlog/server/common"
	"github.com/wonderivan/logger"
	"google.golang.org/grpc"
	"gatherlog/proto"
)

func GRpcStart(){
	c := common.Config{}
	config := c.ParseConfig()

	listener,err := net.Listen("tcp",config.Addr)
	if err != nil {logger.Error(err)}

	common.GetPid()

	logger.Info("server start")
	grpcServer := grpc.NewServer()
	s := gather.Log{}
	gather.RegisterGatherLogServer(grpcServer,&s)
	grpcServer.Serve(listener)
}