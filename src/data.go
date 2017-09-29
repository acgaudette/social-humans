package main

import (
  "bufio"
  "bytes"
  "crypto/sha256"
  "errors"
  "io/ioutil"
  "os"
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
    append([]byte(this.Handle+"\n"), this.Hash...),
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
  file, err := os.Open(userpath(handle))

  if err != nil {
    return nil, err
  }

  defer file.Close()

  scanner := bufio.NewScanner(file)

  scanner.Scan()
  handle = scanner.Text()
  scanner.Scan()
  hash := []byte(scanner.Text()) // safe?

  if err = scanner.Err(); err != nil {
    return nil, err
  }

  return &user{
    Handle: handle,
    Hash:   hash,
  }, nil
}

type statusMessage struct {
  Status string
}
