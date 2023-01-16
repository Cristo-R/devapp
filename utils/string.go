package utils

import "math/rand"

var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyz")

// RandString random string of length n
func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
