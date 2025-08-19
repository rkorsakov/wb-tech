package main

import (
	"fmt"
	"sync"
)

type ConcCounter struct {
	count int
	mu    sync.Mutex
}

func (c *ConcCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *ConcCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func main() {
	counter := &ConcCounter{}
	wg := sync.WaitGroup{}
	numGoroutines := 5000
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			counter.Inc()
		}()
	}
	wg.Wait()
	fmt.Println(counter.Value())
}
