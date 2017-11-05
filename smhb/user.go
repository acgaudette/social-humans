package smhb

import (
	"bytes"
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

// User data wrapper for storage
type userData struct {
	Name string
	Hash []byte
}

// User data wrapper for transmission
type userInfo struct {
	InfoHandle string
	InfoName   string
}

/* Interface implementation */

func (this *userInfo) Handle() string {
	return this.InfoHandle
}

func (this *userInfo) Name() string {
	return this.InfoName
}

// Compare two users
func (this *userInfo) Equals(other User) bool {
	return this.InfoHandle == other.Handle()
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

// Load user info raw buffer with lookup handle
func loadUserInfo(context serverContext, handle string) ([]byte, error) {
	buffer, err := ioutil.ReadFile(prefix(context, handle+".user"))

	if err != nil {
		return nil, err
	}

	// Deserialize user
	loaded := &user{}
	loaded.UnmarshalBinary(buffer)

	// Strip out hash and load into info struct
	info := &userInfo{
		InfoHandle: handle,
		InfoName:   loaded.name,
	}

	log.Printf("Loaded user \"%s\" and returned info", handle)

	return serialize(info)
}

// Deserialize raw buffer with lookup handle
func deserializeUserInfo(handle string, buffer []byte) (*userInfo, error) {
	info := &userInfo{InfoHandle: handle}
	err := deserialize(info, buffer)

	if err != nil {
		return nil, err
	}

	return info, nil
}

// Load user data with lookup handle
func getUser(context serverContext, handle string) (*user, error) {
	buffer, err := loadUser(context, handle)

	if err != nil {
		return nil, err
	}

	account := &user{handle: handle}
	err = account.UnmarshalBinary(buffer)

	if err != nil {
		return nil, err
	}

	return account, nil
}

// Load user info with lookup handle
func getUserInfo(context serverContext, handle string) (*userInfo, error) {
	buffer, err := loadUserInfo(context, handle)

	if err != nil {
		return nil, err
	}

	info, err := deserializeUserInfo(handle, buffer)

	if err != nil {
		return nil, err
	}

	return info, nil
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
