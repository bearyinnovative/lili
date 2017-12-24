package commands

import "testing"

func TestPrettyPriceInWan(t *testing.T) {
	assertEqualString(t, prettyPriceInWan("105151.74"), "10.52w")
	assertEqualString(t, prettyPriceInWan("10191.74"), "1.02w")
	assertEqualString(t, prettyPriceInWan("1011.74"), "1011.74")
}

func TestPriceRound2(t *testing.T) {
	assertEqualString(t, prettyPriceRound2("13925.9"), "13925.90")
	assertEqualString(t, prettyPriceRound2("91564.881385"), "91564.88")
	assertEqualString(t, prettyPriceRound2("685.216"), "685.22")
}
