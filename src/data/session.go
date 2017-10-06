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

func (this *session) checkToken(token string) error {
	if token == this.token {
		return nil
	}

	return fmt.Errorf(
		"token mismatch for user \"%s\" session", this.handle,
	)
}

func (this *session) save() error {
	return ioutil.WriteFile(
		prefix(this.handle+".session"),
		[]byte(this.token),
		0600,
	)
}

// Set session cookie on client
func (this *session) writeToClient(out http.ResponseWriter) {
	cookie := http.Cookie{
		Name:  SESSION_NAME,
		Value: this.handle + DELM + this.token,
	}

	http.SetCookie(out, &cookie)
	log.Printf("Created new session with token \"%s\"", this.token)
}

// Generate a token and create a new session
func AddSession(out http.ResponseWriter, account *User) error {
	this := &session{
		handle: account.Handle,
		token:  generateToken(),
	}

	if err := this.save(); err != nil {
		return err
	}

	this.writeToClient(out)
	return nil
}

// Join an existing session
func JoinSession(out http.ResponseWriter, account *User) error {
	this, err := loadSession(account.Handle)

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

func GetUserFromSession(in *http.Request) (*User, error) {
	cookie, err := in.Cookie(SESSION_NAME)

	if err != nil {
		return nil, err
	}

	split := strings.Split(cookie.Value, DELM)
	s, err := loadSession(split[0])

	if err != nil {
		return nil, err
	}

	if err = s.checkToken(split[1]); err != nil {
		return nil, err
	}

	account, err := LoadUser(s.handle)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func loadSession(handle string) (*session, error) {
	file, err := os.Open(prefix(handle + ".session"))

	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	token := scanner.Text()

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	log.Printf("Loaded session with token \"%s\"", token)

	return &session{
		handle: handle,
		token:  token,
	}, nil
}
