package main

import "fmt"

func main() {
	numbers := make([]int, 2, 5)
	fmt.Println(numbers)
	numbers = append(numbers, 2)
	numbers = append(numbers, 3)
	numbers = append(numbers, 4)

	fmt.Println(numbers)
	fmt.Println(transform(&numbers, double))
	fmt.Println(transform(&numbers, triple))

}

func transform(numbers *[]int, transformNum func(int) int) []int {
	var Arr []int
	for _, val := range *numbers {
		Arr = append(Arr, transformNum(val))
	}
	return Arr
}

func double(num int) int {
	return num * 2
}

func triple(num int) int {
	return num * 3
}
