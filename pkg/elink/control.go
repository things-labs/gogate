package elink

import (
	"errors"
)

// 通道定义
const (
	ChannelRaw      = "raw"
	ChannelRawReply = "raw_reply"
	ChannelRawData  = "raw_data"
	ChannelInternal = "internal"
)

// ErrAbort custom error when user stop request handler manually.
var ErrAbort = errors.New("User stop run")

// ensure Controller implement ControllerInterface
var _ ControllerInterface = (*Controller)(nil)

// ControllerInterface 统一所有控制器处理的接口
type ControllerInterface interface {
	Init(r *Request, p Provider)
	Prepare()
	Get()
	Post()
	Delete()
	Put()
	Finish()
	ErrorResponse(code int) error
}

// Controller 控制器
type Controller struct {
	Input *Request
	p     Provider
}

// Init 前期资源初始化
func (sf *Controller) Init(r *Request, p Provider) {
	sf.Input = r
	sf.p = p
}

// Prepare 预准备,可用于认证等前置准备工作
func (sf *Controller) Prepare() {}

// Get 查询
func (sf *Controller) Get() {
	_ = sf.ErrorResponse(CodeErrSysResourceMethodNotImplemented)
}

// Post 新增
func (sf *Controller) Post() {
	_ = sf.ErrorResponse(CodeErrSysResourceMethodNotImplemented)
}

// Put 更新
func (sf *Controller) Put() {
	_ = sf.ErrorResponse(CodeErrSysResourceMethodNotImplemented)
}

// DeleteRouter 删除
func (sf *Controller) Delete() {
	_ = sf.ErrorResponse(CodeErrSysResourceMethodNotImplemented)
}

// Finish 结束,用于清理和释放资源
func (sf *Controller) Finish() {}

// ErrorResponse 错误回复,默认是加在topic上,仅错误进行回复,
func (sf *Controller) ErrorResponse(code int) error {
	return ErrorDefaultResponse(sf.p, sf.Input.Topic, code)
}

// WriteResponse 应答信息
func (sf *Controller) WriteResponse(tp string, data interface{}) error {
	return sf.p.WriteResponse(tp, data)
}

// StopRun 停止当前任务,用于中断一些处理
func (sf *Controller) StopRun() {
	panic(ErrAbort)
}
