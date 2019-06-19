package synccall

import (
	"time"

	cache "github.com/thinkgos/go-cache"
	"github.com/thinkgos/snowflake"
)

// 默认配置
const (
	DefaultNode            = 2
	DefaultExpiration      = 5 * time.Second
	DefaultCleanUpInterval = 30 * time.Second
)

type item struct {
	value chan interface{} // 传值通道
}

// Manage 管理
type Manage struct {
	c          *cache.Cache
	n          *snowflake.Node
	expiration time.Duration
}

// Option 配置选项
type Option struct {
	node              int64
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
}

// NewOption 新建配置选项,默认置
func NewOption() *Option {
	return &Option{DefaultNode, DefaultExpiration, DefaultCleanUpInterval}
}

// SetNode 设置snowflake的node
func (this *Option) SetNode(node int64) {
	nodeMax := int64(-1 ^ (-1 << snowflake.NodeBits))
	if node > nodeMax {
		node = nodeMax
	} else if node < 0 {
		node = 0
	}
	this.node = node
}

// SetExpiration 设置缓存item默认超时时间
func (this *Option) SetExpiration(t time.Duration) {
	this.defaultExpiration = t
}

// SetCleanUpInterval 设置缓存item清理时间
func (this *Option) SetCleanUpInterval(t time.Duration) {
	this.cleanupInterval = t
}

// New 创建管理
func New(o ...*Option) *Manage {
	var opt *Option

	if len(o) > 0 {
		opt = o[0]
	} else {
		opt = NewOption()
	}

	sf, _ := snowflake.NewNode(opt.node)

	return &Manage{cache.New(
		opt.defaultExpiration, opt.cleanupInterval),
		sf,
		opt.defaultExpiration,
	}
}

// ObainID 获取唯一ID
func (this *Manage) ObainID() string {
	return this.n.Generate().String()
}

// Done 处理完毕,发送数据
func (this *Manage) Done(id string, value interface{}) {
	itm, ok := this.c.Get(id)
	if !ok {
		return
	}
	select {
	case itm.(*item).value <- value:
	default:
	}
}

// Wait 等待
func (this *Manage) Wait(id string, t ...time.Duration) (interface{}, bool) {
	tm := this.expiration
	if len(t) > 0 {
		tm = t[0]
	}
	itm := &item{
		value: make(chan interface{}, 1),
	}
	this.c.Set(id, itm)

	select {
	case v := <-itm.value:
		return v, true
	case <-time.NewTicker(tm).C:
	}
	return nil, false
}
