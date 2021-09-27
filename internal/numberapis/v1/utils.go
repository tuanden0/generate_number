package v1

import (
	"math"
	"math/rand"
	"time"
)

func generateRangeNumber(a, b int64) int64 {

	rand.Seed(time.Now().UnixNano())
	max, min := int64(0), int64(0)

	if a > b {
		min = b
		max = a
	} else {
		min = a
		max = b
	}

	// Avoid panic in rand.Int63n
	abs := math.Abs(float64(max - min))
	return min + rand.Int63n(int64(abs)+1)
}
