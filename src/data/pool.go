package data

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
)

// Store user pool as a set of handles
type pool map[string]string

// Add handle (idempotent)
func (this pool) add(handle string) {
	this[handle] = handle
}

// Remove handle
func (this pool) remove(handle string) {
	delete(this, handle)
}

func newPool() pool {
	return make(pool)
}

// Pool data representation structure
type Pool struct {
	Handle string
	Users  pool
}

// Pool data wrapper for serialization
type poolData struct {
	Users pool
}

// Add a user to the pool, given a handle
func (this *Pool) Add(handle string) error {
	// Confirm that the given user exists
	if _, err := LoadUser(handle); err != nil {
		return err
	}

	// Add handle to user pool
	this.Users.add(handle)

	// Update data
	err := this.save()

	if err == nil {
		log.Printf("Added \"%s\" to \"%s\" pool", handle, this.Handle)
	}

	return err
}

// Remove a user from the pool, given a handle
func (this *Pool) Block(handle string) error {
	// Ignore self
	if handle == this.Handle {
		return fmt.Errorf(
			"user \"%s\" attempted to delete self from pool", this.Handle,
		)
	}

	// Confirm that the given user exists
	if _, err := LoadUser(handle); err != nil {
		return err
	}

	// Remove handle from user pool
	this.Users.remove(handle)

	// Update data
	err := this.save()

	if err == nil {
		log.Printf("Blocked \"%s\" from \"%s\" pool", handle, this.Handle)
	}

	return err
}

// Write pool to file
func (this *Pool) save() error {
	// Serialize
	buffer, err := this.MarshalBinary()

	if err != nil {
		return err
	}

	return ioutil.WriteFile(
		prefix(this.Handle+".pool"), buffer, 0600,
	)
}

// Remove users from the pool that no longer exist
func (this *Pool) clean() {
	// Iterate through handles in user pool
	for _, handle := range this.Users {
		// If user cannot be loaded, remove handle
		if _, err := LoadUser(handle); err != nil {
			this.Users.remove(handle)
		}
	}
}

// Add new pool, given a user handle
func AddPool(handle string) (*Pool, error) {
	this := &Pool{
		Handle: handle,
		Users:  newPool(),
	}

	// Add self to pool
	this.Users.add(handle)

	// Update data
	if err := this.save(); err != nil {
		return nil, err
	}

	log.Printf("Created new pool for user \"%s\"", this.Handle)

	return this, nil
}

// Load pool data with lookup handle
func LoadPool(handle string) (*Pool, error) {
	buffer, err := ioutil.ReadFile(prefix(handle + ".pool"))

	if err != nil {
		return nil, err
	}

	// Create pool struct and deserialize
	loaded := &Pool{Handle: handle}
	err = loaded.UnmarshalBinary(buffer)

	if err != nil {
		return nil, err
	}

	// Clean pool after loading
	loaded.clean()

	log.Printf(
		"Loaded pool for user \"%s\" (%v users)", handle, len(loaded.Users),
	)

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

func (this *Pool) MarshalBinary() ([]byte, error) {
	// Create wrapper from pool struct
	wrapper := &poolData{this.Users}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	// Encode wrapper with gob
	if err := encoder.Encode(wrapper); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (this *Pool) UnmarshalBinary(buffer []byte) error {
	wrapper := poolData{}

	reader := bytes.NewReader(buffer)
	decoder := gob.NewDecoder(reader)

	// Decode wrapper with gob
	if err := decoder.Decode(&wrapper); err != nil {
		return err
	}

	// Load wrapper into new pool struct
	this.Users = wrapper.Users

	return nil
}
