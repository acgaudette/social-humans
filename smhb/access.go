package smhb

import (
	"encoding"
	"io/ioutil"
)

type Storeable interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	GetPath() string
}

type Access interface {
	Save(Storeable, serverContext) error
	Load(Storeable, serverContext) error
	LoadRaw(Storeable, serverContext) ([]byte, error)
}

type FileAccess struct{}

func (this FileAccess) Save(target Storeable, context serverContext) error {
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
