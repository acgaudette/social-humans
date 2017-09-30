package data

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
)

type usermap map[string]string

type Pool struct {
	Handle string
	Users  usermap
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
	file, err := os.Create(poolpath(this.Handle))

	if err != nil {
		return err
	}

	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, handle := range this.Users {
		fmt.Fprintln(writer, handle)
	}

	return writer.Flush()
}

func LoadPool(handle string) (*Pool, error) {
	file, err := os.Open(poolpath(handle))

	if err != nil {
		return nil, err
	}

	defer file.Close()

	users := make(usermap)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		username := scanner.Text()
		users[username] = username
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	log.Printf("Loaded pool for user \"%s\"", handle)

	return &Pool{
		Handle: handle,
		Users:  users,
	}, nil
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
		Users:  make(usermap),
	}

	this.Users[handle] = handle

	if err := this.save(); err != nil {
		return err
	}

	log.Printf("Created new pool for user \"%s\"", this.Handle)

	return nil
}

func poolpath(handle string) string {
	return DATA_PATH + "/" + handle + ".pool"
}
