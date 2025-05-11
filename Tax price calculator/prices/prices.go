package prices

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type TaxIncludedPrices struct {
	TaxRate     float32
	Prices      []float32
	TaxIncluded map[string]float32
}

func (job *TaxIncludedPrices) LoadData() {
	file, err := os.Open("../prices.txt")
	if err != nil {
		fmt.Println("An error occured while opening file.")
		fmt.Println(err)
		return
	}
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		fmt.Print(scanner.Err())
		file.Close()
		return
	}
	prices := make([]float32, len(lines))
	for index, val := range lines {
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			fmt.Println(err)
			file.Close()
			return
		}
		prices[index] = float32(floatVal)
	}
	job.Prices = prices

}
func (job *TaxIncludedPrices) Process() {
	job.LoadData()
	result := make(map[string]string)

	for _, price := range job.Prices {
		resultedPrice := price * (1 + job.TaxRate)
		result[fmt.Sprintf("%.2f", price)] = fmt.Sprintf("%.2f", resultedPrice)
	}
	fmt.Println(result)

}
func New(taxRate float32) *TaxIncludedPrices {
	taxStruct := &TaxIncludedPrices{
		TaxRate: taxRate,
		Prices:  []float32{10, 20, 30},
	}
	return taxStruct
}
