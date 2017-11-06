package smhb

import (
	"crypto/rand"
	"fmt"
)

type Token struct {
	value  string
	handle string
}

func (this Token) Value() string {
	return this.value
}

func (this Token) Handle() string {
	return this.handle
}

/* Interface implementation */

func (this Token) GetPath() string {
	return this.handle + ".key"
}

func (this Token) String() string {
	return "token for user \"" + this.handle + "\""
}

func (this *Token) MarshalBinary() ([]byte, error) {
	return []byte(this.value), nil
}

func (this *Token) UnmarshalBinary(buffer []byte) error {
	this.value = string(buffer[:TOKEN_SIZE])
	return nil
}

// Compare token with input token
func (this Token) compare(token Token) error {
	if token.value == this.value && token.handle == this.handle {
		return nil
	}

	return fmt.Errorf(
		"token mismatch:\nremote: %s (%d) \"%s\"\nlocal:  %s (%d) \"%s\"",
		token.value, len(token.value), token.handle,
		this.value, len(this.value), this.handle,
	)
}

func NewToken(value, handle string) Token {
	return Token{
		value,
		handle,
	}
}

// Authenticate a token, handle pair (e.g. in a server request)
func authenticate(
	token Token,
	context serverContext,
	access Access,
) (error, bool) {
	key, err := getToken(token.handle, context, access)

	if err != nil {
		return err, false
	}

	err = key.compare(token)

	if err == nil {
		return nil, true
	}

	return err, false
}

func addToken(
	handle string, context serverContext, access Access,
) (*Token, error) {
	this := &Token{
		value: generateTokenValue(),
		handle: handle,
	}

	if err := access.Save(this, true, context); err != nil {
		return nil, err
	}

	return this, nil
}

func getToken(
	handle string, context serverContext, access Access,
) (*Token, error) {
	token := &Token{handle: handle}
	err := access.Load(token, context)

	if err != nil {
		return nil, err
	}

	return token, nil
}

// Generate random token value
func generateTokenValue() string {
	var buffer [TOKEN_SIZE / 2]byte
	rand.Read(buffer[:])
	return fmt.Sprintf("%x", buffer)
}
