package smhb

import (
	"encoding"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Accessable interface {
	GetPath() string
	String() string
}

type Storeable interface {
	Accessable
	encoding.BinaryMarshaler
}

type Loadable interface {
	Accessable
	encoding.BinaryUnmarshaler
}

type Access interface {
	Save(Storeable, bool, ServerContext) error
	SaveWithDir(Storeable, string, bool, ServerContext) error
	Load(Loadable, ServerContext) error
	LoadRaw(Loadable, ServerContext) ([]byte, error)
	Remove(Accessable, ServerContext) error
	RemoveDir(string, ServerContext) error
}

type FileAccess struct{}

func (this FileAccess) Save(
	target Storeable, overwrite bool, context ServerContext,
) error {
	_, err := os.Stat(prefix(context, target.GetPath()))

	// Don't overwrite unless specified
	if !os.IsNotExist(err) && !overwrite {
		return fmt.Errorf("data file for %s already exists", target)
	}

	// Serialize
	buffer, err := target.MarshalBinary()

	if err != nil {
		return err
	}

	// Write to file
	err = ioutil.WriteFile(
		prefix(context, target.GetPath()), buffer, 0600,
	)

	if err != nil {
		return err
	}

	log.Printf("Saved %s", target)

	return nil
}

func (this FileAccess) SaveWithDir(
	target Storeable, directory string, overwrite bool, context ServerContext,
) error {
	// Create user directory if it doesn't already exist
	dir := prefix(context, directory)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
	}

	return this.Save(target, overwrite, context)
}

func (this FileAccess) LoadRaw(
	target Loadable, context ServerContext,
) ([]byte, error) {
	buffer, err := ioutil.ReadFile(prefix(context, target.GetPath()))

	if err != nil {
		return nil, err
	}

	log.Printf("Loaded %s", target)

	return buffer, nil
}

func (this FileAccess) Load(
	target Loadable, context ServerContext,
) error {
	buffer, err := this.LoadRaw(target, context)

	if err != nil {
		return err
	}

	err = target.UnmarshalBinary(buffer)

	if err != nil {
		return err
	}

	return nil
}

func (this FileAccess) Remove(
	target Accessable, context ServerContext,
) error {
	if err := os.Remove(prefix(context, target.GetPath())); err != nil {
		return err
	}

	log.Printf("Deleted %s", target)

	return nil
}

func (this FileAccess) RemoveDir(
	directory string, context ServerContext,
) error {
	if err := os.RemoveAll(prefix(context, directory)); err != nil {
		return err
	}

	return nil
}
