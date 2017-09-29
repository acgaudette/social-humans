package main

import (
  "io/ioutil"
  "crypto/sha256"
)

type user struct {
  Handle string
  Hash []byte
}

func userpath(handle string) string {
  return DATA_PATH + "/" + handle + ".user"
}

func hash(cleartext string) []byte {
  hash := sha256.New()
  hash.Write([]byte(cleartext))
  return hash.Sum(nil)
}

func (this *user) save() error {
  return ioutil.WriteFile(
    userpath(this.Handle),
    []byte(this.Handle),
    0600,
  )
}

func loadUser(handle string) (*user, error) {
  in, err := ioutil.ReadFile(userpath(handle))

  if err != nil {
    return nil, err
  }

  return &user{
    Handle: string(in),
  }, nil
}

type statusMessage struct {
  Status string
}
