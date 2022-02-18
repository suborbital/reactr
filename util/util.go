package util

import (
	"crypto/rand"
	"math/big"
)

// GenerateResultID generates an ID string
func GenerateResultID() string {
	available := "abcdefghijklmnopqrstuvwxyz0123456789"

	id := ""

	for i := 0; i < 24; i++ {
		bigint, _ := rand.Int(rand.Reader, big.NewInt(36))
		index := int(bigint.Int64()) // oh, the hoops you have to jump through....

		id += string(available[index])
	}

	return id
}

// ContainsInt returns true if value present in int slice
func ContainsInt(value int, values []int) bool {
	for _, p := range values {
		if p == value {
			return true
		}
	}
	return false
}
