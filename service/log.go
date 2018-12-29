package service

import (
    "golang.org/x/net/context"
    "gatherlog/server/common"
    "github.com/wonderivan/logger"
    "github.com/smallnest/rpcx/server"
    "os"
    "net"
    "strings"
)

type Log struct {
    
}

type LogRequest struct {
    FileName string
    Content string
    IP []string
}

type LogResponse struct {
    Size int64
}

var (
    ClientConn net.Conn
)

func (this *Log) ReciveLog(ctx context.Context,req *LogRequest,reply *LogResponse) error{
    c := common.Config{}
    config := c.ParseConfig()

    ClientConn = ctx.Value(server.RemoteConnContextKey).(net.Conn)
    addr := ClientConn.RemoteAddr().String()
    IP := strings.Split(addr,":")[0]


    f, err := os.OpenFile(config.Logpath + "/" + IP + "-" + req.FileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {logger.Error(err)}
    defer f.Close()

    if req.Content == "" {
        return nil
    }
    _, err = f.Write([]byte(req.Content))
    if err != nil {logger.Error(err)}
    f.Sync()
    
    return nil
}

// func(this *Log) SendMessages(s *server.Server){
//     if ClientConn != nil{
//         err := s.SendMessage(ClientConn,"service_path","ReciveLog",nil,[]byte("call success"))
//         if err != nil{
//             logger.Error("failed to send message to %s:%v\n",ClientConn.RemoteAddr().String(),err)
//             ClientConn = nil
//         }
//     }
// }