package main

import "fmt"

func deleteFromSlice[A any](arr []A, index int) []A {
	copy(arr[index:], arr[index+1:])
	return arr[:len(arr)-1]
}

func main() {
	arr := []int{1, 2, 3, 4}
	arr = deleteFromSlice(arr, 2)
	fmt.Println(arr)
}
