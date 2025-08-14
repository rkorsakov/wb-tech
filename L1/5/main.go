package main

import (
	"context"
	"fmt"
	"time"
)

func runForDuration(N time.Duration) {
	ch := make(chan int)
	ctx, cancel := context.WithTimeout(context.Background(), N)
	defer cancel()

	go func() {
		defer close(ch)
		i := 0
		for {
			select {
			case <-ctx.Done():
				return
			default:
				ch <- i
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	for val := range ch {
		fmt.Println("Got:", val)
	}
	fmt.Println("Timeout after: ", N)
}

func main() {
	runForDuration(1 * time.Second)
}
