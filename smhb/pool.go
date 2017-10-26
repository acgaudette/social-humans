package smhb

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Add handle (idempotent)
func (this userPool) add(handle string) {
	this[handle] = handle
}

// Remove handle
func (this userPool) remove(handle string) {
	delete(this, handle)
}

func newUserPool() userPool {
	return make(userPool)
}

// Pool data representation structure
type pool struct {
	handle string
	users  userPool
}

/* Interface implementation getters */

func (this *pool) Handle() string {
	return this.handle
}

func (this *pool) Users() userPool {
	return this.users
}

// Pool data wrapper for serialization
type poolData struct {
	Users userPool
}

// Add a user to the pool, given a handle
func (this *pool) add(handle string) error {
	// Confirm that the given user exists
	if _, err := getUser(handle); err != nil {
		return err
	}

	// Add handle to user pool
	this.users.add(handle)

	// Update data
	err := this.save()

	if err == nil {
		log.Printf("Added \"%s\" to \"%s\" pool", handle, this.Handle())
	}

	return err
}

// Remove a user from the pool, given a handle
func (this *pool) block(handle string) error {
	// Ignore self
	if handle == this.handle {
		return fmt.Errorf(
			"user \"%s\" attempted to delete self from pool", this.handle,
		)
	}

	// Confirm that the given user exists
	if _, err := getUser(handle); err != nil {
		return err
	}

	// Remove handle from user pool
	this.users.remove(handle)

	// Update data
	err := this.save()

	if err == nil {
		log.Printf("Blocked \"%s\" from \"%s\" pool", handle, this.handle)
	}

	return err
}

// Write pool to file
func (this *pool) save() error {
	// Clean pool before saving
	this.clean()

	// Serialize
	buffer, err := this.MarshalBinary()

	if err != nil {
		return err
	}

	return ioutil.WriteFile(
		prefix(this.handle+".pool"), buffer, 0600,
	)
}

// Remove users from the pool that no longer exist
func (this *pool) clean() {
	// Iterate through handles in user pool
	for _, handle := range this.users {
		// If user cannot be loaded, remove handle
		if _, err := getUser(handle); err != nil {
			this.users.remove(handle)
		}
	}
}

// Add new pool, given a user handle
func addPool(handle string) (*pool, error) {
	this := &pool{
		handle: handle,
		users:  newUserPool(),
	}

	// Add self to pool
	this.users.add(handle)

	// Update data
	if err := this.save(); err != nil {
		return nil, err
	}

	log.Printf("Created new pool for user \"%s\"", this.handle)

	return this, nil
}

// Load pool raw buffer with lookup handle
func loadPool(handle string) ([]byte, error) {
	buffer, err := ioutil.ReadFile(prefix(handle + ".pool"))

	if err != nil {
		return nil, err
	}

	log.Printf("Loaded pool for user \"%s\"", handle)

	return buffer, nil
}

// Deserialize raw buffer with lookup handle
func deserializePool(handle string, buffer []byte) (*pool, error) {
	// Create pool struct and deserialize
	loaded := &pool{handle: handle}
	err := loaded.UnmarshalBinary(buffer)

	if err != nil {
		return nil, err
	}

	return loaded, nil
}

// Load pool data with lookup handle
func getPool(handle string) (*pool, error) {
	buffer, err := loadPool(handle)

	if err != nil {
		return nil, err
	}

	loaded, err := deserializePool(handle, buffer)

	if err != nil {
		return nil, err
	}

	return loaded, nil
}

// Remove pool data with lookup handle
func removePool(handle string) error {
	if err := os.Remove(prefix(handle + ".pool")); err != nil {
		return err
	}

	log.Printf("Deleted pool for user \"%s\"", handle)

	return nil
}

/* Satisfy binary interfaces */

func (this *pool) MarshalBinary() ([]byte, error) {
	// Create wrapper from pool struct
	wrapper := &poolData{this.users}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	// Encode wrapper with gob
	if err := encoder.Encode(wrapper); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (this *pool) UnmarshalBinary(buffer []byte) error {
	wrapper := poolData{}

	reader := bytes.NewReader(buffer)
	decoder := gob.NewDecoder(reader)

	// Decode wrapper with gob
	if err := decoder.Decode(&wrapper); err != nil {
		return err
	}

	// Load wrapper into new pool struct
	this.users = wrapper.Users

	return nil
}
