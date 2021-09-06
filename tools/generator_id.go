package tools

import (
	"math/rand"
	"strings"
	"time"
)

func GenerateID(nb int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const letterBytes = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	sb := strings.Builder{}
	sb.Grow(nb)
	for ; nb > 0; nb-- {
		sb.WriteByte(letterBytes[rand.Intn(len(letterBytes)-1)])
	}
	return sb.String()
}
