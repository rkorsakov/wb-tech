package main

import (
	"fmt"
	"ntpcustom/ntpclient"
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
