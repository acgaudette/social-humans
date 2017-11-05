package smhb

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

// User data representation structure
type user struct {
	handle string
	name   string
	hash   []byte
}

/* Interface implementation */

func (this *user) Handle() string {
	return this.handle
}

func (this *user) Name() string {
	return this.name
}

func (this *user) Equals(other User) bool {
	return this.handle == other.Handle()
}

func (this *user) GetPath() string {
	return this.handle + ".user"
}

func (this *user) String() string {
	return "user \"" + this.handle + "\""
}

// User data wrapper for storage
type userData struct {
	Name string
	Hash []byte
}

// Test given password against user password
func (this *user) validate(cleartext string) error {
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
	cleartext string, context serverContext, access Access,
) error {
	this.setPassword(cleartext)

	// Store hash data
	if err := access.Save(this, true, context); err != nil {
		return err
	}

	log.Printf("Password updated for \"%s\"", this.handle)

	return nil
}

func (this *user) setName(
	name string, context serverContext, access Access,
) error {
	this.name = name

	if err := access.Save(this, true, context); err != nil {
		return err
	}

	log.Printf("Name updated for \"%s\"", this.handle)
	return nil
}

// Set password for user account
func (this *user) setPassword(cleartext string) {
	// Make hash
	this.hash = hash(cleartext)
}

// Add new user, given a handle
func addUser(
	handle, password, name string, context serverContext, access Access,
) (*user, error) {
	account := &user{
		handle: handle,
		name:   name,
	}

	account.setPassword(password)

	// Save, but throw error if overwriting
	if err := access.Save(account, false, context); err != nil {
		return nil, err
	}

	// Create new pool with handle
	if _, err := addPool(handle, context, access); err != nil {
		return nil, err
	}

	log.Printf("Added user \"%s\"", handle)

	return account, nil
}

// Load user data with lookup handle
func getUser(
	handle string, context serverContext, access Access,
) (*user, error) {
	account := &user{handle: handle}
	err := access.Load(account, context)

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

	return serialize(wrapper)
}

func (this *user) UnmarshalBinary(buffer []byte) error {
	wrapper := &userData{}
	err := deserialize(wrapper, buffer)

	if err != nil {
		return err
	}

	// Load wrapper into new user struct
	this.name = wrapper.Name
	this.hash = wrapper.Hash

	return nil
}
