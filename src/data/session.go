package data

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type session struct {
	handle string
	token  string
}

// Compare session token with input token
func (this *session) checkToken(token string) error {
	if token == this.token {
		return nil
	}

	return fmt.Errorf(
		"token mismatch for user \"%s\" session", this.handle,
	)
}

// Set session cookie on the client
func (this *session) writeToClient(out http.ResponseWriter) {
	cookie := http.Cookie{
		Name:  SESSION_NAME,
		Value: this.handle + DELM + this.token,
	}

	http.SetCookie(out, &cookie)

	log.Printf("Created new session with token \"%s\"", this.token)
}

// Write session token to file
func (this *session) save() error {
	return ioutil.WriteFile(
		prefix(this.handle+".session"),
		[]byte(this.token),
		0600,
	)
}

// Load session with lookup handle
func loadSession(handle string) (*session, error) {
	// Open file, if it exists
	file, err := os.Open(prefix(handle + ".session"))

	if err != nil {
		return nil, err
	}

	defer file.Close()

	// Use scanner--we may read more lines in the future
	scanner := bufio.NewScanner(file)
	scanner.Scan()

	// Get token from file
	token := scanner.Text()

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	log.Printf("Loaded session with token \"%s\"", token)

	// Build new session structure
	return &session{
		handle: handle,
		token:  token,
	}, nil
}

// Generate a token and create a new session
func AddSession(out http.ResponseWriter, account User) error {
	this := &session{
		handle: account.Handle(),
		token:  generateToken(),
	}

	if err := this.save(); err != nil {
		return err
	}

	this.writeToClient(out)
	return nil
}

// Join an existing session
func JoinSession(out http.ResponseWriter, account User) error {
	// Attempt to load user session
	this, err := loadSession(account.Handle())

	// Add new session if existing session was not found
	if err != nil {
		log.Printf("%s", err)
		return AddSession(out, account)
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
func GetUserFromSession(in *http.Request) (*user, error) {
	// Get session cookie
	cookie, err := in.Cookie(SESSION_NAME)

	if err != nil {
		return nil, err
	}

	// Get handle from cookie and load session
	split := strings.Split(cookie.Value, DELM)
	this, err := loadSession(split[0])

	if err != nil {
		return nil, err
	}

	// Compare token from cookie with token from loaded session
	if err = this.checkToken(split[1]); err != nil {
		return nil, err
	}

	// Load user from from loaded session
	account, err := LoadUser(this.handle)

	if err != nil {
		return nil, err
	}

	return account, nil
}
