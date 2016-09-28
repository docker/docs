package signature

import (
	"math/rand"
	"time"
)

// RandomKeyCharacters is a []byte of the characters to choose from when generating
// random keys.
var RandomKeyCharacters []byte = []byte("abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")

// randomised gets whether the rand.Seed has been set or not.
var randomised bool

// RandomKey generates a random key at the given length.
//
// The first time this is called, the rand.Seed will be set
// to the current time.
func RandomKey(length int) string {

	// randomise the seed
	if !randomised {
		rand.Seed(time.Now().UTC().UnixNano())
		randomised = true
	}

	// Credit: http://stackoverflow.com/questions/12321133/golang-random-number-generator-how-to-seed-properly

	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		randInt := randInt(0, len(RandomKeyCharacters))
		bytes[i] = RandomKeyCharacters[randInt : randInt+1][0]
	}
	return string(bytes)

}

// randInt generates a random integer between min and max.
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
