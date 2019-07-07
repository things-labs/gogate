package elinkctls

import (
	"os/exec"

	"github.com/thinkgos/elink"
	"github.com/thinkgos/gogate/apps/elinkch/ctrl"

	jsoniter "github.com/json-iterator/go"
)

// GwCmd 网关命令
type GwCmd struct {
	Command string `json:"Command"`
}

// GwCmdRequest 网关命令请求
type GwCmdRequest struct {
	ctrl.BaseRequest
	Payload GwCmd `json:"payload"`
}

// GatewayCommands 网关命令控制器
type GatewayCommands struct {
	ctrl.Controller
}

// Post 接收网关控制命令
func (this *GatewayCommands) Post() {
	code := elink.CodeSuccess
	defer func() {
		this.ErrorResponse(code)
	}()

	req := &GwCmdRequest{}
	if err := jsoniter.Unmarshal(this.Input.Payload, req); err != nil {
		code = elink.CodeErrSysInvalidParameter
		return
	}

	switch req.Payload.Command {
	case "reboot":
		if err := exec.Command("reboot").Run(); err != nil {
			code = elink.CodeErrSysOperationFailed
			return
		}
	case "factoryReset":
	case "identify":
	default:
		code = elink.CodeErrSysNotSupport
		return
	}

	err := this.WriteResponsePyServerJSON(elink.CodeSuccess, nil)
	if err != nil {
		code = elink.CodeErrSysException
	}
}
