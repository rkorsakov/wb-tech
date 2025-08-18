package main

import (
	"fmt"
	"reflect"
)

func detectType(v interface{}) string {
	switch v.(type) {
	case int:
		return "int"
	case string:
		return "string"
	case bool:
		return "bool"
	default:
		if reflect.TypeOf(v).Kind() == reflect.Chan {
			return "chan"
		}
		return "unknown"
	}

}

func main() {
	var (
		i     = 42
		s     = "hello"
		b     = true
		chInt = make(chan int)
		chStr = make(chan string)
		f     = 3.14
	)

	fmt.Println("Type of i:", detectType(i))
	fmt.Println("Type of s:", detectType(s))
	fmt.Println("Type of b:", detectType(b))
	fmt.Println("Type of chInt:", detectType(chInt))
	fmt.Println("Type of chStr:", detectType(chStr))
	fmt.Println("Type of f:", detectType(f))
}
