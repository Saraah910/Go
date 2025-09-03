package main

import (
	"fmt"

	"example.com/performance/resource_util"
)

var result []int

func heavyComputation() []int {
	sum := 0
	for i := 0; i < 1e6; i++ {
		sum += i
	}
	return result
}

func lightComputation() {
	sum := 0
	for i := 0; i < 1e4; i++ {
		sum += i
		result = append(result, sum)
	}
}
func main() {
	// resource_util.ComparePerformanceMetrics(heavyComputation, lightComputation)
	resource_util.CheckMemoryUsage(lightComputation)
	fmt.Printf("Time taken for execution: %v", resource_util.CheckCPUTime(lightComputation))
}
