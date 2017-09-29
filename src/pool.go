package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
)

type usermap map[string]string

type pool struct {
	handle string
	users  usermap
}

func (this *pool) add(handle string) error {
	if _, err := loadUser(handle); err != nil {
		return err
	}

	this.users[handle] = handle

	err := this.save()
	return err
}

func (this *pool) block(handle string) error {
	if handle == this.handle {
		return errors.New("attempted to delete self from pool")
	}

	if _, err := loadUser(handle); err != nil {
		return err
	}

	delete(this.users, handle)

	err := this.save()
	return err
}

func (this *pool) save() error {
	file, err := os.Create(poolpath(this.handle))

	if err != nil {
		return err
	}

	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, handle := range this.users {
		fmt.Fprintln(writer, handle)
	}

	return writer.Flush()
}

func addPool(handle string) error {
	this := &pool{
		handle: handle,
		users:  make(usermap),
	}

	this.users[handle] = handle

	if err := this.save(); err != nil {
		return err
	}

	log.Printf("Created new pool for \"%s\"", this.handle)

	return nil
}

func loadPool(handle string) (*pool, error) {
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

	log.Printf("Loaded pool for \"%s\"", handle)

	return &pool{
		handle: handle,
		users:  users,
	}, nil
}

func loadPoolAndAdd(handle string, username string) error {
	this, err := loadPool(handle)

	if err != nil {
		return err
	}

	err = this.add(username)

	log.Printf("Added %s to %s pool", username, handle)

	return err
}

func loadPoolAndBlock(handle string, username string) error {
	this, err := loadPool(handle)

	if err != nil {
		return err
	}

	err = this.block(username)

	log.Printf("Blocked %s from %s pool", username, handle)

	return err
}

func poolpath(handle string) string {
	return DATA_PATH + "/" + handle + ".pool"
}
