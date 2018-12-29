package gather

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/peer"
	"strings"
	"github.com/wonderivan/logger"
	"net"
	"os"
	"gatherlog/server/common"
)



type Log struct {}


func (this *Log) Get(ctx context.Context, req *RequestLog) (*ResponseLog,error){
	c := common.Config{}
	config := c.ParseConfig()

	pr, ok := peer.FromContext(ctx)
	if !ok {logger.Error("get client ip failed")}

	if pr.Addr == net.Addr(nil){logger.Error("Client ip is nil")}

	addr := pr.Addr.String()
	ip := strings.Split(addr,":")[0]

	f, err:= os.OpenFile(config.Logpath + "/" + ip + "-" + req.FileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {logger.Error(err)}
	defer f.Close()

	logger.Info("connect from %s: ",ip)

	if req.Content == "" {
		return &ResponseLog{},nil
	}
	_, err = f.Write([]byte(req.Content))
	if err != nil {logger.Error(err)}
	f.Sync()

	return &ResponseLog{},nil
}