package smhb

import (
	"bytes"
	"encoding/gob"
)

type userStore struct {
	Password string
	Name     string
}

func (this userStore) serialize() ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	if err := encoder.Encode(this); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (this *userStore) deserialize(buffer []byte) error {
	reader := bytes.NewReader(buffer)
	decoder := gob.NewDecoder(reader)

	if err := decoder.Decode(&this); err != nil {
		return err
	}

	return nil
}

func (this client) AddUser(handle, password, name string) (User, error) {
	data, err := userStore{password, name}.serialize()

	if err != nil {
		return nil, err
	}

	err = this.store(USER, handle, data)

	if err != nil {
		return nil, err
	}

	return this.GetUser(handle)
}
