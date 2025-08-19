package main

import (
	"fmt"
	"unicode/utf8"
)

func reverseString(s string) string {
	size := utf8.RuneCountInString(s)
	reversed := make([]rune, size)
	tail := size - 1
	for _, val := range s {
		reversed[tail] = val
		tail--
	}
	return string(reversed)
}

func main() {
	str := "привет мир"
	fmt.Println(reverseString(str))
}
