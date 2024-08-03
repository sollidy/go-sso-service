package main

import (
	"context"
	"runtime"
	"sync"
	"time"
)

func main() {
	workerPool()
}

func workerPool() {
	ctx, cancel := context.WithCancel(context.Background())
	// ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	wg := &sync.WaitGroup{}

	numbersToProcess, processNumbers := make(chan int, 5), make(chan int, 5)

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(ctx, numbersToProcess, processNumbers)
		}()
	}

	go func() {
		for i := 0; i <= 1000; i++ {
			numbersToProcess <- i
		}
		close(numbersToProcess)
	}()

	// This goroutine waits for all the worker goroutines to finish and then closes the processNumbers channel.
	// This is done to signal to the main goroutine that it can stop waiting for values from the channel.
	go func() {
		wg.Wait()
		close(processNumbers)
	}()

	var counter int
	for value := range processNumbers {
		counter++
		println(value)
	}
	println(counter)

}

// It runs in a loop until the context is canceled or the channel is closed.
// It repeatedly receives a number from the toProcess channel.
// If the context is canceled or the channel is closed, the function returns.
// Otherwise, it squares the received number and sends it to the processed channel after a short delay.
func worker(ctx context.Context, toProcess <-chan int, processed chan<- int) {
	for {
		select {
		case <-ctx.Done():
			return
		case value, ok := <-toProcess:
			if !ok {
				return
			}
			time.Sleep(time.Millisecond)
			processed <- value * value
		}
	}
}
