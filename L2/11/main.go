package main

import (
	"11/anagram"
	"fmt"
)

func main() {
	var n int //кол-во слов
	fmt.Scan(&n)
	words := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&words[i])
	}
	res := anagram.Search(words)
	fmt.Println(res)
}
