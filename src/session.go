package main

import (
  "net/http"
  "strings"
  "log"
  "time"
)

func newSession(writer http.ResponseWriter, account *user) {
  token := account.refreshToken()

  session := http.Cookie{
    Name:  SESSION_NAME,
    Value: account.Handle + DELM + token,
  }

  http.SetCookie(writer, &session)
  log.Printf("Created new session with token %s", token)
}

func clearSession(writer http.ResponseWriter) {
  session := http.Cookie{
    Name:    SESSION_NAME,
    Value:   "",
    Expires: time.Now().Add(-time.Minute),
  }

  http.SetCookie(writer, &session)
  log.Printf("Cleared session")
}

func loadSession(request *http.Request) (*user, error) {
  session, err := request.Cookie(SESSION_NAME)

  if err != nil {
    return nil, err
  }

  split := strings.Split(session.Value, DELM)
  account, err := loadUser(split[0])

  if err != nil {
    return nil, err
  }

  if err = account.checkToken(split[1]); err != nil {
    return nil, err
  }

  log.Printf("Loaded session with token %s", split[1])
  return account, nil
}
