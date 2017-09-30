package data

import (
	"bufio"
	"crypto/rand"
	"errors"
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

	return errors.New("token mismatch")
}

func (this *session) save() error {
	return ioutil.WriteFile(
		sessionpath(this.handle),
		[]byte(this.token),
		0600,
	)
}

func AddSession(writer http.ResponseWriter, account *user) error {
	this := &session{
		handle: account.Handle,
		token:  generateToken(),
	}

	if err := this.save(); err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:  SESSION_NAME,
		Value: this.handle + DELM + this.token,
	}

	http.SetCookie(writer, &cookie)
	log.Printf("Created new session with token \"%s\"", this.token)

	return nil
}

func ClearSession(writer http.ResponseWriter) {
	cookie := http.Cookie{
		Name:    SESSION_NAME,
		Value:   "",
		Expires: time.Now().Add(-time.Minute),
	}

	http.SetCookie(writer, &cookie)
	log.Printf("Cleared session")
}

func GetUserFromSession(request *http.Request) (*user, error) {
	cookie, err := request.Cookie(SESSION_NAME)

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
	file, err := os.Open(sessionpath(handle))

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

func sessionpath(handle string) string {
	return DATA_PATH + "/" + handle + ".session"
}

func generateToken() string {
	buffer := make([]byte, 32)
	rand.Read(buffer)
	return fmt.Sprintf("%x", buffer)
}
