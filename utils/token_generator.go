package utils

import (
	"math/rand"
	"time"
)

func GenerateToken() string {
	str := []rune("qwertyuiop[]asdfghjkl;1234567890-=QWERTYUIOP[]ASDFGHJKL;ZXCVBNM!@#$%^&*()_+zxcvbnm,./")
	tokenRune := []rune{}

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	
	for i := 0 ; i <= 95 ; i++{
		tokenRune = append(tokenRune, str[rng.Intn(len(str))])
	}
	return string(tokenRune)
}


