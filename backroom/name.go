package backroom

import (
	"math/rand"
	"regexp"
	"time"
)

// validName matches a valid name string.
var validName = regexp.MustCompile(`^[a-zA-Z0-9\-]+$`)

// randId returns a string of random letters.
func randId(l int) string {
	n := make([]byte, l)
	for i := range n {
		n[i] = 'a' + byte(rand.Intn(26))
	}
	return string(n)
}

func init() {
	// Seed number generator with the current time.
	rand.Seed(time.Now().UnixNano())
}
