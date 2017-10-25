package smhb

import (
	"bytes"
	"encoding/gob"
)

type store interface{}

type userStore struct {
	Password string
	Name     string
}

func serialize(this store) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	if err := encoder.Encode(this); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func deserialize(this store, buffer []byte) error {
	reader := bytes.NewReader(buffer)
	decoder := gob.NewDecoder(reader)

	if err := decoder.Decode(this); err != nil {
		return err
	}

	return nil
}

func (this client) AddUser(handle, password, name string) (User, error) {
	data, err := serialize(userStore{password, name})

	if err != nil {
		return nil, err
	}

	err = this.store(USER, handle, data)

	if err != nil {
		return nil, err
	}

	return this.GetUser(handle)
}
