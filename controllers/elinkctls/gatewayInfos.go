package elinkctls

import (
	"net"

	"github.com/thinkgos/gogate/apps/elinkch/ctrl"
	"github.com/thinkgos/gogate/misc"
	"github.com/thinkgos/gogate/pkg/elink"
	"github.com/thinkgos/memlog"
)

// GwInfos 网关信息
type GwInfos struct {
	BuildDate   string `json:"buildDate"`
	Version     string `json:"version"`
	RunningTime int64  `json:"runningTime"`
	LocalIP     string `json:"localIP"`
}

// GatewayInfosController 网关信息控制器
type GatewayInfosController struct {
	ctrl.Controller
}

// Get 获取网关信息
func (this *GatewayInfosController) Get() {
	err := this.WriteResponsePyServerJSON(elink.CodeSuccess, &GwInfos{
		misc.BuildDate(),
		misc.Version(),
		misc.RunningTime(),
		GetOutboundIP(),
	})
	if err != nil {
		this.ErrorResponse(elink.CodeErrSysException)
		memlog.Error("GatewayInfo: ", err)
	}
}

// GetOutboundIP Get preferred outbound ip of this machine
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}
