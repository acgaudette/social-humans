package smhb

import (
	"encoding"
	"io/ioutil"
	"os"
	"fmt"
)

type Storeable interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	GetPath() string
}

type Access interface {
	Save(Storeable, bool, serverContext) error
	Load(Storeable, serverContext) error
	LoadRaw(Storeable, serverContext) ([]byte, error)
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
	return ioutil.WriteFile(
		prefix(context, target.GetPath()), buffer, 0600,
	)
}

func (this FileAccess) LoadRaw(
	target Storeable, context serverContext,
) ([]byte, error) {
	buffer, err := ioutil.ReadFile(prefix(context, target.GetPath()))

	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func (this FileAccess) Load(
	target Storeable, context serverContext,
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
