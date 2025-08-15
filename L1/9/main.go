package main

import (
	"fmt"
	"sync"
)

func main() {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	wg := sync.WaitGroup{}
	numsCh := make(chan int)
	squared := make(chan int)
	wg.Add(len(arr))
	for i := range arr {
		go func(i int) {
			defer wg.Done()
			numsCh <- arr[i]
		}(i)
	}
	go func() {
		wg.Wait()
		close(numsCh)
	}()
	go func() {
		for num := range numsCh {
			squared <- num * num
		}
		close(squared)
	}()

	for res := range squared {
		fmt.Println(res)
	}
}
