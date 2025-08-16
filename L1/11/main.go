package main

import "fmt"

type void struct{}

func setsIntersection(firstSet, secondSet map[int]void) map[int]void {
	intersection := make(map[int]void)
	for number := range firstSet {
		if _, ok := secondSet[number]; ok {
			intersection[number] = firstSet[number]
		}
	}
	return intersection
}

func main() {
	var n int
	fmt.Scan(&n)
	firstSet := make(map[int]void, n)
	secondSet := make(map[int]void, n)
	var number int
	var member void
	for i := 0; i < n; i++ {
		fmt.Scan(&number)
		if _, ok := firstSet[number]; !ok {
			firstSet[number] = member
		}
	}
	for i := 0; i < n; i++ {
		fmt.Scan(&number)
		if _, ok := secondSet[number]; !ok {
			secondSet[number] = member
		}
	}
	fmt.Println(firstSet)
	fmt.Println(secondSet)
	fmt.Println(setsIntersection(firstSet, secondSet))
}
