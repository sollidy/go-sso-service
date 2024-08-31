package main

import (
	"fmt"
	"math/rand/v2"
)

// This Go code snippet runs a Monte Carlo simulation to estimate the average number of attempts it takes to guess a random number (37) between 1 and 100. It does this by:

// 1. Creating a channel to receive the number of attempts from each goroutine.
// 2. Starting 1 million goroutines, each trying to guess the number.
// 3. Collecting the results from each goroutine and storing them in a slice.
// 4. Calculating the average number of attempts and printing the result.

func main() {
	attemptsChannel := make(chan int)

	const (
		targetNumber  = 37
		numGoroutines = 1000000
	)

	for i := 0; i < numGoroutines; i++ {
		go guessRandom(targetNumber, attemptsChannel)
	}

	attempts := make([]int, numGoroutines)
	for i := range attempts {
		attempts[i] = <-attemptsChannel
	}

	averageAttempts := avarage(attempts)
	fmt.Printf("Average guess attempts: %v\n", averageAttempts)
}

func guessRandom(num int, channel chan int) {
	var i int = 0
	for {
		i++
		var rand int = rand.IntN(100) + 1
		switch {
		case num == rand:
			// fmt.Println("quantity of tries: ", i)
			channel <- i
			return
		case num < rand:
			// fmt.Println("less you guess")
		case num > rand:
			// fmt.Println("more you guess")
		}
	}
}

func avarage(numbers []int) int {
	var sum int
	for _, number := range numbers {
		sum += number
	}
	return sum / len(numbers)
}
