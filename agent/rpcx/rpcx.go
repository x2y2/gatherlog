package rpcx

import (
    "golang.org/x/net/context"
    "flag"
    "github.com/smallnest/rpcx/client"
    // "github.com/smallnest/rpcx/protocol"
    "github.com/wonderivan/logger"
    "gatherlog/service"
    "gatherlog/agent/gather"
    "gatherlog/agent/common"
    "sync"
    "time"
)

var (
    zkAddr = flag.String("zkAddr","192.168.20.71:2181","zookeeper address")
    basePath = flag.String("base","/rpcx_file","prefix path")
)

func RpcxClient(){
    c := common.Config{}
    config := c.ParseConfig()

    flag.Parse()

    // xch := make(chan *protocol.Message)

    d := client.NewZookeeperDiscovery(*basePath,"Log",[]string{*zkAddr},nil)
    xclient := client.NewXClient("Log",client.Failtry,client.RandomSelect,d,client.DefaultOption)
    defer xclient.Close()

    var wg sync.WaitGroup
    for {
        wg.Add(1)
        go callServer(xclient,&wg)
        time.Sleep(time.Second * config.Interval)
    }
}



func callServer(xclient client.XClient,wg *sync.WaitGroup){
    defer wg.Done()

    for ch := range gather.SendLog(){
        if ch.Data == "" {
            continue
        }
        req := &service.LogRequest{FileName: ch.Path,Content: ch.Data}
        reply := &service.LogResponse{}
        err := xclient.Call(context.Background(),"ReciveLog",req,reply)
        if err != nil{
            logger.Error("call fail")
        }
    }
}
