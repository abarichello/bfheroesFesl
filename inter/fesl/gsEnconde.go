package fesl

import (
	"math/rand"
	"time"
)

//gs=GameSpy
//rand=rand
const gsLetters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ]["

const (
	gsLettersIdxBits = 6                       // 6 bits to represent a letter index
	gsLettersIdxMask = 1<<gsLettersIdxBits - 1 // All 1-bits, as many as letterIdxBits
	gsLettersIdxMax  = 63 / gsLettersIdxBits   // # of letter indices fitting in 63 bits
)

var randSrc = rand.NewSource(time.Now().UnixNano())

// BF2randUnsafe is safe// you should make your own
func BF2randUnsafe(randLen int) string {
	return BF2rand(randLen, randSrc)
}

// BF2rand generates rand valid BF2string
func BF2rand(randLen int, source rand.Source) string {
	b := make([]byte, randLen)
	// Generates 63 rand bits
	for i, cache, remain := randLen-1, source.Int63(), gsLettersIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = source.Int63(), gsLettersIdxMax
		}
		if idx := int(cache & gsLettersIdxMask); idx < len(gsLetters) {
			b[i] = gsLetters[idx]
			i--
		}
		cache >>= gsLettersIdxBits
		remain--
	}

	return string(b)
}
