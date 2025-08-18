package main

import "fmt"

func swap(a, b *int) {
	*a += *b
	*b = *a - *b
	*a -= *b
}

func main() {
	var a, b int
	fmt.Scan(&a, &b)
	fmt.Printf("%d %d\n", a, b)
	swap(&a, &b)
	fmt.Printf("%d %d\n", a, b)
}
