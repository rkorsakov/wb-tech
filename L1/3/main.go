package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func work(id int, ch chan string) {
	for msg := range ch {
		fmt.Printf("Worker %d got msg: %s\n", id, msg)
	}
}

func main() {
	n := os.Args[1]
	amountOfWorkers, err := strconv.Atoi(n)
	if err != nil {
		return
	}
	ch := make(chan string, amountOfWorkers)
	for i := 0; i < amountOfWorkers; i++ {
		go work(i, ch)
	}
	counter := 0
	for {
		ch <- fmt.Sprintf("hello i am msg number - %d", counter)
		counter++
		time.Sleep(1 * time.Second)
	}
}
