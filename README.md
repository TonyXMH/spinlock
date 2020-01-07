# spinlock

spinlock 自旋锁，在等待锁资源时会占用cpu，不会让占用锁的进程进入睡眠状态。因为这点spinlock效率比mutex高，但是会占用cpu资源导致cpu占用过高。
PS:请尽可能的减少临界区

spinlock 基于go实现，
atomic.CompareAndSwapInt32()cas 值与old比较相等则更新新值
atomic.StoreInt32()原子性的覆盖新值
runtime.Gosched()让出cpu时间片

```go test -bench=.```
```go
BenchmarkUncontested-4                  100000000               16.5 ns/op
BenchmarkSyncAtomicUncontested-4        100000000               12.0 ns/op
BenchmarkParallelAdd-4                  50000000                30.4 ns/op
BenchmarkAtomicParallelAdd-4            10000000               172 ns/op
BenchmarkMutexParallelAdd-4             30000000                75.9 ns/op

```
多加runtime.Gosched()的性能更差了。
其实这个压测没有实际意义由于没有比较cpu占用率

