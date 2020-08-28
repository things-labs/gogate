/*
 package tsm 提供普通传输序号的生成,它在Acquire 会获取下一个空闲的Id,一般用于生成8bit或
 16bit的传输号生成,对于复杂的,32位以上的传输序号生成将占用大量内存,不建议使用这个包,或将里面的
 bitset实现,改成高效率RoaringBitmap 或采用分布式唯一Id生成包如snowflake,sonyflake
*/
package tsm

import (
	"errors"
	"sync"

	"github.com/willf/bitset"
)

// Tsm 传输序列号对象
type Tsm struct {
	mu     sync.Mutex
	lastID uint
	bitmap *bitset.BitSet
}

// New 创建一个tsm对象,给定传输序列号的长度(即个数)
func New(length uint) *Tsm {
	return &Tsm{lastID: 0, bitmap: bitset.New(length)}
}

// Acquire 获取下一个空闲的id
func (this *Tsm) Acquire() (uint, error) {
	this.mu.Lock()
	id, ok := this.bitmap.NextClear(this.lastID)
	if !ok {
		this.mu.Unlock()
		return 0, errors.New("not free id")
	}
	this.bitmap.Set(id)
	this.lastID = id + 1
	if this.lastID >= this.bitmap.Len() {
		this.lastID = 0
	}
	this.mu.Unlock()
	return id, nil
}

// Release 释放对应Id
func (this *Tsm) Release(id uint) *Tsm {
	this.mu.Lock()
	this.bitmap.Clear(id)
	this.mu.Unlock()
	return this
}

// Deprecate 弃用一个id,用于某个Id保留不用
func (this *Tsm) Deprecate(id uint) *Tsm {
	this.mu.Lock()
	if id < this.bitmap.Len() {
		this.bitmap.Set(id)
	}
	this.mu.Unlock()
	return this
}

// Recover 恢复弃用的Id
func (this *Tsm) Recover(id uint) *Tsm {
	this.mu.Lock()
	this.bitmap.Clear(id)
	this.mu.Unlock()
	return this
}

// ReleaseAll 释放所有Id
func (this *Tsm) ReleaseAll() *Tsm {
	this.mu.Lock()
	this.lastID = 0
	this.bitmap.ClearAll()
	this.mu.Unlock()
	return this
}

// TestUsed 测试对应id是否使用中
func (this *Tsm) TestUsed(id uint) bool {
	this.mu.Lock()
	b := this.bitmap.Test(id)
	this.mu.Unlock()
	return b
}

// TestIdle 测试对应id是否未被使用
func (this *Tsm) TestIdle(id uint) bool {
	this.mu.Lock()
	b := this.bitmap.Test(id)
	this.mu.Unlock()
	return !b
}

// IsFull 是否id全使用完了
func (this *Tsm) IsFull() bool {
	this.mu.Lock()
	b := this.bitmap.All()
	this.mu.Unlock()
	return b
}

// IsEmpty 是否id全未被使用
func (this *Tsm) IsEmpty() bool {
	this.mu.Lock()
	b := this.bitmap.None()
	this.mu.Unlock()
	return b
}
