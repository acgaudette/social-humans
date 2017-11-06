package smhb

import (
	"fmt"
	"log"
)

// Store user pool as a set of handles
type userPool map[string]string

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

/* Interface implementation */

func (this pool) Handle() string {
	return this.handle
}

func (this pool) Users() userPool {
	return this.users
}

func (this pool) GetPath() string {
	return this.handle + ".pool"
}

func (this pool) String() string {
	return "\"" + this.handle + "\" pool"
}

// Pool data wrapper for serialization
type poolData struct {
	Users userPool
}

// Add a user to the pool, given a handle
func (this pool) add(
	handle string, context serverContext, access Access,
) error {
	// Confirm that the given user exists
	if _, err := getRawUserInfo(handle, context, access); err != nil {
		return err
	}

	// Add handle to user pool
	this.users.add(handle)

	// Update data
	err := access.Save(this, true, context)

	if err == nil {
		log.Printf("Added \"%s\" to %s", handle, this)
	}

	return err
}

// Remove a user from the pool, given a handle
func (this pool) block(
	handle string, context serverContext, access Access,
) error {
	// Ignore self
	if handle == this.handle {
		return fmt.Errorf(
			"user \"%s\" attempted to delete self from pool", this.handle,
		)
	}

	// Confirm that the given user exists
	if _, err := getRawUserInfo(handle, context, access); err != nil {
		return err
	}

	// Remove handle from user pool
	this.users.remove(handle)

	// Update data
	err := access.Save(this, true, context)

	if err == nil {
		log.Printf("Blocked \"%s\" from %s", handle, this)
	}

	return err
}

// Remove users from the pool that no longer exist
func (this pool) clean(context serverContext, access Access) (modified bool) {
	// Iterate through handles in user pool
	for _, handle := range this.users {
		// If user cannot be loaded, remove handle
		if _, err := getRawUserInfo(handle, context, access); err != nil {
			log.Printf("Cleaned \"%s\" from %s", handle, this)
			this.users.remove(handle)
			modified = true
		}
	}
	return
}

// Add new pool, given a user handle
func addPool(
	handle string, context serverContext, access Access,
) (*pool, error) {
	this := &pool{
		handle: handle,
		users:  newUserPool(),
	}

	// Add self to pool
	this.users.add(handle)

	// Update data
	if err := access.Save(this, true, context); err != nil {
		return nil, err
	}

	log.Printf("Created new pool for user \"%s\"", this.handle)

	return this, nil
}

// Load pool raw buffer with lookup handle
func getRawPool(
	handle string, context serverContext, access Access,
) ([]byte, error) {
	loaded := pool{handle: handle}
	buffer, err := access.LoadRaw(loaded, context)

	if err != nil {
		return nil, err
	}

	// Deserialize (unfortunately has to be done) for cleaning
	err = loaded.UnmarshalBinary(buffer)

	if err != nil {
		return nil, err
	}

	// If modified after clean, save and reserialize
	if loaded.clean(context, access) {
		err = access.Save(loaded, true, context)

		if err != nil {
			return nil, err
		}

		// Reserialize
		buffer, err = loaded.MarshalBinary()

		if err != nil {
			return nil, err
		}
	}

	return buffer, nil
}

// Load pool data with lookup handle
func getPool(
	handle string, context serverContext, access Access,
) (*pool, error) {
	buffer, err := getRawPool(handle, context, access)

	if err != nil {
		return nil, err
	}

	loaded := &pool{handle: handle}
	err = loaded.UnmarshalBinary(buffer)

	if err != nil {
		return nil, err
	}

	return loaded, nil
}

// Remove pool data with lookup handle
func removePool(
	handle string, context serverContext, access Access,
) error {
	target := pool{handle: handle}
	err := access.Remove(target, context)

	if err != nil {
		return err
	}

	return nil
}

/* Satisfy binary interfaces */

func (this pool) MarshalBinary() ([]byte, error) {
	// Create wrapper from pool struct and serialize
	wrapper := &poolData{this.users}
	return serialize(wrapper)
}

func (this pool) UnmarshalBinary(buffer []byte) error {
	wrapper := &poolData{}
	err := deserialize(wrapper, buffer)

	if err != nil {
		return err
	}

	// Load wrapper into new pool struct
	this.users = wrapper.Users

	return nil
}
