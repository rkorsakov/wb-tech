package main

import (
	"fmt"
	"math"
)

type Point struct {
	x, y float64
}

func NewPoint(x, y float64) *Point {
	return &Point{x, y}
}

func (p *Point) Distance(other Point) float64 {
	return math.Sqrt(math.Pow(other.x-p.x, 2) + math.Pow(other.y-p.y, 2))
}

func main() {
	point1 := NewPoint(1, 4)
	point2 := NewPoint(3, 15)
	fmt.Println(point1.Distance(*point2))
}
