package fn

import (
	"runtime/debug"
	"sync"
)

// GoFns
func GoFns(num int, fns []func()) *sync.WaitGroup {
	size := len(fns)
	step := size/num + 1
	wg := &sync.WaitGroup{}
	for i := 0; i < num; i++ {
		start := step * i
		end := step * (i + 1)
		if end > size {
			end = size
		}
		wg.Add(1) //添加一个计数
		go func(fns []func()) {
			for _, fn := range fns {
				fn()
				// GC 太慢了
				debug.FreeOSMemory()
			}
			wg.Done() //减去一个计数
		}(fns[start:end])
	}
	return wg
}

// GoFnsGroupBy
func GoFnsGroupBy(count int, num int, fns []func()) *sync.WaitGroup {
	size := len(fns)
	if size > count {
		start := 0
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			for start < size {
				end := start + count
				if end > size {
					end = size
				}
				GoFns(num, fns[start:end]).Wait()
				start += count
			}
			wg.Done()
		}()
		return wg
	}
	return GoFns(num, fns)
}
