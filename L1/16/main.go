package main

import "fmt"

func quickSort(nums []int) []int {
	length := len(nums)
	pivot := nums[length/2]
	left, right := 0, length-1
	for left <= right {
		for nums[left] < pivot {
			left++
		}
		for nums[right] > pivot {
			right--
		}
		if left <= right {
			nums[left], nums[right] = nums[right], nums[left]
			left++
			right--
		}
	}
	if right > 0 {
		quickSort(nums[:right+1])
	}
	if left < len(nums) {
		quickSort(nums[left:])
	}

	return nums
}

func main() {
	nums := []int{15, 30, 21, 1, 20, 3, 43, 125, 2}
	fmt.Println(quickSort(nums))
}
