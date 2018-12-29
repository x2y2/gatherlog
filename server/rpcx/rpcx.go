package rpcx

import (
    "flag"
    "github.com/smallnest/rpcx/server"
    "github.com/smallnest/rpcx/serverplugin"
    "gatherlog/service"
    "time"
    "github.com/wonderivan/logger"
    "github.com/rcrowley/go-metrics"
    // "gatherlog/server/common"
)

var (
    addr = flag.String("addr","192.168.80.100:8888","server1 address")
    zkAddr = flag.String("zkAddr","192.168.20.71:2181","zookeeper address")
    basePath = flag.String("base","/rpcx_file","prefix path")
)

func RpcxStart(){
    flag.Parse()

    s := server.NewServer()
    
    addRegistryPlugin(s)

    s.RegisterName("Log",new(service.Log),"")
    s.Serve("tcp",*addr)

}

func addRegistryPlugin(s *server.Server){

    r := &serverplugin.ZooKeeperRegisterPlugin{
        ServiceAddress: "tcp@" + *addr,
        ZooKeeperServers: []string{*zkAddr},
        BasePath: *basePath,
        Metrics: metrics.NewRegistry(),
        UpdateInterval: time.Minute,
    }
    err := r.Start()
    if err != nil{logger.Error(err)}
    s.Plugins.Add(r)
}