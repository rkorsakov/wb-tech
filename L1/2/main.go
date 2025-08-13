package main

import (
	"fmt"
	"sync"
)

func main() {
	numbers := []int{2, 4, 6, 8, 10}
	wg := sync.WaitGroup{}
	wg.Add(len(numbers))
	for i := range numbers {
		go func(i int) {
			defer wg.Done()
			numbers[i] *= numbers[i]
		}(i)
	}
	wg.Wait()
	fmt.Println(numbers)
}
