package elinkctls

import (
	"net/http"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/thinkgos/gogate/apps/elinkch/ctrl"
	"github.com/thinkgos/gomo/elink"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/validation"
	update "github.com/inconshreveable/go-update"
	jsoniter "github.com/json-iterator/go"
)

// GwUpReqPy 负载
type GwUpReqPy struct {
	URL       string `json:"url"`
	Checksum  string `json:"checksum"`
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature"`
	IsPatcher bool   `json:"isPatcher"`
}

// GwUpRequest 升级请求
type GwUpRequest struct {
	ctrl.BaseRequest
	Payload GwUpReqPy `json:"payload"`
}

// GatewayUpgradeController 网关应用程序升级控制器
type GatewayUpgradeController struct {
	ctrl.Controller
}

var isUpgradeInProcess bool

// Post 更新程序
func (this *GatewayUpgradeController) Post() {
	code := elink.CodeSuccess
	defer func() {
		if code != elink.CodeErrSysInProcess {
			isUpgradeInProcess = false
		}
		this.ErrorResponse(code)
	}()

	if isUpgradeInProcess {
		code = elink.CodeErrSysInProcess
		return
	}
	isUpgradeInProcess = true

	req := &GwUpRequest{}
	if err := jsoniter.Unmarshal(this.Input.Payload, req); err != nil {
		code = elink.CodeErrSysInvalidParameter
		return
	}
	rpl := req.Payload
	// check request parameter valid
	valid := validation.Validation{}
	valid.Required(rpl.URL, "url")
	//	valid.Required(rpl.Checksum, "checksum")
	//	valid.Required(rpl.Signature, "signature")
	//	valid.Required(rpl.PublicKey, "publicKey")
	if valid.HasErrors() {
		code = elink.CodeErrSysInvalidParameter
		return
	}

	if err := doUpdate(&rpl); err != nil {
		code = elink.CodeErrSysOperationFailed
		return
	}
	if err := this.WriteResponsePyServerJSON(elink.CodeSuccess, nil); err != nil {
		logs.Error("response failed!", err)
	}
	time.Sleep(time.Second) // give enough time to send the message to client
	bin, err := os.Executable()
	if err != nil {
		code = elink.CodeErrSysException
		logs.Error("path: find failed!", err)
		return
	}
	_, file := filepath.Split(bin)
	err = syscall.Exec(bin, []string{file}, os.Environ())
	if err != nil {
		code = elink.CodeErrSysException
		logs.Error("exec failed!", err.Error())
		return
	}
}

func doUpdate(iop *GwUpReqPy) error {
	//	ck, err := hex.DecodeString(iop.Checksum)
	//	if err != nil {
	//		return err
	//	}

	//	sign, err := hex.DecodeString(iop.Signature)
	//	if err != nil {
	//		return err
	//	}
	// default crypto.sha256 and ECDSAVerifier
	opt := update.Options{
		//		Checksum:  ck,
		//		Signature: sign,
	}
	//	if err = opt.SetPublicKeyPEM([]byte(iop.PublicKey)); err != nil {
	//		return err
	//	}
	//	if iop.IsPatcher {
	//		opt.Patcher = update.NewBSDiffPatcher()
	//	}

	resp, err := http.Get(iop.URL) // get the new file
	if err != nil {
		logs.Debug("failed go get:　%s\n", err)
		return err
	}
	defer resp.Body.Close()
	err = update.Apply(resp.Body, opt)
	if err != nil {
		if rerr := update.RollbackError(err); rerr != nil {
			logs.Debug("Failed to rollback from bad update: %v", rerr)
		}
	}

	return err
}
