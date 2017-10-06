package data

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// User data representation structure
type User struct {
	Handle string
	Name   string
	hash   []byte
}

// User data wrapper for serialization
type userData struct {
	Name string
	Hash []byte
}

// Test given password against user password
func (this *User) Validate(cleartext string) error {
	// Compare hashes
	if bytes.Equal(hash(cleartext), this.hash) {
		return nil
	}

	return fmt.Errorf(
		"password hash mismatch for user \"%s\"", this.Handle,
	)
}

// Set password for user account
func (this *User) setPassword(cleartext string) {
	// Make hash
	this.hash = hash(cleartext)
}

// Change password for user account
func (this *User) UpdatePassword(cleartext string) error {
	this.setPassword(cleartext)

	// Store hash data
	if err := this.save(true); err != nil {
		return err
	}

	log.Printf("Password updated for \"%s\"", this.Handle)

	return nil
}

func (this *User) SetName(name string) error {
	this.Name = name

	if err := this.save(true); err != nil {
		return err
	}

	log.Printf("Name updated for \"%s\"", this.Handle)
	return nil
}

// Write user to file
func (this *User) save(overwrite bool) error {
	_, err := os.Stat(prefix(this.Handle + ".user"))

	// Don't overwrite unless specified
	if !os.IsNotExist(err) && !overwrite {
		return fmt.Errorf(
			"data file for user \"%s\" already exists", this.Handle,
		)
	}

	// Serialize
	buffer, err := this.MarshalBinary()

	if err != nil {
		return err
	}

	// Write to file
	return ioutil.WriteFile(
		prefix(this.Handle+".user"), buffer, 0600,
	)
}

// Add new user, given a handle
func AddUser(handle, password, name string) (*User, error) {
	account := &User{
		Handle: handle,
		Name:   name,
	}

	account.setPassword(password)

	// Save, but throw error if overwriting
	if err := account.save(false); err != nil {
		return nil, err
	}

	// Create new pool with handle
	if _, err := AddPool(handle); err != nil {
		return nil, err
	}

	log.Printf("Added user \"%s\"", handle)

	return account, nil
}

// Load user data with lookup handle
func LoadUser(handle string) (*User, error) {
	buffer, err := ioutil.ReadFile(prefix(handle + ".user"))

	if err != nil {
		return nil, err
	}

	// Create user struct and deserialize
	account := &User{Handle: handle}
	err = account.UnmarshalBinary(buffer)

	if err != nil {
		return nil, err
	}

	log.Printf("Loaded user \"%s\"", handle)

	return account, nil
}

// Remove user data with lookup handle
func RemoveUser(handle string) error {
	if err := os.Remove(prefix(handle + ".user")); err != nil {
		return err
	}

	if err != removePool(handle); err != nil {
		return err
	}

	log.Printf("Deleted user \"%s\"", handle)

	return nil
}

/* Satisfy binary interfaces */

func (this *User) MarshalBinary() ([]byte, error) {
	// Create wrapper from user struct
	wrapper := &userData{
		this.Name,
		this.hash,
	}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	// Encode wrapper with gob
	if err := encoder.Encode(wrapper); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (this *User) UnmarshalBinary(buffer []byte) error {
	wrapper := userData{}

	reader := bytes.NewReader(buffer)
	decoder := gob.NewDecoder(reader)

	// Decode wrapper with gob
	if err := decoder.Decode(&wrapper); err != nil {
		return err
	}

	// Load wrapper into new user struct
	this.Name = wrapper.Name
	this.hash = wrapper.Hash

	return nil
}
