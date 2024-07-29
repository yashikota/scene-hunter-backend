package util

import (
	"math/rand"
)

// GenerateRandomScore generates a random score
// between 10.0000 ~ 99.9999
func GenerateRandomScore() float32 {
	return float32(rand.Intn(899999)+100000) / 10000
}
