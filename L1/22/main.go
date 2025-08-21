package main

import (
	"fmt"
	"math/big"
)

func sum(a, b *big.Int) *big.Int {
	result := new(big.Int)
	return result.Add(a, b)
}

func product(a, b *big.Int) *big.Int {
	result := new(big.Int)
	return result.Mul(a, b)
}

func diff(a, b *big.Int) *big.Int {
	result := new(big.Int)
	return result.Sub(a, b)
}

func quotient(a, b *big.Int) *big.Int {
	result := new(big.Int)
	return result.Quo(a, b)
}

func main() {
	a := new(big.Int)
	b := new(big.Int)
	fmt.Scan(a, b)
	fmt.Printf("a + b = %d\n", sum(a, b))
	fmt.Printf("a * b = %d\n", product(a, b))
	fmt.Printf("a - b = %d\n", diff(a, b))
	fmt.Printf("a / b = %d\n", quotient(a, b))
}
