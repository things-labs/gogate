package elinkctls

import (
	"net"

	"github.com/thinkgos/gogate/misc"
	"github.com/thinkgos/gogate/protocol/elinkch/ctrl"
	"github.com/thinkgos/gomo/elink"
	"github.com/thinkgos/utils"

	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
)

type GwInfosRspPy struct {
	BuildDate   string `json:"buildDate"`
	Version     string `json:"version"`
	RunningTime string `json:"runningTime"`
	LocalIP     string `json:"localIP"`
}

type GatewayInfos struct {
	ctrl.Controller
}

func (this *GatewayInfos) Get() {
	rsp := &GwInfosRspPy{
		misc.BuildDate(),
		misc.Version(),
		utils.RunningTime(),
		GetOutboundIP(),
	}

	out, err := jsoniter.Marshal(rsp)
	if err != nil {
		this.ErrorResponse(elink.CodeErrSysException)
		return
	}
	err = this.WriteResponsePy(elink.CodeSuccess, out)
	if err != nil {
		logs.Error(err)
	}
}

// Get preferred outbound ip of this machine
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}
