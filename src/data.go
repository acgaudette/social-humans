package main

import (
  "bytes"
  "crypto/sha256"
  "io/ioutil"
)

type user struct {
  Handle string
  Hash   []byte
}

func userpath(handle string) string {
  return DATA_PATH + "/" + handle + ".user"
}

func hash(cleartext string) []byte {
  hash := sha256.New()
  hash.Write([]byte(cleartext))
  return hash.Sum(nil)
}

func (this *user) setPassword(cleartext string) {
  this.Hash = hash(cleartext)
}

func (this *user) validate(cleartext string) bool {
  return bytes.Equal(hash(cleartext), this.Hash)
}

func (this *user) save() error {
  return ioutil.WriteFile(
    userpath(this.Handle),
    this.Hash,
    0600,
  )
}

func addUser(handle string, password string) (*user, error) {
  account := &user{
    Handle: handle,
  }

  account.setPassword(password)

  if err := account.save(); err != nil {
    return nil, err
  }

  return account, nil
}

func loadUser(handle string) (*user, error) {
  in, err := ioutil.ReadFile(userpath(handle))

  if err != nil {
    return nil, err
  }

  return &user{
    Handle: handle,
    Hash:   in,
  }, nil
}

type statusMessage struct {
  Status string
}
