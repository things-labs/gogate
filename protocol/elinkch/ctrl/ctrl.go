// ctrl通道的实现
package ctrl

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/thinkgos/gogate/apps/broad"

	"github.com/thinkgos/gogate/models"
	"github.com/thinkgos/gomo/elink"

	jsoniter "github.com/json-iterator/go"
)

const (
	Signature_salt0 = `@#$%`
	Signature_salt1 = `^&*()`
)

const (
	ChannelData      = "data"
	ChannelCtrl      = "ctrl"
	ChannelCtrlReply = "ctrl_reply"
)

type BaseRequest struct {
	Topic     string `json:"topic,omitempty"`
	Timestamp string `json:"timestamp"`
	Signature string `json:"signature"`
	PacketID  int    `json:"packetID"`
}

type BaseResponse struct {
	Topic      string `json:"topic,omitempty"`
	PacketID   int    `json:"packetID"`
	Code       int    `json:"code"`
	CodeDetail string `json:"codeDetail,omitempty"`
	Message    string `json:"message,omitempty"`
}

type BaseData struct {
	Topic string `json:"topic,omitempty"`
}

type BaseRawPayload struct {
	Payload jsoniter.RawMessage `json:"payload,omitempty"`
}

type Request struct {
	*BaseRequest
	*BaseRawPayload
}

type Response struct {
	*BaseResponse
	*BaseRawPayload
}
type Data struct {
	*BaseData
	*BaseRawPayload
}
type Controller struct {
	elink.Controller
}

func init() {
	elink.RegisterChannelSelector(ChannelCtrl)
}

func (this *Controller) Prepare() {
	if !models.HasUser(this.Input.Topic.UserID) {
		this.ErrorResponse(elink.CodeErrCommonUserNoAccess)
		this.StopRun()
	}
	utimes := jsoniter.Get(this.Input.Payload, "timestamp").ToString()
	Expsign := jsoniter.Get(this.Input.Payload, "signature").ToString()
	sign := GenerateSignature(this.Input.Topic.Mac, utimes)
	if !strings.EqualFold(sign, Expsign) { // 验证签名
		this.ErrorResponse(elink.CodeErrCommonAuthorizationSignatureVerificationFailed)
		this.StopRun()
	}
}

func (this *Controller) Get() {
	this.ErrorResponse(elink.CodeErrCommonResourceMethodNotImplemented)
}

func (this *Controller) Post() {
	this.ErrorResponse(elink.CodeErrCommonResourceMethodNotImplemented)
}

func (this *Controller) Put() {
	this.ErrorResponse(elink.CodeErrCommonResourceMethodNotImplemented)
}

func (this *Controller) Patch() {
	this.ErrorResponse(elink.CodeErrCommonResourceMethodNotImplemented)
}

func (this *Controller) Delete() {
	this.ErrorResponse(elink.CodeErrCommonResourceMethodNotImplemented)
}

// 不带Payload错误回复,code为CodeSuccess将不进行回复
func (this *Controller) ErrorResponse(code int) error {
	if code != elink.CodeSuccess {
		return this.WriteResponsePy(code, nil)
	}
	return nil
}

// 回复,只关注payload即可
func (this *Controller) WriteResponsePy(code int, payload []byte) error {
	tp := elink.FromatRspTopic(this.Input.Topic)
	brsp := &BaseResponse{
		Topic:    tp,
		PacketID: jsoniter.Get(this.Input.Payload, "packetID").ToInt(),
		Code:     code,
	}
	if code != elink.CodeSuccess {
		errMsg := elink.CodeErrorMessage(code)
		brsp.CodeDetail = errMsg.Detail
		brsp.Message = errMsg.Message
	}

	out, err := jsoniter.Marshal(
		Response{
			brsp,
			&BaseRawPayload{payload},
		})
	if err != nil {
		return err
	}

	return this.WriteResponse(tp, out)
}

// 推送数据,向对应的推送通道推送数据
func Publish(resourse, method, messageType string, payload []byte) error {
	tp := elink.FormatPshTopic(ChannelData, resourse, method, messageType)
	out, err := jsoniter.Marshal(
		Data{
			&BaseData{tp},
			&BaseRawPayload{payload},
		})
	if err != nil {
		return err
	}

	return broad.Publish(tp, out)
}

//  签名mac + `@#$%` + timeStamp + `^&*()`拼接后md5 ,加盐值加密验证
func GenerateSignature(mac, timestamp string) string {
	h := md5.New()
	io.WriteString(h, mac)
	io.WriteString(h, Signature_salt0)
	io.WriteString(h, timestamp)
	io.WriteString(h, Signature_salt1)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// 获取主题上的productID参数,resource.productID
func (this *Controller) AcquireParamPid() (int, error) {
	spid := this.Input.Param.Get("productID")
	if spid == "" { // never happen but deal,may be other used
		return 0, errors.New("resource productID invalid")
	}

	pid, err := strconv.Atoi(spid)
	if err != nil { //never happen but deal
		return 0, errors.New("resource productID invalid")
	}

	return int(pid), nil
}
