package house

import "testing"

func TestConvertPrice(t *testing.T) {
	assertEqualString(t, convertPrice(1000000), "100w")
	assertEqualString(t, convertPrice(1000010), "100.001w")
	assertEqualString(t, convertPrice(1230), "0.123w")
}

func assertEqualString(t *testing.T, a, b string) {
	if a == b {
		return
	}

	t.Errorf("expect %s == %s", a, b)
}
