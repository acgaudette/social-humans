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

	return fmt.Errorf(
		"token mismatch:\nremote: %s (%d)\nlocal:  %s (%d)",
		token.value, len(token.value),
		this.value, len(this.value),
	)
}

// Write token to file
func (this Token) save(context serverContext, handle string) error {
	return ioutil.WriteFile(
		prefix(context, handle+".key"),
		[]byte(this.value),
		0600,
	)
}

func NewToken(value string) Token {
	return Token{value}
}

// Authenticate a token, handle pair (e.g. in a server request)
func authenticate(
	token Token,
	handle string,
	context serverContext,
) (error, bool) {
	key, err := getToken(context, handle)

	if err != nil {
		return err, false
	}

	err = key.compare(token)

	if err == nil {
		return nil, true
	}

	return err, false
}

func addToken(context serverContext, handle string) (*Token, error) {
	this := generateToken()

	if err := this.save(context, handle); err != nil {
		return nil, err
	}

	log.Printf("Added token for user \"%s\"", handle)

	return &this, nil
}

func getToken(context serverContext, handle string) (*Token, error) {
	// Open file, if it exists
	file, err := os.Open(prefix(context, handle+".key"))

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var buffer [TOKEN_SIZE + 1]byte
	_, err = file.Read(buffer[:])

	if err != nil {
		return nil, err
	}

	value := string(buffer[:TOKEN_SIZE])

	log.Printf("Loaded token for user \"%s\"", handle)

	return &Token{value}, nil
}

// Generate random token
func generateToken() Token {
	var buffer [TOKEN_SIZE / 2]byte
	rand.Read(buffer[:])
	out := fmt.Sprintf("%x", buffer)
	return NewToken(out)
}
