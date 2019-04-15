package elinkctls

import (
	"github.com/thinkgos/gogate/models"
	"github.com/thinkgos/gogate/protocol/elinkch/ctrl"
	"github.com/thinkgos/gomo/elink"

	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"
)

type SysMultiUserPy struct {
	Uid []int64 `json:"uid"`
}

type SysUserPy struct {
	Uid int64 `json:"uid"`
}

type SysMultiUserRequest struct {
	ctrl.BaseRequest
	Payload SysMultiUserPy `json:"payload,omitempty"`
}

type SysUserController struct {
	ctrl.Controller
}

func (this *SysUserController) Get() {
	out, err := jsoniter.Marshal(&SysMultiUserPy{Uid: models.GetUsers()})
	if err != nil {
		this.ErrorResponse(elink.CodeErrSysOperationFailed)
		return
	}
	err = this.WriteResponse(elink.CodeSuccess, out)
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
		uid = req.Payload.Uid
	default:
		code = elink.CodeErrSysInvalidParameter
		return
	}

	sucUid := make([]int64, 0, len(uid))
	for _, v := range uid {
		if isDel {
			err = models.DeleteUser(v)
		} else {
			err = models.AddUser(v)
		}
		if err != nil {
			continue
		}
		sucUid = append(sucUid, v)
	}
	if len(sucUid) == 0 {
		code = elink.CodeErrSysOperationFailed
		return
	}

	var rspPy interface{}

	if isArray {
		rspPy = &SysMultiUserPy{Uid: sucUid}
	} else {
		rspPy = &SysUserPy{Uid: sucUid[0]}
	}
	out, err := jsoniter.Marshal(rspPy)
	if err != nil {
		code = elink.CodeErrSysOperationFailed
		return
	}

	err = this.WriteResponse(elink.CodeSuccess, out)
	if err != nil {
		logs.Error(err)
	}
}
