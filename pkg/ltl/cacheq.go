package ltl

import (
	"errors"
	"strconv"
	"strings"

	"github.com/spf13/cast"
	"github.com/thinkgos/gogate/pkg/tsm"
)

const (
	TsmIDSize     = 256
	TsmIDReseverd = 0
)

type cacheItem struct {
	l     *Ltl_t
	cmdId byte
	val   interface{}
}

// 获取传输Id
func (this *Ltl_t) acquireID(nwkAddr uint16) (uint8, error) {
	var ts *tsm.Tsm

	nwkStr := cast.ToString(nwkAddr)
	v, ok := this.tsmb.Load(nwkStr)
	if !ok {
		ts = tsm.New(TsmIDSize)
		ts.Deprecate(TsmIDReseverd) // 弃用0
		this.tsmb.Store(nwkStr, ts)
	} else {
		ts = v.(*tsm.Tsm)
	}
	id, err := ts.Acquire()
	if err != nil {
		return TsmIDReseverd, err
	}
	return uint8(id), nil
}

// 释放传输Id
func (this *Ltl_t) releaseID(nwkAddr uint16, id uint8) {
	nwkStr := cast.ToString(nwkAddr)
	v, ok := this.tsmb.Load(nwkStr)
	if !ok {
		return
	}
	v.(*tsm.Tsm).Release(uint(id))
}

func (this *Ltl_t) hang(nwk uint16, id uint8, cmd byte, val interface{}) {
	keystr := toKeyString(nwk, id)
	this.c.SetDefault(keystr, &cacheItem{this, cmd, val})
}

func (this *Ltl_t) FindItem(nwk uint16, seq uint8) (byte, interface{}, error) {
	keystr := toKeyString(nwk, seq)
	v, ok := this.c.Get(keystr)
	if !ok {
		return 0, nil, errors.New("not found")
	}
	this.c.Delete(keystr)
	this.releaseID(nwk, seq)
	itm := v.(*cacheItem)
	return itm.cmdId, itm.val, nil
}

func expireCb(k string, val interface{}) {
	nwk, id, err := toSplitKey(k)
	if err != nil {
		return
	}
	sl := val.(*cacheItem).l

	nwkStr := cast.ToString(nwk)
	v, ok := sl.tsmb.Load(nwkStr)
	if !ok {
		return
	}
	ts := v.(*tsm.Tsm)
	ts.Release(uint(id))
	ts.Recover(TsmIDReseverd)
	if ts.IsEmpty() {
		sl.tsmb.Delete(nwkStr)
		return
	}
	ts.Deprecate(TsmIDReseverd)
}

func toKeyString(nwk uint16, id uint8) string {
	return strings.Join([]string{
		cast.ToString(nwk),
		cast.ToString(id),
	}, "_")
}

func toSplitKey(k string) (uint16, uint8, error) {
	s := strings.Split(k, "_")
	if len(s) != 2 {
		return 0, 0, errors.New("split invalid key")
	}

	nwk, err := strconv.Atoi(s[0])
	if err != nil {
		return 0, 0, err
	}

	id, err := strconv.Atoi(s[1])
	if err != nil {
		return 0, 0, err
	}
	return uint16(nwk), uint8(id), nil
}
