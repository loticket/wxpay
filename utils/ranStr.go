package utils

import (
	"crypto/rand"
)

func GenerateNonceStr(lens int) string {
	bytes := make([]byte, lens)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	var symbals string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	symbolsByteLength := byte(len(symbals))
	for i, b := range bytes {
		bytes[i] = symbals[b%symbolsByteLength]
	}
	return string(bytes)
}
