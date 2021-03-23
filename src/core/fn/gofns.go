package fn

import (
	"sync"
)

// GoFns
func GoFns(count int, fns []func()) *sync.WaitGroup {
	size := len(fns)
	step := size/count + 1
	wg := &sync.WaitGroup{}
	for i := 0; i < count; i++ {
		start := step * i
		end := step * (i + 1)
		if end > size {
			end = size
		}
		wg.Add(1) //添加一个计数
		go func(start, end int) {
			for j := start; j < end; j++ {
				fn := fns[j]
				if fn != nil {
					fn()
				}
			}
			wg.Done() //减去一个计数

		}(start, end)
	}
	return wg
}
