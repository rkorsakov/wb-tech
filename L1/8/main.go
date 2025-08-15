package main

import "fmt"

func setBit(num int64, i uint, value int) int64 {
	if value == 1 {
		return num | (1 << i)
	} else {
		return num &^ (1 << i)
	}
}

func main() {
	var num int64 = 5
	var i uint = 1
	result := setBit(num, i-1, 0)
	fmt.Printf("Число после установки %d-го бита в 0: %d (%b)\n", i, result, result)
}
