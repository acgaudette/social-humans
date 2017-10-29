package smhb

import (
	"crypto/rand"
	"fmt"
)

type Token struct {
	value string
}

func (this Token) Value() string {
	return this.value
}

func NewToken(value string) Token {
	return Token{value}
}

// Generate random token
func generateToken() Token {
	var buffer [TOKEN_SIZE / 2]byte
	rand.Read(buffer[:])
	out := fmt.Sprintf("%x", buffer)
	return NewToken(out)
}
