package g

import (
	"google.golang.org/grpc"
	"github.com/wonderivan/logger"
	"gatherlog/proto"
	logs "gatherlog/agent/gather"
	"gatherlog/agent/common"
	"sync"
	"time"
	"context"
	"runtime"
)


func GRpcStart(host string)error{
	conf := common.Config{}
	config := conf.ParseConfig()

	conn,err := grpc.Dial(host,grpc.WithInsecure())
	if err != nil{
		logger.Error(err)
	}
	defer conn.Close()

	c := gather.NewGatherLogClient(conn)

	runtime.GOMAXPROCS(runtime.NumCPU())

	var wg sync.WaitGroup
	wg.Add(1)
	go func(){
		defer  wg.Done()
		for {
			for ch := range logs.SendLog(){
				if ch.Data == "" {
					continue
				}

				req := gather.RequestLog{FileName: ch.Path, Content: ch.Data}
				_,err := c.Get(context.Background(),&req)
				if err != nil{
					logger.Error("Connect failed")
				}
			}
			time.Sleep(time.Second * config.Interval)
		}
	}()
	wg.Wait()
	return nil
}