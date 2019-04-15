package elinkctls

import (
	"os/exec"

	"github.com/astaxie/beego/logs"

	"github.com/thinkgos/gogate/protocol/elinkch/ctrl"
	"github.com/thinkgos/gomo/elink"

	jsoniter "github.com/json-iterator/go"
)

type GwCmdReqPy struct {
	Command string `json:"url"`
}

type GwCmdRequest struct {
	ctrl.BaseRequest
	Payload GwCmdReqPy `json:"payload"`
}

type GatewayCommands struct {
	ctrl.Controller
}

func (this *GatewayCommands) Post() {
	req := &GwCmdRequest{}
	if err := jsoniter.Unmarshal(this.Input.Payload, req); err != nil {
		this.ErrorResponse(elink.CodeErrSysInvalidParameter)
		return
	}

	switch req.Payload.Command {
	case "reboot":
		if err := exec.Command("reboot").Run(); err != nil {
			this.ErrorResponse(elink.CodeErrSysOperationFailed)
			return
		}
	case "factoryReset":
	case "identify":
	}

	err := this.WriteResponse(elink.CodeSuccess, nil)
	if err != nil {
		logs.Error(err)
	}
}
