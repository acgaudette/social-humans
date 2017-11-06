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
	Save(Storeable, bool, serverContext) error
	SaveWithDir(Storeable, string, bool, serverContext) error
	Load(Loadable, serverContext) error
	LoadRaw(Loadable, serverContext) ([]byte, error)
	Remove(Accessable, serverContext) error
	RemoveDir(string, serverContext) error
}

type FileAccess struct{}

func (this FileAccess) Save(
	target Storeable, overwrite bool, context serverContext,
) error {
	_, err := os.Stat(prefix(context, target.GetPath()))

	// Don't overwrite unless specified
	if !os.IsNotExist(err) && !overwrite {
		return fmt.Errorf("data file already exists")
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
	target Storeable, directory string, overwrite bool, context serverContext,
) error {
	// Create user directory if it doesn't already exist
	dir := prefix(context, directory)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
	}

	return this.Save(target, overwrite, context)
}

func (this FileAccess) LoadRaw(
	target Loadable, context serverContext,
) ([]byte, error) {
	buffer, err := ioutil.ReadFile(prefix(context, target.GetPath()))

	if err != nil {
		return nil, err
	}

	log.Printf("Loaded %s", target)

	return buffer, nil
}

func (this FileAccess) Load(
	target Loadable, context serverContext,
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
	target Accessable, context serverContext,
) error {
	if err := os.Remove(prefix(context, target.GetPath())); err != nil {
		return err
	}

	log.Printf("Deleted %s", target)

	return nil
}

func (this FileAccess) RemoveDir(
	directory string, context serverContext,
) error {
	if err := os.RemoveAll(prefix(context, directory)); err != nil {
		return err
	}

	return nil
}
