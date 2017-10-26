package smhb

import (
	"bytes"
	"encoding/gob"
)

type userStore struct {
	Password string
	Name     string
}

type postStore struct {
	Content string
	Author  string
}

func serialize(this interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	if err := encoder.Encode(this); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func deserialize(this interface{}, buffer []byte) error {
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

func (this client) AddPost(title, content, author string) error {
	data, err := serialize(postStore{content, author})

	if err != nil {
		return err
	}

	return this.store(POST, title, data)
}
