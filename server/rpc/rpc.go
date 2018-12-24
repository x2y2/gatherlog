
package rpc

import (
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
	"gatherlog/server/common"
	"net"
	"github.com/wonderivan/logger"
	"strings"
)


func ServerStart() {
	server := rpc.NewServer()
	server.Register(new(Log))

	c := common.Config{}
	config := c.ParseConfig()

	listener ,_ := net.Listen("tcp",config.Addr)
	logger.Info("start")
	common.GetPid()
	for {
		conn,err := listener.Accept()
		if err != nil{
			logger.Error("rpc listener accept fail %s",err)
			time.Sleep(time.Duration(2)*time.Millisecond)
			continue
		}
		defer conn.Close()

		addr := conn.RemoteAddr().String()
		ip := strings.Split(addr,":")[0]

		logger.Info("new connect is comming: %s",ip)
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}