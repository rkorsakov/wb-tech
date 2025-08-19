package main

import (
	"fmt"
	"sort"
)

func binarySearch(nums []int, number int) int {
	left, right := 0, len(nums)-1
	for left <= right {
		mid := (left + right) / 2
		if nums[mid] == number {
			return mid
		}
		if nums[mid] > number {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	return -1
}

func main() {
	nums := []int{10, 12, 1, 15, 23, 102, 234, 123, 531, 124, 22, 15}
	sort.Ints(nums)
	fmt.Println(binarySearch(nums, 12))
	fmt.Println(binarySearch(nums, 999))
}
