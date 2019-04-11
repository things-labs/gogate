package elinkctls

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/thinkgos/gogate/misc"
	"github.com/thinkgos/gogate/protocol/elinkch/ctrl"
	"github.com/thinkgos/gomo/elink"
	"github.com/thinkgos/utils"
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
	rsp := GwInfosRspPy{
		misc.BuildDate(),
		misc.Version(),
		utils.RunningTime(),
		"",
	}

	out, err := jsoniter.Marshal(rsp)
	if err != nil {
		this.ErrorResponse(elink.CodeErrSysException)
		return
	}

	this.WriteResponse(elink.CodeSuccess, out)
}
