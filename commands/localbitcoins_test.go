package commands

import "testing"

func TestPrettyPrice(t *testing.T) {
	assertEqualString(t, prettyPrice("105151.74"), "10.52w")
	assertEqualString(t, prettyPrice("10191.74"), "1.02w")
	assertEqualString(t, prettyPrice("1011.74"), "1011.74")
}

func assertEqualString(t *testing.T, a, b string) {
	if a == b {
		return
	}

	t.Errorf("expect %s == %s", a, b)
}
