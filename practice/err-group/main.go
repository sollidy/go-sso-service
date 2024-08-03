package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	errorGroup()
}

func errorGroup() {

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		time.Sleep(time.Second)

		select {
		case <-ctx.Done():
		default:
			fmt.Println("first started")
			time.Sleep(time.Second)
		}
		return nil
	})

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		default:
			fmt.Println("second started")
		}
		return fmt.Errorf("second goroutin error")
	})

	eg.Go(func() error {
		select {
		case <-ctx.Done():
		default:
			fmt.Println("third started")
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		fmt.Println(err)
	}
}
