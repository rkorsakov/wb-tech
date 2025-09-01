package main

import (
	"9/unpacker"
	"fmt"
	"os"
)

func main() {
	var str string
	fmt.Scan(&str)
	str, err := unpacker.UnpackString(str)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(str)
}
