package random

import (
	"math/rand"
	"time"
)

// GetRandom returns a random int
func GetRandom(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}
