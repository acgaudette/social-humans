package data

import (
	"bytes"
	"encoding/gob"
	"fmt"
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

func (this *Pool) Add(handle string) error {
	if _, err := LoadUser(handle); err != nil {
		return err
	}

	this.Users[handle] = handle
	err := this.save()

	if err == nil {
		log.Printf("Added \"%s\" to \"%s\" pool", handle, this.Handle)
	}

	return err
}

func (this *Pool) Block(handle string) error {
	if handle == this.Handle {
		return fmt.Errorf(
			"user \"%s\" attempted to delete self from pool", this.Handle,
		)
	}

	if _, err := LoadUser(handle); err != nil {
		return err
	}

	delete(this.Users, handle)
	err := this.save()

	if err == nil {
		log.Printf("Blocked \"%s\" from \"%s\" pool", handle, this.Handle)
	}

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

func (this *Pool) clean() {
	for _, handle := range this.Users {
		_, err := LoadUser(handle)

		if err != nil {
			delete(this.Users, handle)
		}
	}
}

func AddPool(handle string) (*Pool, error) {
	this := &Pool{
		Handle: handle,
		Users:  make(userPool),
	}

	this.Users[handle] = handle

	if err := this.save(); err != nil {
		return nil, err
	}

	log.Printf("Created new pool for user \"%s\"", this.Handle)

	return this, nil
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

	loaded.clean()

	log.Printf("Loaded pool for user \"%s\"", handle)

	return loaded, nil
}
