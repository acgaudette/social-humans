package main

import (
  "io/ioutil"
)

type user struct {
  Handle string
}

func userpath(handle string) string {
  return DATA_PATH + "/" + handle + ".user"
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
