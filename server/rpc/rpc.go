
package rpc

import (
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
	"log"
	"gatherlog/server/common"
	"os"
	"net"
	"fmt"
)

type Log struct {}

type LogRequest struct {
	FileName string
	Content string
	IP []string
}

type LogResponse struct {
	Size int64
}

type RpcServer struct {
	Host string
	Protocol string
}

func (this *Log) ShowLog(req LogRequest,res *LogResponse) error{
	c := common.Config{}
	config := c.ParseConfig()

	f, err := os.OpenFile(config.Logpath + "/" + req.IP[0] + "-" + req.FileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {log.Fatal(err)}
	defer f.Close()

	if req.Content == "" {
		return nil
	}
	_, err = f.Write([]byte(req.Content))
	if err != nil {log.Fatal(err);return err}
	f.Sync()

	return nil
}

func ServerStart() {
	server := rpc.NewServer()
	server.Register(new(Log))

	c := common.Config{}
	config := c.ParseConfig()

	lis ,_ := net.Listen("tcp",config.Addr)
	fmt.Fprintf(os.Stdout,"%s","start connect\n")
	for {
		conn,err := lis.Accept()
		defer conn.Close()

		if err != nil{
			log.Fatal("rpc listener accept fail:",err)
			time.Sleep(time.Duration(2)*time.Millisecond)
			continue
		}
		fmt.Fprintf(os.Stdout,"%s","new connect is comming\n")
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}