package tsm

import (
	"testing"
)

var tsm = New(100)

func TestTsm(t *testing.T) {
	t.Run("tsm", func(t *testing.T) {
		id, err := tsm.Acquire()
		if err != nil {
			t.Errorf("Acquire() err = %v, want %v", err, nil)
		}
		if !tsm.TestUsed(id) {
			t.Errorf("TestUsed() = %v, want %v", false, true)
		}
		if !tsm.TestIdle(99) {
			t.Errorf("TestIdle() = %v, want %v", false, true)
		}
		if tsm.IsFull() {
			t.Errorf("IsFull() = %v, want %v", true, false)
		}
		if tsm.IsEmpty() {
			t.Errorf("IsEmpty() = %v, want %v", true, false)
		}
		tsm.Release(id)
		if tsm.TestUsed(id) {
			t.Errorf("TestUsed() = %v, want %v", true, false)
		}
		if tsm.IsFull() {
			t.Errorf("IsFull() = %v, want %v", true, false)
		}
		if !tsm.IsEmpty() {
			t.Errorf("IsEmpty() = %v, want %v", false, true)
		}
		for i := 0; i < 100; i++ {
			_, _ = tsm.Acquire()
		}

		_, err = tsm.Acquire()
		if err == nil {
			t.Errorf("Acquire() should be = %v, want %v", nil, err)
		}
		if !tsm.ReleaseAll().IsEmpty() {
			t.Errorf("Acquire() should be = %v, want %v", true, false)
		}
	})

	t.Run("弃用和恢复id", func(t *testing.T) {
		if tsm.Deprecate(0).IsEmpty() {
			t.Errorf("Deprecate 0 then IsEmpty() = %v, want %v", true, false)
		}
		if !tsm.Recover(0).IsEmpty() {
			t.Errorf("Recover 0 then IsEmpty() = %v, want %v", false, true)
		}
	})
}
