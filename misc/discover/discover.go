package discover

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/thinkgos/gogate/misc"
	"github.com/thinkgos/gogate/protocol/elinkmd"
	"github.com/thinkgos/utils"

	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"
)

const (
	DefaultDiscoverPort = 8091
)

// Run discover application.
// discover.Run() default run on HttpPort
// discover.Run("localhost")
// discover.Run(":8091")
// discover.Run("127.0.0.1:8091")
func Run(params ...string) {
	var addrstr string

	listenAddr := "localhost"
	listenPort := int(0)

	if len(params) > 0 && params[0] != "" {
		strs := strings.Split(params[0], ":")
		if len(strs) > 0 && strs[0] != "" {
			listenAddr = strs[0]
		}
		if len(strs) > 1 && strs[1] != "" {
			listenPort, _ = strconv.Atoi(strs[1])
		}
	}
	if listenPort != 0 {
		addrstr = fmt.Sprintf("%s:%d", listenAddr, listenPort)
	} else {
		addrstr = fmt.Sprintf("%s:%d", listenAddr, DefaultDiscoverPort)
	}

	addr, err := net.ResolveUDPAddr("udp", addrstr)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	logs.Info("discover server Running on %s", addrstr)
	for {
		handleClient(conn)
	}
}

type GatewayDiscoverReq struct {
	Topic      string `json:"topic"`
	ProductKey string `json:"productKey"`
	Mac        string `json:"mac"`
}

type GatewayDiscoverRsp struct {
	Topic      string `json:"topic"`
	ProductKey string `json:"productKey"`
	Mac        string `json:"mac"`
	BuildDate  string `json:"buildDate"`
	Version    string `json:"version"`
}

func handleClient(conn *net.UDPConn) {
	buf := make([]byte, 2048)
	m, remoteAddr, err := conn.ReadFromUDP(buf)
	if err != nil {
		logs.Error("read failed!", err)
		return
	}
	rawData := buf[:m]
	req := &GatewayDiscoverReq{}
	err = jsoniter.Unmarshal(rawData, req)
	if err != nil {
		logs.Error(err)
		return
	}

	if req.ProductKey != elinkmd.ProductKey {
		logs.Error("productkey not match")
		return
	}

	rsp := GatewayDiscoverRsp{
		Topic:      req.Topic,
		ProductKey: req.ProductKey,
		Mac:        misc.Mac(),
		BuildDate:  utils.BuildDateTime(),
		Version:    misc.Version(),
	}

	out, err := jsoniter.Marshal(rsp)
	if err != nil {
		logs.Error(err)
		return
	}

	conn.WriteToUDP(out, remoteAddr)
}
