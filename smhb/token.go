package smhb

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Token struct {
	value string
}

func (this Token) Value() string {
	return this.value
}

// Compare token with input token
func (this Token) compare(token Token) error {
	if token.value == this.value {
		return nil
	}

	return fmt.Errorf("token mismatch")
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
