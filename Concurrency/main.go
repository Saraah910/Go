package main

import "sync"

func main() {
	area := make(chan float64)
	vol := make(chan float64)
	peri := make(chan float64)

	inputString := DataInput("Enter the side: ")

	var wg sync.WaitGroup

	wg.Add(1)
	go func(side float64) {
		defer wg.Done()
		area <- CalArea(side)
		close(area)
	}(inputString)

	wg.Add(1)
	go func(side float64) {
		defer wg.Done()
		vol <- CalCube(side)
		close(vol)
	}(inputString)

	wg.Add(1)
	go func(side float64) {
		defer wg.Done()
		peri <- CalPeri(side)
		close(peri)
	}(inputString)

	a := <-area
	v := <-vol
	p := <-peri
	wg.Wait()

	process(a, v, p)
}
