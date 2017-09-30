package data

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

type User struct {
	Handle string
	Name   string
	hash   []byte
}

type userData struct {
	Name string
	Hash []byte
}

func (this *User) MarshalBinary() ([]byte, error) {
	wrapper := &userData{
		this.Name,
		this.hash,
	}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	if err := encoder.Encode(wrapper); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (this *User) UnmarshalBinary(buffer []byte) error {
	wrapper := userData{}

	reader := bytes.NewReader(buffer)
	decoder := gob.NewDecoder(reader)

	if err := decoder.Decode(&wrapper); err != nil {
		return err
	}

	this.Name = wrapper.Name
	this.hash = wrapper.Hash

	return nil
}

func (this *User) Validate(cleartext string) error {
	if bytes.Equal(hash(cleartext), this.hash) {
		return nil
	}

	return errors.New("password hash mismatch")
}

func (this *User) SetPassword(cleartext string) error {
	this.hash = hash(cleartext)
	return this.save(true)
}

func (this *User) SetName(name string) error {
	this.Name = name
	return this.save(true)
}

func (this *User) save(overwrite bool) error {
	_, err := os.Stat(path(this.Handle, "user"))

	if !os.IsNotExist(err) && !overwrite {
		return errors.New("user file already exists")
	}

	buffer, err := this.MarshalBinary()

	if err != nil {
		return err
	}

	return ioutil.WriteFile(
		path(this.Handle, "user"), buffer, 0600,
	)
}

func AddUser(handle, password, name string) (*User, error) {
	account := &User{
		Handle: handle,
		Name:   name,
	}

	account.SetPassword(password)

	if err := account.save(false); err != nil {
		return nil, err
	}

	if _, err := AddPool(handle); err != nil {
		return nil, err
	}

	log.Printf("Added user \"%s\"", handle)
	return account, nil
}

func LoadUser(handle string) (*User, error) {
	buffer, err := ioutil.ReadFile(path(handle, "user"))

	if err != nil {
		return nil, err
	}

	loaded := &User{Handle: handle}
	err = loaded.UnmarshalBinary(buffer)

	if err != nil {
		return nil, err
	}

	log.Printf("Loaded user \"%s\"", handle)

	return loaded, nil
}

func RemoveUser(handle string) error {
	return os.Remove(path(handle, "user"))
}

func hash(cleartext string) []byte {
	hash := sha256.New()
	hash.Write([]byte(cleartext))
	return hash.Sum(nil)
}
