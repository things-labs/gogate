package elinkctls

import (
	"github.com/thinkgos/gogate/models"
	"github.com/thinkgos/gogate/protocol/elinkch/ctrl"
	"github.com/thinkgos/gomo/elink"

	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
)

type SysMultiUserPy struct {
	UID []int64 `json:"uid"`
}

type SysUserPy struct {
	UID int64 `json:"uid"`
}

type SysMultiUserRequest struct {
	ctrl.BaseRequest
	Payload SysMultiUserPy `json:"payload,omitempty"`
}

type SysUserController struct {
	ctrl.Controller
}

func (this *SysUserController) Get() {
	err := this.WriteResponsePyServerJSON(elink.CodeSuccess,
		&SysMultiUserPy{UID: models.GetUsers()})
	if err != nil {
		logs.Error(err)
	}
}

func (this *SysUserController) Post() {
	this.userDeal(false)
}

func (this *SysUserController) Delete() {
	this.userDeal(true)
}

func (this *SysUserController) userDeal(isDel bool) {
	var uid []int64
	var isArray bool
	var err error

	code := elink.CodeSuccess
	defer func() {
		this.ErrorResponse(code)
	}()

	juid := jsoniter.Get(this.Input.Payload, "payload", "uid")
	switch juid.ValueType() {
	case jsoniter.NumberValue:
		uid = append(uid, juid.ToInt64())
	case jsoniter.ArrayValue:
		isArray = true
		req := &SysMultiUserRequest{}
		err = jsoniter.Unmarshal(this.Input.Payload, req)
		if err != nil {
			code = elink.CodeErrSysInvalidParameter
			return
		}
		uid = req.Payload.UID
	default:
		code = elink.CodeErrSysInvalidParameter
		return
	}

	sucUID := make([]int64, 0, len(uid))
	for _, v := range uid {
		if isDel {
			err = models.DeleteUser(v)
		} else {
			err = models.AddUser(v)
		}
		if err != nil {
			continue
		}
		sucUID = append(sucUID, v)
	}
	if len(sucUID) == 0 {
		code = elink.CodeErrSysOperationFailed
		return
	}

	var rspPy interface{}
	if isArray {
		rspPy = &SysMultiUserPy{UID: sucUID}
	} else {
		rspPy = &SysUserPy{UID: sucUID[0]}
	}

	err = this.WriteResponsePyServerJSON(elink.CodeSuccess, rspPy)
	if err != nil {
		code = elink.CodeErrSysException
	}
}
