package elinkctls

import (
	"github.com/thinkgos/elink"
	"github.com/thinkgos/gogate/apps/elinkch/ctrl"
	"github.com/thinkgos/gogate/models"

	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
)

// SysMultiUserPy 多用户负载
type SysMultiUserPy struct {
	UID []string `json:"uid"`
}

// SysUserPy 单用户负载
type SysUserPy struct {
	UID string `json:"uid"`
}

// SysMultiUserRequest 多用户请求
type SysMultiUserRequest struct {
	ctrl.BaseRequest
	Payload SysMultiUserPy `json:"payload,omitempty"`
}

// SysUserController 用户控制控制器
type SysUserController struct {
	ctrl.Controller
}

// Get 获取用户
func (this *SysUserController) Get() {
	err := this.WriteResponsePyServerJSON(elink.CodeSuccess,
		&SysMultiUserPy{UID: models.GetUsers()})
	if err != nil {
		this.ErrorResponse(elink.CodeErrSysException)
		logs.Error("response", err)
	}
}

// Post 增加用户
func (this *SysUserController) Post() {
	this.userDeal(false)
}

// Delete 删除用户
func (this *SysUserController) Delete() {
	this.userDeal(true)
}

func (this *SysUserController) userDeal(isDel bool) {
	var uid []string
	var isArray bool
	var err error

	code := elink.CodeSuccess
	defer func() {
		this.ErrorResponse(code)
	}()

	juid := jsoniter.Get(this.Input.Payload, "payload", "uid")
	switch juid.ValueType() {
	case jsoniter.StringValue:
		uid = append(uid, juid.ToString())
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

	sucUID := make([]string, 0, len(uid))
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
