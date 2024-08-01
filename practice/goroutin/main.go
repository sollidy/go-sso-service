package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// withWaiting()
	// rangeBuffered()
	// selectChan()
	writeChanByTimer()
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

func rangeBuffered() {
	numbers := []int{1, 2, 3, 4, 5}
	bufferdChannel := make(chan int, 2)
	go func() {
		for _, number := range numbers {
			bufferdChannel <- number
		}
		close(bufferdChannel)
	}()
	for number := range bufferdChannel {
		fmt.Println(number)
	}
}

func selectChan() {
	unbufferedChannel := make(chan int)
	go func() {
		time.Sleep(time.Millisecond * 500)
		unbufferedChannel <- 1
	}()

	select {
	case val := <-unbufferedChannel:
		fmt.Println("Unbuffered:", val)

	case <-time.After(time.Millisecond * 600):
		fmt.Println("timeout")
	}
}

func writeChanByTimer() {

	unbufferedChannel := make(chan int)

	go func() {
		timer := time.After(time.Millisecond * 1)
		defer close(unbufferedChannel)
		for i := 0; i < 1000; i++ {
			select {
			case val := <-timer:
				fmt.Println("time is up", val)
				return
			default:
				time.Sleep(time.Nanosecond)
				unbufferedChannel <- i
			}
		}
	}()

	for val := range unbufferedChannel {
		fmt.Println(val)
	}
}
