package main

import (
	"context"
	"fmt"
	"time"
)

func sleep(duration time.Duration) {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	<-timer.C
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				fmt.Println("goroutine is working")
				sleep(1 * time.Second)
			}
		}
	}()
	sleep(5 * time.Second)
}
