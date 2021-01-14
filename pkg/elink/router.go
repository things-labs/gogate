package elink

// PublicCenter 全局管理中心
var PublicCenter = NewCenter()

// RegisterChannelSelector 注册通道
func RegisterChannelSelector(channel string) {
	if err := PublicCenter.Put(channel, NewControllerRegister()); err != nil {
		panic(err.Error())
	}
}

// Router 向指定通道注册资源控制器,错误panic
// 添加路由到资源
// 第三个参数就是用来设置对应 method 到函数名，即自定义method到功能函数,定义如下
//-  *表示任意的elink method 都执行该功能函数,不包含其它的自定义的method
//-  使用 method:funcname 格式来展示
//-  多个不同的, 格式使用 ; 分割
//-  多个 method 对应同一个funcname，method之间通过,来分割
//- 如果同时存在 * 和自定义elnik Method对应的功能函数，那么优先执行自定义功能函数
func Router(channel, resource string, c ControllerInterface, mapMethods ...string) *Center {
	if err := PublicCenter.Router(channel, resource, c, mapMethods...); err != nil {
		panic("router: " + err.Error())
	}
	return PublicCenter
}

// ChannelSelectorList 返回已注册的通道列表
func ChannelSelectorList() []string {
	return PublicCenter.ChannelSelectorList()
}

// Server 服务
func Server(p Provider, topic string, payload []byte) {
	if err := PublicCenter.Server(p, topic, payload); err != nil {
		PublicCenter.errorHandler(err)
	}
}

func SetPanicHandler(f func(interface{})) {
	PublicCenter.SetPanicHandler(f)
}
