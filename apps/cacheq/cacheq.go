package cacheq

import (
	"errors"
	"strconv"
	"time"

	"github.com/slzm40/common"
	"github.com/slzm40/go-cache"
	"github.com/slzm40/tsmanage"

	"github.com/astaxie/beego/logs"
)

type Cacheq struct {
	*tsmanage.Tsm
	c *cache.Cache
}

type CacheqItem struct {
	Pkid    int
	IsLocal bool // if local do not send message to up
	Cb      func(*CacheqItem) error
	Val     interface{}
	cq      *Cacheq
}

var chcheq *Cacheq

func init() {
	c := cache.New(1*time.Second, 30*time.Second)
	c.OnEvicted(expireCb)
	ts := tsmanage.New(256)
	ts.Deprecate(0) // 保留节点0
	chcheq = &Cacheq{ts, c}
}

func AllocID() (uint8, error) {
	v, err := chcheq.Alloc()
	return uint8(v), err
}

func FreeID(id uint8) {
	chcheq.Free(uint(id))
}

func Hang(id uint8, ci *CacheqItem) {
	ci.cq = chcheq
	chcheq.c.Set(common.FormatBaseTypes(id), ci)
}

func Excute(id uint8) (*CacheqItem, error) {
	ids := common.FormatBaseTypes(id)
	v, find := chcheq.c.Get(ids)
	if !find {
		logs.Debug("excute not find")
		return nil, errors.New("no this transfer id")
	}
	chcheq.Free(uint(id))
	chcheq.c.Delete(ids)
	cv := v.(*CacheqItem)
	//	if cv.Cb != nil {
	//		return cv.Cb(cv)
	//	}

	return cv, nil
}

func expireCb(k string, val interface{}) {
	c := val.(*CacheqItem).cq
	logs.Debug("free key: %s", k)
	id, err := strconv.ParseUint(k, 0, 8)
	if err != nil {
		return
	}
	c.Free(uint(id))
}
