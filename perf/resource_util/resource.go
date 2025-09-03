package resource_util

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

func CheckMemoryUsage(obj any) float64 {
	var mBefore, mAfter runtime.MemStats
	runtime.ReadMemStats(&mBefore)

	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Func {
		val.Call(nil)
	}

	runtime.ReadMemStats(&mAfter)
	memUsage := float64(mAfter.TotalAlloc-mBefore.TotalAlloc) / 1024

	if val.Kind() != reflect.Func {
		memUsage = float64(int(val.Type().Size())) / 1024
	}

	fmt.Printf("Memory used: %.2f KB\n", memUsage)
	return memUsage
}

func CheckCPUTime(obj any) time.Duration {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Func {
		start := time.Now()
		val.Call(nil)
		return time.Since(start)
	}

	return 0
}

func ComparePerformanceMetrics(obj1, obj2 any) {
	mem1 := CheckMemoryUsage(obj1)
	cpu1 := CheckCPUTime(obj1)

	mem2 := CheckMemoryUsage(obj2)
	cpu2 := CheckCPUTime(obj2)

	fmt.Printf("Object 1 - Memory: %.2f KB, CPU Time: %v\n", mem1, cpu1)
	fmt.Printf("Object 2 - Memory: %.2f KB, CPU Time: %v\n", mem2, cpu2)

	if mem1 < mem2 {
		fmt.Println("Object 1 is more memory efficient.")
	} else if mem2 < mem1 {
		fmt.Println("Object 2 is more memory efficient.")
	}

	if cpu1 < cpu2 {
		fmt.Println("Object 1 is more CPU efficient.")
	} else if cpu2 < cpu1 {
		fmt.Println("Object 2 is more CPU efficient.")
	} else {
		fmt.Println("Both have similar performance.")
	}
}
