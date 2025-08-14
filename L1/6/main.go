package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	//1. Выход по условию
	flag := true
	go func() {
		for flag {
			fmt.Println("goroutine is working...")
			time.Sleep(1 * time.Second)
		}
	}()
	time.Sleep(2 * time.Second)
	flag = false

	//2. Канал завершения
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				fmt.Println("goroutine is working...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
	time.Sleep(2 * time.Second)
	close(quit)

	//3. Контекст(cancel)
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				fmt.Println("goroutine is working...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}(ctx)
	time.Sleep(2 * time.Second)
	cancel()

	//4. Контекст(timeout)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				fmt.Println("goroutine is working...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
	time.Sleep(2 * time.Second)

	//5. runtime.Goexit()
	timeout := time.After(3 * time.Second)
	go func() {
		for {
			select {
			default:
				fmt.Println("goroutine working...")
				time.Sleep(1 * time.Second)
			case <-timeout:
				runtime.Goexit()
			}
		}
	}()
	time.Sleep(5 * time.Second)
	//6. panic()
	timeout := time.After(3 * time.Second)
	go func() {
		for {
			select {
			default:
				fmt.Println("goroutine working...")
				time.Sleep(1 * time.Second)
			case <-timeout:
				panic("goroutine is stopped")
			}
		}
	}()
	time.Sleep(5 * time.Second)

	//6. Через закрытие канала данных
	wg := sync.WaitGroup{}
	dataCh := make(chan int)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case _, ok := <-dataCh:
				if !ok {
					return
				}
			default:
				fmt.Println("goroutine is working...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
	time.Sleep(1 * time.Second)
	dataCh <- 42
	time.Sleep(1 * time.Second)
	close(dataCh)
	wg.Wait()
}
