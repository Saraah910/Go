package main

import (
	"fmt"
	"time"
)

func DataInput(inputData string) float64 {
	var side float64
	for {
		fmt.Println(inputData)
		s, err := fmt.Scanln(&side)
		if err != nil || s == 0 || inputData == "" {
			fmt.Println("Input cannot be empty. Please try again.")
			continue
		}
		break
	}
	return side

}

func process(area chan float64, cube chan float64, peri chan float64) {
	fmt.Printf("The details: \nArea: %v\nCube: %v\nPerimeter: %v", <-area, <-cube, <-peri)
}

func CalArea(side float64) float64 {
	time.Sleep(time.Millisecond * 200)
	return side * side
}

func CalCube(side float64) float64 {
	time.Sleep(time.Millisecond * 300)
	return side * side * side
}

func CalPeri(side float64) float64 {
	time.Sleep(time.Millisecond * 400)
	return 4 * side
}
