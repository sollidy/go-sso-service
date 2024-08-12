package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	compareAndSwap()
}

func compareAndSwap() {
	var (
		counter int64
		wg      sync.WaitGroup
	)
	wg.Add(100)

	// This way, each goroutine can swap the value of `counter` only once.
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			if !atomic.CompareAndSwapInt64(&counter, 0, 2) {
				return
			}
			fmt.Println("Goroutine swapped number is: ", i)
		}(i)
	}

	wg.Wait()
	fmt.Println("final number is: ", counter)
}
