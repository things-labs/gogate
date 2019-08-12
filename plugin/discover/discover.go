// UDP 本地发现插件,默认端口8091

package discover

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/thinkgos/gogate/apps/elinkch/ctrl"
	"github.com/thinkgos/gogate/misc"
	"github.com/thinkgos/memlog"

	jsoniter "github.com/json-iterator/go"
)

// 默认端口
const (
	DefaultDiscoverPort = 8091
)

// Run discover application.
// discover.Run("localhost:8091")
// discover.Run(":8091")
// discover.Run("127.0.0.1:8091")
func Run(params ...string) {
	var listenAddr string
	var listenPort int
	var err error

	if len(params) > 0 && params[0] != "" {
		strs := strings.Split(params[0], ":")
		if len(strs) > 0 && strs[0] != "" {
			listenAddr = strs[0]
		}
		if len(strs) > 1 && strs[1] != "" {
			// 转换错误使用默认的端口
			if listenPort, err = strconv.Atoi(strs[1]); err != nil {
				listenPort = DefaultDiscoverPort
			}
		}
	}
	listenAddr = fmt.Sprintf("%s:%d", listenAddr, listenPort)

	addr, err := net.ResolveUDPAddr("udp", listenAddr)
	if err != nil {
		memlog.Critical("discover: ", err)
		return
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		memlog.Critical("discover: ", err)
		return
	}
	defer conn.Close()
	memlog.Debug("discover: server Running on %s", listenAddr)
	for {
		if err := handleClient(conn); err != nil {
			memlog.Error("discover handle,", err)
		}
	}
}

// GatewayDiscoverReq 网关发现请求
type GatewayDiscoverReq struct {
	Topic      string `json:"topic"`
	ProductKey string `json:"productKey"`
	Mac        string `json:"mac"`
}

// GatewayDiscoverRsp 网关发现回复
type GatewayDiscoverRsp struct {
	Topic      string `json:"topic"`
	ProductKey string `json:"productKey"`
	Mac        string `json:"mac"`
	BuildDate  string `json:"buildDate"`
	Version    string `json:"version"`
}

func handleClient(conn *net.UDPConn) error {
	buf := make([]byte, 2048)
	m, remoteAddr, err := conn.ReadFromUDP(buf)
	if err != nil {
		return fmt.Errorf("read failed,%v", err)
	}
	rawData := buf[:m]
	req := &GatewayDiscoverReq{}
	if err = jsoniter.Unmarshal(rawData, req); err != nil {
		return err
	}

	if req.ProductKey != ctrl.TpInfos.ProductKey {
		return fmt.Errorf("productkey not match, %v", err)
	}

	out, err := jsoniter.Marshal(&GatewayDiscoverRsp{
		Topic:      req.Topic,
		ProductKey: req.ProductKey,
		Mac:        misc.Mac(),
		BuildDate:  misc.BuildDate(),
		Version:    misc.Version(),
	})
	if err != nil {
		memlog.Error(err)
		return err
	}

	_, err = conn.WriteToUDP(out, remoteAddr)
	return err
}
