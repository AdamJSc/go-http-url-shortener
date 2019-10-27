package shortcodeservice

import (
	"math/rand"
	"time"
)

const shortCodeLength int = 4
const shortCodeSource string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Generate a new short code
func Generate() string {
	var generated string

	// New randomised source
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	index := 0
	var char byte

	for i := 0; i < shortCodeLength; i++ {
		index = r.Intn(len(shortCodeSource))
		char = []byte(shortCodeSource)[index]

		generated += string(char)
	}

	return generated
}
