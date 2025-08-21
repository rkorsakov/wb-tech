package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func noRepeatStr(s string) bool {
	letters := make(map[rune]struct{}, utf8.RuneCountInString(s))
	for _, char := range s {
		char = unicode.ToLower(char)
		if _, ok := letters[char]; ok {
			return false
		}
		letters[char] = struct{}{}
	}
	return true
}

func main() {
	var s string
	fmt.Scan(&s)
	check := noRepeatStr(s)
	fmt.Println(check)
}
