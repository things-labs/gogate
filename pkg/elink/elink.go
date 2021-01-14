package elink

import (
	"errors"
	"reflect"
	"strconv"
	"sync"
)

// Provider 接入源接口
type Provider interface {
	// 错误追加在主题上回复,针对mqtt
	ErrorDefaultResponse(topic string) error
	// 写回复
	WriteResponse(topic string, data interface{}) error
}

// Center 所有应用控制中心,任何低层协议源满足一定的接口就可以进行通信
type Center struct {
	l            sync.RWMutex
	ch           map[string]*ControllerRegister // 通道与控制器的映射
	panicHandler func(interface{})
	errorHandler func(error)
}

// NewCenter 创建一个新的数据中心
func NewCenter() *Center {
	return &Center{
		ch:           make(map[string]*ControllerRegister),
		panicHandler: func(interface{}) {},
		errorHandler: func(error) {},
	}
}
func (sf *Center) SetPanicHandler(f func(interface{})) {
	if f != nil {
		sf.panicHandler = f
	}
}

func (sf *Center) SetErrorHandler(f func(error)) {
	if f != nil {
		sf.errorHandler = f
	}
}

// Put 增加一个通道的控制器,有则报错,否则有新的
func (sf *Center) Put(channel string, c *ControllerRegister) error {
	_, exist := sf.Get(channel)
	if exist {
		return errors.New("channel has exist")
	}
	sf.Set(channel, c)
	return nil
}

// Set 增加一个新通道的控制器,不管是否有,都使用新的
func (sf *Center) Set(channel string, c *ControllerRegister) {
	sf.l.Lock()
	sf.ch[channel] = c
	sf.l.Unlock()
}

// Get 获取通道注册控制器
func (sf *Center) Get(channel string) (*ControllerRegister, bool) {
	sf.l.RLock()
	v, ok := sf.ch[channel]
	sf.l.RUnlock()
	return v, ok
}

// DeleteRouter 删除通道的所有资源 或 删除通道的某个资源
func (sf *Center) Delete(channel string, resources ...string) {
	sf.l.Lock()
	if len(resources) == 0 {
		delete(sf.ch, channel)
		sf.l.Unlock()
		return
	}

	cr, isExists := sf.ch[channel]
	if !isExists {
		sf.l.Unlock()
		return
	}
	sf.l.Unlock()

	// 删除指定通道的指定资源实列
	for _, v := range resources {
		cr.DeleteRouter(v)
	}
}

// ChannelSelectorList 获取通道列表
func (sf *Center) ChannelSelectorList() []string {
	sf.l.RLock()
	s := make([]string, 0, len(sf.ch))
	for k := range sf.ch {
		s = append(s, k)
	}
	sf.l.RUnlock()
	return s
}

// Contains 是否有指定的通道
func (sf *Center) Contains(channel string) bool {
	_, exist := sf.Get(channel)
	return exist
}

// Router 向指定通道注册资源控制器
// 添加路由到资源
// 第三个参数就是用来设置对应 method 到函数名，即自定义method到功能函数,定义如下
//-  *表示任意的elink method 都执行该功能函数,不包含其它的自定义的method
//-  使用 method:funcname 格式来展示
//-  多个不同的, 格式使用 ; 分割
//-  多个 method 对应同一个funcname，method之间通过,来分割
//- 如果同时存在 * 和自定义elnik Method对应的功能函数，那么优先执行自定义功能函数
func (sf *Center) Router(channel, resource string, c ControllerInterface, mapMethods ...string) error {
	if len(channel) == 0 || len(resource) == 0 {
		return errors.New("channel or resource empty")
	}

	cr, ok := sf.Get(channel)
	if !ok {
		return errors.New("channel selector no found")
	}

	return cr.AddToRouter(resource, c, mapMethods...)
}

// Server 服务,处理请求
func (sf *Center) Server(pr Provider, tp string, payload []byte) error {
	topic, err := DecodeTopic(tp) //主题不符合要求抛弃此信息
	if err != nil {
		return err
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				sf.panicHandler(err)
			}
		}()

		// 匹配通道
		execCR, ok := sf.Get(topic.Channel)
		if !ok {
			sf.errorHandler(ErrChannelNotMatch)
			return
		}
		// 找到资源后调用找到控制信息,路由匹配
		ctlInfo, resParam, err := execCR.MatchRouter(topic.Resource)
		if err != nil {
			_ = ErrorDefaultResponse(pr, topic, CodeErrSysResourceNotSupport)
			return
		}
		// 由控制信息获得,处理自定义方法的自定义函数
		runMode := ctlInfo.RunMode(topic.Method)
		if runMode == MethodUnknown {
			_ = ErrorDefaultResponse(pr, topic, CodeErrSysMethodNotSupport)
			return
		}
		r := &Request{
			Topic:   topic,
			Param:   resParam,
			Values:  topic.Query(),
			Payload: payload,
		}

		// 获取控制器实例
		execController := ctlInfo.initialize()
		execController.Init(r, pr) // 执行任何实例方法,先调用Init

		// prepare something
		execController.Prepare()
		switch runMode {
		case MethodGet:
			execController.Get()
		case MethodPost:
			execController.Post()
		case MethodDelete:
			execController.Delete()
		case MethodPut:
			execController.Put()
		default:
			vc := reflect.ValueOf(execController)
			method := vc.MethodByName(runMode)
			// if !method.IsValid() { // 方法不支持
			// 	logs.Error("serverControl:resource(%s) Method( %s ) not implemented", r.Topic.Resource, r.Topic.Method)
			// 	execController.ErrorResponse(101)
			// 	break
			// }
			method.Call(nil)
		}
		// 处理完所有路由,释放相应资源
		execController.Finish()
	}()
	return nil
}

// ErrorDefaultResponse 默认级别的回复,主要用于主题错误的回复,错误直接加在主题上进行返回
func ErrorDefaultResponse(p Provider, tp *TopicLayer, code int) error {
	errMsg := CodeErrorMessage(code)
	v := NewQueryValues()
	v.Set("code", strconv.FormatInt(int64(code), 10))
	v.Set("codeDetail", errMsg.Detail)
	v.Set("message", errMsg.Message)

	return p.ErrorDefaultResponse(EncodeReplyTopic(tp, v.EncodeQuery()))
}
