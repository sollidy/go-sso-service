package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	withWaiting()
	fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))
}

func withWaiting() {
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(i)
		}()
	}
	wg.Wait()
}
