package rpc

import (
	"gatherlog/agent/ip"
	"sync"
	"gatherlog/agent/gather"
	"gatherlog/agent/common"
	"net/rpc"
	"time"
	"github.com/wonderivan/logger"
	"net/rpc/jsonrpc"
	"math"
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

func (this RpcServer)RpcConn() (*rpc.Client,error){
	var retry int = 1
	for {
		conn, err := jsonrpc.Dial(this.Protocol,this.Host)
		if err != nil{
			logger.Error("dailing err %s",err)
			if retry > 3 {
				return nil,err
			}
			time.Sleep(time.Duration(math.Pow(2.0, float64(retry))) * time.Second)
			retry++
			continue
		}
		logger.Info("connect to %s\n",this.Host)
		return conn,nil
	}
}


func (this RpcServer)RpcStart() error{
	IP := ip.GetIP()

	c := common.Config{}
	config := c.ParseConfig()
	conn, _ := this.RpcConn()
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go func(){
		defer  wg.Done()
		for {
			for ch := range logs.SendLog(){
				if ch.Data == "" {
					continue
				}
				req := LogRequest{ch.Path,ch.Data,IP}
				res := LogResponse{}
				err := conn.Call("Log.ReciveLog",req,&res)
				if err != nil{
					logger.Error("call fail")
					conn.Close()
				}
			}
			time.Sleep(time.Second * config.Interval)
		}
	}()
	wg.Wait()
	return nil
}

