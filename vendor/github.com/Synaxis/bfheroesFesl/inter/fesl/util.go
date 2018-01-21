package fesl

import (
	"math/rand"
	"time"
)

const gamespyLetters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ]["
const (
	gamespyLettersIdxBits = 6                            // 6 bits to represent a letter index
	gamespyLettersIdxMask = 1<<gamespyLettersIdxBits - 1 // All 1-bits, as many as letterIdxBits
	gamespyLettersIdxMax  = 63 / gamespyLettersIdxBits   // # of letter indices fitting in 63 bits
)

var randSrc = rand.NewSource(time.Now().UnixNano())

// BF2RandomUnsafe is a not thread-safe version of BF2Random
// For thread-safety you should use BF2Random with your own seed
func BF2RandomUnsafe(randomLen int) string {
	return BF2Random(randomLen, randSrc)
}

// BF2Random generates a random string with valid BF2 random chars
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang/31832326
func BF2Random(randomLen int, source rand.Source) string {
	b := make([]byte, randomLen)
	// A source.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := randomLen-1, source.Int63(), gamespyLettersIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = source.Int63(), gamespyLettersIdxMax
		}
		if idx := int(cache & gamespyLettersIdxMask); idx < len(gamespyLetters) {
			b[i] = gamespyLetters[idx]
			i--
		}
		cache >>= gamespyLettersIdxBits
		remain--
	}

	return string(b)
}
