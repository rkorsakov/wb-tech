package main

import (
	"15/internal/shell"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Simple Unix Shell")
	fmt.Println("Type 'exit' to quit")

	sh := shell.NewShell()
	sh.Start()

	fmt.Println("Shell terminated")
	os.Exit(0)
}
