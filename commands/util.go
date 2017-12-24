package commands

import (
	"fmt"
	"math"
	"strconv"
	"testing"
)

func prettyPriceInWan(priceStr string) string {
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return priceStr
	}

	if price < 10000 {
		return priceStr
	}

	priceW := round(float64(price)/10000.0, .5, 2)
	return fmt.Sprintf("%.2fw", priceW)
}

func prettyPriceRound2(priceStr string) string {
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return priceStr
	}

	priceRound := round(float64(price), .5, 2)
	return fmt.Sprintf("%.2f", priceRound)
}

func assertEqualString(t *testing.T, a, b string) {
	if a == b {
		return
	}

	t.Errorf("expect %s == %s", a, b)
}

func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
