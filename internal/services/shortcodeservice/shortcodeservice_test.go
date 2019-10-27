package shortcodeservice

import (
	"strings"
	"testing"
)

func TestItGeneratesAValidShortCode(t *testing.T) {
	shortCode := Generate()

	expectedLength := 4
	if len(shortCode) != expectedLength {
		t.Errorf("Expected shortcode of length %d, instead received %d", expectedLength, len(shortCode))
	}

	expectedSource := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	for _, char := range []byte(shortCode) {
		if strings.Contains(expectedSource, string(char)) == false {
			t.Errorf("Unexpected character '%s' in shortcode", string(char))
		}
	}
}
