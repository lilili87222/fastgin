package util

import (
	"math/rand"
	"regexp"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) string {
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

// IsPhoneNumber checks if the input string is a valid phone number with country code
func IsPhoneNumber(phone string) bool {
	// Define the regex pattern for a valid international phone number with country code
	// This pattern matches phone numbers with an optional country code prefix
	var re = regexp.MustCompile(`^\+?[1-9]\d{0,2}[- ]?\d{10,14}$`)
	return re.MatchString(phone)
}
