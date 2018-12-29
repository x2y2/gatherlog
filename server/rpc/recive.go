package rpc

import (
	"os"
	"log"
	"gatherlog/server/common"
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

func (this *Log) ReciveLog(req LogRequest,res *LogResponse) error{
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