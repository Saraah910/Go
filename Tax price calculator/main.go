package main

import (
	"example.com/TaxPrice/prices"
)

func main() {

	taxRates := []float32{0, 0.7, 0.1, 0.15}

	for _, taxRate := range taxRates {
		priceJob := prices.New(taxRate)
		priceJob.Process()
	}
}
