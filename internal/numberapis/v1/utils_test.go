package v1

import "testing"

func Test_generateRangeNumber(t *testing.T) {

	min := int64(-10)
	max := int64(10)

	res := generateRangeNumber(min, max)

	if res < min || res > max {
		t.Errorf("generate number failed, got %v not in range %v, %v", res, min, max)
	}
}
