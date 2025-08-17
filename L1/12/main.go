package main

import "fmt"

type void struct{}

func main() {
	var n int
	fmt.Scan(&n)
	words := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&words[i])
	}
	set := make(map[string]void)
	for _, word := range words {
		if _, ok := set[word]; !ok {
			set[word] = void{}
		}
	}
	fmt.Println(set)
}
