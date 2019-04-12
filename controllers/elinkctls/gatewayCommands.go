package elinkctls

import (
	"os/exec"

	"github.com/thinkgos/gogate/protocol/elinkch/ctrl"
	"github.com/thinkgos/gomo/elink"

	jsoniter "github.com/json-iterator/go"
)

type GatewayCommands struct {
	ctrl.Controller
}

type GwCmdReqPy struct {
	Command string `json:"url"`
}

type GwCmdRequest struct {
	ctrl.BaseRequest
	Payload GwCmdReqPy `json:"payload"`
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

	this.WriteResponse(elink.CodeSuccess, nil)
}
