package data

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io/ioutil"
	"log"
)

type userPool map[string]string

type Pool struct {
	Handle string
	Users  userPool
}

type poolData struct {
	Users userPool
}

func (this *Pool) MarshalBinary() ([]byte, error) {
	wrapper := &poolData{this.Users}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	if err := encoder.Encode(wrapper); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (this *Pool) UnmarshalBinary(buffer []byte) error {
	wrapper := poolData{}

	reader := bytes.NewReader(buffer)
	decoder := gob.NewDecoder(reader)

	if err := decoder.Decode(&wrapper); err != nil {
		return err
	}

	this.Users = wrapper.Users
	return nil
}

func (this *Pool) add(handle string) error {
	if _, err := LoadUser(handle); err != nil {
		return err
	}

	this.Users[handle] = handle

	err := this.save()
	return err
}

func (this *Pool) block(handle string) error {
	if handle == this.Handle {
		return errors.New("attempted to delete self from pool")
	}

	if _, err := LoadUser(handle); err != nil {
		return err
	}

	delete(this.Users, handle)

	err := this.save()
	return err
}

func (this *Pool) save() error {
	buffer, err := this.MarshalBinary()

	if err != nil {
		return err
	}

	return ioutil.WriteFile(
		path(this.Handle, "pool"), buffer, 0600,
	)
}

func LoadPool(handle string) (*Pool, error) {
	buffer, err := ioutil.ReadFile(path(handle, "pool"))

	if err != nil {
		return nil, err
	}

	loaded := &Pool{Handle: handle}
	err = loaded.UnmarshalBinary(buffer)

	if err != nil {
		return nil, err
	}

	log.Printf("Loaded pool for user \"%s\"", handle)

	return loaded, nil
}

func LoadPoolAndAdd(handle string, username string) error {
	this, err := LoadPool(handle)

	if err != nil {
		return err
	}

	err = this.add(username)

	if err == nil {
		log.Printf("Added %s to %s pool", username, handle)
	}

	return err
}

func LoadPoolAndBlock(handle string, username string) error {
	this, err := LoadPool(handle)

	if err != nil {
		return err
	}

	err = this.block(username)

	if err == nil {
		log.Printf("Blocked %s from %s pool", username, handle)
	}

	return err
}

func addPool(handle string) error {
	this := &Pool{
		Handle: handle,
		Users:  make(userPool),
	}

	this.Users[handle] = handle

	if err := this.save(); err != nil {
		return err
	}

	log.Printf("Created new pool for user \"%s\"", this.Handle)

	return nil
}
