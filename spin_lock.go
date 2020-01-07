package spinlock

import (
	"runtime"
	"sync"
	"sync/atomic"
)

// sl值为0的时候表示锁未被锁住，可以lock，值为1时 表示已上锁待unlock
// 重复调用unlock不会导致死锁。
type spinLock int32

func (sl *spinLock) Lock() {
	//cas操作不成功调用Gosched让出cup时间片。cas成功就结束
	for !atomic.CompareAndSwapInt32((*int32)(sl), 0, 1) {
		runtime.Gosched()
	}
}

func (sl *spinLock) Unlock() {
	atomic.StoreInt32((*int32)(sl), 0)
}

func (sl *spinLock) TryLock() bool {
	return atomic.CompareAndSwapInt32((*int32)(sl), 0, 1)
}

func NewSpinLock() sync.Locker {
	return new(spinLock)
}
