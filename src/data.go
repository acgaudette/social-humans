package main

import (
  "bytes"
  "crypto/sha256"
  "errors"
  "io/ioutil"
  "os"
)

type user struct {
  Handle string
  hash   []byte
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
  this.hash = hash(cleartext)
}

func (this *user) validate(cleartext string) bool {
  return bytes.Equal(hash(cleartext), this.hash)
}

func (this *user) save(overwrite bool) error {
  _, err := os.Stat(userpath(this.Handle))

  if !os.IsNotExist(err) && !overwrite {
    return errors.New("user file already exists")
  }

  return ioutil.WriteFile(
    userpath(this.Handle),
    this.hash,
    0600,
  )
}

func addUser(handle string, password string) (*user, error) {
  account := &user{
    Handle: handle,
  }

  account.setPassword(password)

  if err := account.save(false); err != nil {
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
    hash:   in,
  }, nil
}

type statusMessage struct {
  Status string
}
