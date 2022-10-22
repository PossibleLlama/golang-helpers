package strings

import (
	"crypto/rand"
	"strings"
)

const (
	alphabeticBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	hexadecimalBytes = "abcdef1234567890"
	letterIdxBits    = 6
	letterIdxMask    = 1<<letterIdxBits - 1
)

// RandAlphabeticString Generator function of a random series of characters
// Uses a-zA-Z character set
func RandAlphabeticString(n int) string {
	return randString(n, alphabeticBytes)
}

// RandHexAlphaNumericString Generator function of a random series of characters
// Uses a-f0-9 character set
func RandHexAlphaNumericString(n int) string {
	return randString(n, hexadecimalBytes)
}

// From https://stackoverflow.com/a/35615565
func randString(n int, characterSet string) string {
	if n <= 0 {
		return ""
	}
	result := make([]byte, n)
	bufferSize := int(float64(n) * 1.3)
	for i, j, randomBytes := 0, 0, []byte{}; i < n; j++ {
		if j%bufferSize == 0 {
			randomBytes = mustSecureRandomBytes(bufferSize)
		}
		if idx := int(randomBytes[j%n] & letterIdxMask); idx < len(characterSet) {
			result[i] = characterSet[idx]
			i++
		}
	}

	return string(result)
}

// mustSecureRandomBytes returns the requested number of bytes using crypto/rand
func mustSecureRandomBytes(length int) []byte {
	var randomBytes = make([]byte, length)
	// #nosec G104 -- Ignoring for now
	rand.Read(randomBytes)
	return randomBytes
}

// Append takes a string and appends a string to it
func Append(s string, args ...string) string {
	var sb strings.Builder
	sb.WriteString(s)

	for _, arg := range args {
		sb.WriteString(arg)
	}
	return sb.String()
}
