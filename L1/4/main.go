package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func work(id int, ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range ch {
		fmt.Printf("Worker %d got msg: %s\n", id, msg)
	}
}

/*
	1. Используем signal.NotifyContext для отслеживания сигнала SIGINT(классический подход)
	2. sync.WaitGroup для гарантии того, что каждый воркер закончит свою работу
	3. select{<-ctx.Done} как классический способ поймать наш SIGINT и обработать его
	4. close(ch) говорит воркерам закончить
*/

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()
	n := os.Args[1]
	amountOfWorkers, err := strconv.Atoi(n)
	if err != nil {
		return
	}
	ch := make(chan string, amountOfWorkers)
	var wg sync.WaitGroup
	for i := 0; i < amountOfWorkers; i++ {
		wg.Add(1)
		go work(i, ch, &wg)
	}
	counter := 0
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case ch <- fmt.Sprintf("hello i am msg number - %d", counter):
			counter++
			time.Sleep(250 * time.Millisecond)
		}
	}
	close(ch)
	wg.Wait()
}
