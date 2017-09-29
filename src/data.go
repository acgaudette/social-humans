package main

import (
  "bufio"
  "bytes"
  "crypto/rand"
  "crypto/sha256"
  "errors"
  "fmt"
  "io/ioutil"
  "log"
  "os"
)

func hash(cleartext string) []byte {
  hash := sha256.New()
  hash.Write([]byte(cleartext))
  return hash.Sum(nil)
}

type statusMessage struct {
  Status string
}

type user struct {
  Handle string
  hash   []byte
  token  string
}

func (this *user) refreshToken() string {
  buffer := make([]byte, 32)
  rand.Read(buffer)
  this.token = fmt.Sprintf("%x", buffer)

  this.save(true)
  return this.token
}

func (this *user) checkToken(token string) error {
  if token == this.token {
    return nil
  }

  return errors.New("token mismatch")
}

func (this *user) setPassword(cleartext string) {
  this.hash = hash(cleartext)
}

func (this *user) validate(cleartext string) error {
  if bytes.Equal(hash(cleartext), this.hash) {
    return nil
  }

  return errors.New("password hash mismatch")
}

func (this *user) save(overwrite bool) error {
  _, err := os.Stat(userpath(this.Handle))

  if !os.IsNotExist(err) && !overwrite {
    return errors.New("user file already exists")
  }

  return ioutil.WriteFile(
    userpath(this.Handle),
    append([]byte(this.token+"\n"), this.hash...),
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

  log.Printf("Added user \"%s\"", handle)
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
  token := scanner.Text()
  scanner.Scan()
  hash := scanner.Bytes()

  if err = scanner.Err(); err != nil {
    return nil, err
  }

  log.Printf("Loaded user \"%s\"", handle)

  return &user{
    Handle: handle,
    hash:   hash,
    token:  token,
  }, nil
}

func userpath(handle string) string {
  return DATA_PATH + "/" + handle + ".user"
}
