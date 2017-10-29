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

// Generate random string
func generateToken() Token {
	buffer := make([]byte, TOKEN_LENGTH / 2)
	rand.Read(buffer)
	out := fmt.Sprintf("%x", buffer)
	return NewToken(out)
}
