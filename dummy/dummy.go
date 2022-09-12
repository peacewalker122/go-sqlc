package dummy

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomint(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}
func randomstring(n int) string {
	var sb strings.Builder
	l := len(alphabet)
	for i := 0; i < l; i++ {
		c := alphabet[rand.Intn(l)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func Randomowner() string {
	return randomstring(6)
}

func Randommoney() int64 {
	return randomint(0, 1000)
}

func Randomcurrency() string {
	c := []string{"IDR", "USD", "GBP", "EUR"}
	l := len(c)
	return c[rand.Intn(l)]
}

// it's nothing just a dummy to test the git
// second test for the git
func Nothing() {
	println("nothing")
}
