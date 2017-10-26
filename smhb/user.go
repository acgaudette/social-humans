package smhb

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// User data representation structure
type user struct {
	handle string
	name   string
	hash   []byte
}

/* Interface implementation getters */

func (this *user) Handle() string {
	return this.handle
}

func (this *user) Name() string {
	return this.name
}

// User data wrapper for serialization
type userData struct {
	Name string
	Hash []byte
}

// Test given password against user password
func (this *user) Validate(cleartext string) error {
	// Compare hashes
	if bytes.Equal(hash(cleartext), this.hash) {
		return nil
	}

	return fmt.Errorf(
		"password hash mismatch for user \"%s\"", this.handle,
	)
}

// Change password for user account
func (this *user) updatePassword(
	context serverContext, cleartext string,
) error {
	this.setPassword(cleartext)

	// Store hash data
	if err := this.save(context, true); err != nil {
		return err
	}

	log.Printf("Password updated for \"%s\"", this.handle)

	return nil
}

func (this *user) setName(context serverContext, name string) error {
	this.name = name

	if err := this.save(context, true); err != nil {
		return err
	}

	log.Printf("Name updated for \"%s\"", this.handle)
	return nil
}

// Compare two users
func (this *user) Equals(other User) bool {
	return this.handle == other.Handle()
}

// Set password for user account
func (this *user) setPassword(cleartext string) {
	// Make hash
	this.hash = hash(cleartext)
}

// Write user to file
func (this *user) save(context serverContext, overwrite bool) error {
	_, err := os.Stat(prefix(context, this.handle+".user"))

	// Don't overwrite unless specified
	if !os.IsNotExist(err) && !overwrite {
		return fmt.Errorf(
			"data file for user \"%s\" already exists", this.handle,
		)
	}

	// Serialize
	buffer, err := this.MarshalBinary()

	if err != nil {
		return err
	}

	// Write to file
	return ioutil.WriteFile(
		prefix(context, this.handle+".user"), buffer, 0600,
	)
}

// Add new user, given a handle
func addUser(
	context serverContext, handle, password, name string,
) (*user, error) {
	account := &user{
		handle: handle,
		name:   name,
	}

	account.setPassword(password)

	// Save, but throw error if overwriting
	if err := account.save(context, false); err != nil {
		return nil, err
	}

	// Create new pool with handle
	if _, err := addPool(context, handle); err != nil {
		return nil, err
	}

	log.Printf("Added user \"%s\"", handle)

	return account, nil
}

// Load user raw buffer with lookup handle
func loadUser(context serverContext, handle string) ([]byte, error) {
	buffer, err := ioutil.ReadFile(prefix(context, handle+".user"))

	if err != nil {
		return nil, err
	}

	log.Printf("Loaded user \"%s\"", handle)

	return buffer, nil
}

// Deserialize raw buffer with lookup handle
func deserializeUser(handle string, buffer []byte) (*user, error) {
	account := &user{handle: handle}
	err := account.UnmarshalBinary(buffer)

	if err != nil {
		return nil, err
	}

	return account, nil
}

// Load user data with lookup handle
func getUser(context serverContext, handle string) (*user, error) {
	buffer, err := loadUser(context, handle)

	if err != nil {
		return nil, err
	}

	account, err := deserializeUser(handle, buffer)

	if err != nil {
		return nil, err
	}

	return account, nil
}

// Remove user data with lookup handle
func removeUser(context serverContext, handle string) error {
	if err := os.Remove(prefix(context, handle+".user")); err != nil {
		return err
	}

	if err := removePool(context, handle); err != nil {
		return err
	}

	if err := removePostsByAuthor(context, handle); err != nil {
		return err
	}

	log.Printf("Deleted user \"%s\"", handle)

	return nil
}

/* Satisfy binary interfaces */

func (this *user) MarshalBinary() ([]byte, error) {
	// Create wrapper from user struct
	wrapper := &userData{
		this.name,
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

func (this *user) UnmarshalBinary(buffer []byte) error {
	wrapper := userData{}

	reader := bytes.NewReader(buffer)
	decoder := gob.NewDecoder(reader)

	// Decode wrapper with gob
	if err := decoder.Decode(&wrapper); err != nil {
		return err
	}

	// Load wrapper into new user struct
	this.name = wrapper.Name
	this.hash = wrapper.Hash

	return nil
}
