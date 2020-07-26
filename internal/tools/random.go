package tools

import (
	"math/rand"
	"time"
)

const startBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~@%^_+=-][}{:/.,"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	b[0] = startBytes[rand.Intn(len(startBytes))]
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
