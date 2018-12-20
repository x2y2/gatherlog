package ip

import (
	"net"
	"fmt"
	"os"
)

//获取本机IP
func GetIP() []string{
	var IP []string = make([]string,0)
	addr,err := net.InterfaceAddrs()
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	for _,address := range addr {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				IP = append(IP,ipnet.IP.String())
			}
		}
	}
	return IP
}