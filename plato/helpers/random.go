package helpers

import (
	"math/rand"
	"time"
)

func GenerateRandomNumber(length int) int {
	localRandomizer := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(localRandomizer)
	return r1.Intn(length)
}
