package main

import (
	"fmt"
	"math"
)

func main() {
	tmps := []float32{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}
	groups := make(map[int][]float32)
	var group int
	for _, temp := range tmps {
		if temp < 0 {
			group = int(math.Ceil(float64(temp/10))) * 10
		} else {
			group = int(math.Floor(float64(temp/10))) * 10
		}
		groups[group] = append(groups[group], temp)
	}
	fmt.Println(groups)
}
