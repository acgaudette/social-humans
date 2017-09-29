package main

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

func hash(cleartext string) []byte {
	hash := sha256.New()
	hash.Write([]byte(cleartext))
	return hash.Sum(nil)
}

type statusMessage struct {
	Status string
}

type user struct {
	Handle string
	hash   []byte
}

func (this *user) setPassword(cleartext string) {
	this.hash = hash(cleartext)
}

func (this *user) validate(cleartext string) error {
	if bytes.Equal(hash(cleartext), this.hash) {
		return nil
	}

	return errors.New("password hash mismatch")
}

func (this *user) save(overwrite bool) error {
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

func addUser(handle string, password string) (*user, error) {
	account := &user{
		Handle: handle,
	}

	account.setPassword(password)

	if err := account.save(false); err != nil {
		return nil, err
	}

	log.Printf("Added user \"%s\"", handle)
	return account, nil
}

func loadUser(handle string) (*user, error) {
	hash, err := ioutil.ReadFile(userpath(handle))

	if err != nil {
		return nil, err
	}

	log.Printf("Loaded user \"%s\"", handle)

	return &user{
		Handle: handle,
		hash:   hash,
	}, nil
}

func userpath(handle string) string {
	return DATA_PATH + "/" + handle + ".user"
}
