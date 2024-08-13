package pkg

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandom6Char() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, 6)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}
	return string(result)
}
