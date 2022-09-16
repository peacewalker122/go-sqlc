package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Randomint(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}
func Randomstring(n int) string {
	var sb strings.Builder
	l := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(l)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func Randomowner() string {
	return Randomstring(6)
}

func Randommoney() int64 {
	return Randomint(0, 1000)
}

func Randomcurrency() string {
	c := []string{"IDR", "USD", "GBP", "EUR"}
	l := len(c)
	return c[rand.Intn(l)]
}

func Randomemail() string {
	return fmt.Sprintf("%v@test.com",Randomstring(6))
}