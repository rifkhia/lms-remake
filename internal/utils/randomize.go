package utils

import "math/rand"

func GenerateClassKey(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = rune(LETTER_RUNES[rand.Intn(len(LETTER_RUNES))])
	}
	return string(b)
}
