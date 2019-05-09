// package ctrl 通道的实现

package ctrl

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/thinkgos/gogate/models"
	"github.com/thinkgos/gomo/elink"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

// 签名加盐值
const (
	SignatureSalt0 = `@#$%`
	SignatureSalt1 = `^&*()`
)

// 通道定义
const (
	ChannelData      = "data"
	ChannelCtrl      = "ctrl"
	ChannelCtrlReply = "ctrl_reply"
)

// BaseRequest 请求基本格式
type BaseRequest struct {
	Topic     string `json:"topic,omitempty"`
	Timestamp string `json:"timestamp"`
	Signature string `json:"signature"`
	PacketID  int    `json:"packetID"`
}

// BaseResponse 回复基本格式
type BaseResponse struct {
	Topic      string `json:"topic,omitempty"`
	PacketID   int    `json:"packetID"`
	Code       int    `json:"code"`
	CodeDetail string `json:"codeDetail,omitempty"`
	Message    string `json:"message,omitempty"`
}

// BasePublishData 推送基本格式
type BasePublishData struct {
	Topic string `json:"topic,omitempty"`
}

// BaseRawPayload Raw payload
type BaseRawPayload struct {
	Payload jsoniter.RawMessage `json:"payload,omitempty"`
}

// Request 请求含payload
type Request struct {
	*BaseRequest
	Payload interface{} `json:"payload,omitempty"`
}

// Response 回复含payload
type Response struct {
	*BaseResponse
	Payload interface{} `json:"payload,omitempty"`
}

// PublishData 推送含payload
type PublishData struct {
	*BasePublishData
	Payload interface{} `json:"payload,omitempty"`
}

// RawRequest 请求含raw payload
type RawRequest struct {
	*BaseRequest
	*BaseRawPayload
}

// RawResponse 回复含raw payload
type RawResponse struct {
	*BaseResponse
	*BaseRawPayload
}

// RawPublishData 推送含raw payload
type RawPublishData struct {
	*BasePublishData
	*BaseRawPayload
}

// Controller 控制器
type Controller struct {
	elink.Controller
}

func init() {
	elink.RegisterChannelSelector(ChannelCtrl)
}

// Prepare 前期准备
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

// Get 方法
func (this *Controller) Get() {
	this.ErrorResponse(elink.CodeErrCommonResourceMethodNotImplemented)
}

// Post 方法
func (this *Controller) Post() {
	this.ErrorResponse(elink.CodeErrCommonResourceMethodNotImplemented)
}

// Put 方法
func (this *Controller) Put() {
	this.ErrorResponse(elink.CodeErrCommonResourceMethodNotImplemented)
}

// Patch 方法
func (this *Controller) Patch() {
	this.ErrorResponse(elink.CodeErrCommonResourceMethodNotImplemented)
}

// Delete 方法
func (this *Controller) Delete() {
	this.ErrorResponse(elink.CodeErrCommonResourceMethodNotImplemented)
}

// ErrorResponse 不带Payload错误回复,code为CodeSuccess将不进行回复
func (this *Controller) ErrorResponse(code int) error {
	if code != elink.CodeSuccess {
		return this.WriteResponsePyServerJSON(code, nil)
	}
	return nil
}

// WriteResponsePyServerJSON 回复,只关注payload即可,json序列化由底层处理
func (this *Controller) WriteResponsePyServerJSON(code int, payload interface{}) error {
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

	out, err := jsoniter.Marshal(&Response{brsp, payload})
	if err != nil {
		return errors.Wrap(err, "json marshal failed")
	}

	return this.WriteResponse(tp, out)
}

// GenerateSignature 签名mac + `@#$%` + timeStamp + `^&*()`拼接后md5 ,加盐值加密验证
func GenerateSignature(mac, timestamp string) string {
	h := md5.New()
	io.WriteString(h, mac)
	io.WriteString(h, SignatureSalt0)
	io.WriteString(h, timestamp)
	io.WriteString(h, SignatureSalt1)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// AcquireParamPid 获取主题上的productID参数,格式resource.productID
func (this *Controller) AcquireParamPid() (int, error) {
	spid := this.Input.Param.Get("productID")
	if spid == "" { // never happen but deal,may be other used
		return 0, errors.New("resource productID invalid")
	}

	pid, err := strconv.Atoi(spid)
	if err != nil { //never happen but deal
		return 0, errors.New("resource productID invalid")
	}

	return pid, nil
}
