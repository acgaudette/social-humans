package data

import (
	"../../smhb"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type session struct {
	token string
	key   smhb.Token
}

/* Interface implementation */

func (this *session) GetPath() string {
	return this.key.Handle() + ".session"
}

func (this *session) String() string {
	return "session for user \"" + this.key.Handle() + "\""
}

func (this *session) MarshalBinary() ([]byte, error) {
	return append([]byte(this.token), []byte(this.key.Value())...), nil
}

func (this *session) UnmarshalBinary(buffer []byte) error {
	reader := bytes.NewBuffer(buffer)

	var tokenBuffer [TOKEN_SIZE]byte
	var keyBuffer [smhb.TOKEN_SIZE + 1]byte

	_, err := reader.Read(tokenBuffer[:])

	if err != nil {
		return err
	}

	_, err = reader.Read(keyBuffer[:])

	if err != nil {
		return err
	}

	this.token = string(tokenBuffer[:TOKEN_SIZE])
	this.key = smhb.NewToken(string(keyBuffer[:smhb.TOKEN_SIZE]), "")

	return nil
}

// Compare session token with input token
func (this *session) checkToken(token string) error {
	if token == this.token {
		return nil
	}

	return fmt.Errorf(
		"token mismatch for user \"%s\" session", this.key.Handle(),
	)
}

// Set session cookie on the client
func (this *session) writeToClient(out http.ResponseWriter) {
	cookie := http.Cookie{
		Name:  SESSION_NAME,
		Value: this.key.Handle() + DELM + this.token,
	}

	http.SetCookie(out, &cookie)

	log.Printf("Created new session with token \"%s\"", this.token)
}

// Generate a token and create a new session
func AddSession(
	out http.ResponseWriter, handle string, token *smhb.Token,
) error {
	this := &session{
		token: generateToken(),
		key:   *token,
	}

	if err := access.Save(this, true, accessContext); err != nil {
		return err
	}

	this.writeToClient(out)
	return nil
}

// Join an existing session
func JoinSession(out http.ResponseWriter, handle, password string) error {
	// Add new session if existing session was not found
	rewrite := func(err error) error {
		log.Printf("%s", err)

		// Get new key from backend
		key, err := Backend.GetToken(handle, password)

		if err != nil {
			return err
		}

		return AddSession(out, handle, key)
	}

	// Attempt to load user session
	err := Backend.Validate(handle, password)

	if err != nil {
		return rewrite(err)
	}

	this := &session{
		key: smhb.NewToken("", handle),
	}

	err = access.Load(this, accessContext)

	if err != nil {
		return rewrite(err)
	}

	// Rebuild key
	this.key = smhb.NewToken(this.key.Value(), handle)

	// Check that token exists
	err = Backend.CheckToken(handle)

	if err != nil {
		return rewrite(err)
	}

	// Overwrite the existing cookie
	// Worst case scenario, it writes the same cookie twice
	this.writeToClient(out)
	return nil
}

// Clear a session on the client (but not the server)
func ClearSession(out http.ResponseWriter) {
	cookie := http.Cookie{
		Name:    SESSION_NAME,
		Value:   "",
		Expires: time.Now().Add(-time.Minute),
	}

	http.SetCookie(out, &cookie)

	log.Printf("Cleared session")
}

// Load a user structure from the current session
func GetUserFromSession(in *http.Request) (smhb.User, *smhb.Token, error) {
	// Get session cookie
	cookie, err := in.Cookie(SESSION_NAME)

	if err != nil {
		return nil, nil, err
	}

	// Get handle from cookie and load session
	split := strings.Split(cookie.Value, DELM)
	this := &session{
		key: smhb.NewToken("", split[0]),
	}

	err = access.Load(this, accessContext)

	if err != nil {
		return nil, nil, err
	}

	// Rebuild key
	this.key = smhb.NewToken(this.key.Value(), split[0])

	// Compare token from cookie with token from loaded session
	if err = this.checkToken(split[1]); err != nil {
		return nil, nil, err
	}

	// Load user from from loaded session
	account, err := Backend.GetUser(this.key.Handle())

	if err != nil {
		return nil, nil, err
	}

	return account, &this.key, nil
}
