package data

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

type User struct {
	Handle string
	hash   []byte
}

func (this *User) setPassword(cleartext string) {
	this.hash = hash(cleartext)
}

func (this *User) Validate(cleartext string) error {
	if bytes.Equal(hash(cleartext), this.hash) {
		return nil
	}

	return errors.New("password hash mismatch")
}

func (this *User) save(overwrite bool) error {
	_, err := os.Stat(userpath(this.Handle))

	if !os.IsNotExist(err) && !overwrite {
		return errors.New("user file already exists")
	}

	return ioutil.WriteFile(
		userpath(this.Handle),
		this.hash,
		0600,
	)
}

func AddUser(handle string, password string) (*User, error) {
	account := &User{
		Handle: handle,
	}

	account.setPassword(password)

	if err := account.save(false); err != nil {
		return nil, err
	}

	if err := addPool(handle); err != nil {
		return nil, err
	}

	log.Printf("Added user \"%s\"", handle)
	return account, nil
}

func LoadUser(handle string) (*User, error) {
	hash, err := ioutil.ReadFile(userpath(handle))

	if err != nil {
		return nil, err
	}

	log.Printf("Loaded user \"%s\"", handle)

	return &User{
		Handle: handle,
		hash:   hash,
	}, nil
}

func userpath(handle string) string {
	return DATA_PATH + "/" + handle + ".user"
}

func hash(cleartext string) []byte {
	hash := sha256.New()
	hash.Write([]byte(cleartext))
	return hash.Sum(nil)
}
