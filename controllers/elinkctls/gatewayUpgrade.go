package elinkctls

import (
	"net/http"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/thinkgos/gogate/protocol/elinkch/ctrl"
	"github.com/thinkgos/gomo/elink"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/validation"
	"github.com/inconshreveable/go-update"
	"github.com/json-iterator/go"
)

type GatewayUpgrade struct {
	ctrl.Controller
}

type GwUpReqPayload struct {
	Url       string `json:"url"`
	Checksum  string `json:"checksum"`
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature"`
	IsPatcher bool   `json:"isPatcher"`
}

type GwUpRequest struct {
	ctrl.BaseRequest
	Payload GwUpReqPayload `json:"payload,omitempty"`
}

func (this *GatewayUpgrade) Post() {
	code := elink.CodeSuccess
	defer func() {
		this.ErrorResponse(code)
	}()

	req := &GwUpRequest{}
	if err := jsoniter.Unmarshal(this.Input.Payload, req); err != nil {
		code = elink.CodeErrSysInvalidParameter
		return
	}
	rpl := req.Payload
	valid := validation.Validation{}
	valid.Required(rpl.Url, "url")
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
	this.WriteResponse(elink.CodeSuccess, nil)
	time.Sleep(3) // give enough time to send the message
	bin, err := os.Executable()
	if err != nil {
		code = elink.CodeErrSysException
		logs.Error("path: find failed!", err)
		return
	}
	_, file := filepath.Split(bin)
	if err = syscall.Exec(bin, []string{file}, os.Environ()); err != nil {
		code = elink.CodeErrSysException
		logs.Error("exec failed!%s", err.Error())
		return
	}
}

func doUpdate(iop *GwUpReqPayload) error {
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

	resp, err := http.Get(iop.Url) // get the new file
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
