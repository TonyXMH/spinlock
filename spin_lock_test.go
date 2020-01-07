package spinlock

import (
	"sync"
	"sync/atomic"
	"testing"
)

const(
	IncNum = 100000
)


//开100个goroutine 分别对money字段进行++操作。
func TestConcurrentLock(t *testing.T)  {
	l:=NewSpinLock()
	var money int64
	var wg sync.WaitGroup
	incFunc:= func() {
		for i:=0;i< IncNum ;i++  {
			l.Lock()
			money++
			l.Unlock()
		}
		wg.Done()
	}
	gNum:=100
	wg.Add(gNum)
	for i:=0;i<gNum ;i++  {
		go incFunc()
	}
	wg.Wait()
	want:=int64(gNum * IncNum)
	if money != want{
		t.Fatal("want",want,"got",money)
	}
}

//####################################################################
//压测内容go test -bench=.
func BenchmarkUncontested(b*testing.B)  {
	l:=NewSpinLock()
	for i:=0;i<b.N ;i++  {
		l.Lock()
		l.Unlock()
	}
}
//压测不调用runtime.Gosched
func BenchmarkSyncAtomicUncontested(b *testing.B)  {
	var l int32
	for i:=0;i<b.N ;i++  {
		for !atomic.CompareAndSwapInt32(&l,0,1){}
		atomic.StoreInt32(&l,0)
	}
}

func BenchmarkParallelAdd(b *testing.B)  {
	l:=NewSpinLock()
	var i int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next(){
			l.Lock()
			i++
			l.Unlock()
		}
	})
}

func BenchmarkAtomicParallelAdd(b *testing.B)  {
	var l int32
	var i int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next(){
			for !atomic.CompareAndSwapInt32(&l,0,1){}
			i++
			atomic.StoreInt32(&l,0)
		}
	})
}

func BenchmarkMutexParallelAdd(b *testing.B)  {
	lock:=new(sync.Mutex)
	var i int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next(){
			lock.Lock()
			i++
			lock.Unlock()
		}

	})
}