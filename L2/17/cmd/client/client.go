package main

import (
	"17/internal/client"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := parseFlag()
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go HOST PORT")
		os.Exit(1)
	}
	host := os.Args[1]
	port := os.Args[2]
	cl := client.NewClient(timeout)
	err := cl.Connect(host, port)
	if err != nil {
		fmt.Printf("Error connecting to %s: %s\n", host, err)
		os.Exit(1)
	}
	defer cl.Close()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go cl.StartInputHandler()
	go cl.StartConnectionHandler()

	select {
	case <-cl.Done():
		fmt.Println("Connection closed")
	case <-sigCh:
		fmt.Println("\nReceived shutdown signal")
		cl.Close()
	}
}

func parseFlag() time.Duration {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
	flag.Parse()
	return timeout
}
