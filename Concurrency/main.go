package main

import "sync"

func main() {

	side := make(chan float64)
	area := make(chan float64)
	vol := make(chan float64)
	peri := make(chan float64)

	var inputwg sync.WaitGroup
	inputwg.Add(1)
	go func() {
		defer inputwg.Done()
		s := DataInput("Enter the side:")
		side <- s
		close(side)
	}()
	inputwg.Wait()

	var wg sync.WaitGroup
	s := <-side
	wg.Add(1)
	go func(side float64) {
		defer wg.Done()
		area <- CalArea(side)
	}(s)

	wg.Add(1)
	go func(side float64) {
		defer wg.Done()
		vol <- CalCube(side)
	}(s)

	wg.Add(1)
	go func(side float64) {
		defer wg.Done()
		peri <- CalPeri(side)
	}(s)

	wg.Wait()

	process(area, vol, peri)
}
