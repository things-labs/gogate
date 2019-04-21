package ewait

import (
	"sync/atomic"
	"time"

	"github.com/thinkgos/go-cache"
	"github.com/thinkgos/snowflake"
)

const defaultTime = 5 * time.Second

type EvWait struct {
	value chan interface{} // 传值通道
	once  int32            // 通道只关一次
	tm    time.Duration
}
type evManage struct {
	*cache.Cache
	*snowflake.Node
}

var em *evManage

func init() {
	sf, _ := snowflake.NewNode(2)
	c := cache.New(5*time.Second, 30*time.Second)
	c.OnEvicted(cleanUp)
	em = &evManage{
		Cache: c,
		Node:  sf,
	}
}

func cleanUp(id string, value interface{}) {
	close(value.(*EvWait).value)
}

func ObainID() string {
	return em.Generate().String()
}

func Add(id string, t ...time.Duration) *EvWait {
	ew := &EvWait{value: make(chan interface{}), tm: defaultTime}
	if len(t) > 0 {
		ew.tm = t[0]
	}
	em.Set(id, ew)
	return ew
}

func Done(id string, value interface{}) {
	itm, ok := em.Get(id)
	if !ok {
		return
	}
	ev := itm.(*EvWait)
	if atomic.CompareAndSwapInt32(&ev.once, 0, 1) {
		ev.value <- value
	}
}

func (this *EvWait) Wait() (interface{}, bool) {
	select {
	case v, ok := <-this.value:
		return v, ok
	case <-time.NewTicker(this.tm).C:
	}

	return nil, false
}
