package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"
)

func HashToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}

func GenerateOTP(length int) string {
	rand.NewSource(time.Now().UnixNano())
	digits := "0123456789"
	result := make([]byte, length)
	for i := range length {
		result[i] = digits[rand.Intn(len(digits))]
	}
	return string(result)
}
