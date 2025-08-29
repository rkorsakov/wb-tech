package main

import (
	"fmt"
	"ntp-time/ntpclient"
	"os"
)

func main() {
	time, err := ntpclient.GetCurrentTime()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting current time: %v", err)
		os.Exit(1)
	}
	fmt.Println(time)
}
