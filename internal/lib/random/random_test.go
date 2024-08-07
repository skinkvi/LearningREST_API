package random

import (
	"math/rand"
	"testing"
	"time"
)

func TestNewRandomString(t *testing.T) {
	length := 10
	res := NewRandomString(length)

	if len(res) != length {
		t.Errorf("Expected length %d, but got %d", length, len(res))
	}
}

// Только допустимые символы
func TestNewRandomStringCharacters(t *testing.T) {
	length := 10
	res := NewRandomString(length)

	for _, char := range res {
		if !contains(letterRunes, char) {
			t.Errorf("String contains invalid character: %c", char)
		}
	}
}

func contains(slice []rune, char rune) bool {
	for _, r := range slice {
		if r == char {
			return true
		}
	}
	return false
}

func TestNewRandomStringSeed(t *testing.T) {
	length := 10
	rand.Seed(time.Now().UnixNano())
	result1 := NewRandomString(length)
	rand.Seed(time.Now().UnixNano())
	result2 := NewRandomString(length)
	if result1 == result2 {
		t.Errorf("Expected different strings, but got the same: %s", result1)
	}
}
