package main

import (
  "errors"
  "html/template"
  "log"
  "net/http"
  "time"
)

type router struct {
  mux *http.ServeMux
}

func newRouter() *router {
  this := &router{}
  this.mux = http.NewServeMux()

  this.mux.HandleFunc("/", index)
  this.mux.HandleFunc("/login", login)
  this.mux.HandleFunc("/logout", logout)

  return this
}

func loadSession(request *http.Request) (*user, error) {
  session, err := request.Cookie(SESSION_NAME)

  if err != nil {
    return nil, err
  }

  data, err := loadUser(session.Value)

  if err != nil {
    return nil, err
  }

  return data, nil
}

func serveTemplate(
  writer http.ResponseWriter, path string, data interface{},
) error {
  t, err := template.ParseFiles(ROOT + path)

  if err != nil {
    return err
  }

  err = t.Execute(writer, data)

  return err
}

func error501(writer http.ResponseWriter) {
  http.Error(
    writer,
    http.StatusText(http.StatusInternalServerError),
    http.StatusInternalServerError,
  )
}

func index(writer http.ResponseWriter, request *http.Request) {
  switch request.Method {
  case "GET":
    if request.URL.Path == "/" {

      data, err := loadSession(request)

      if err != nil {
        log.Printf("%s", err)
        http.Redirect(writer, request, "/login", http.StatusFound)
        break
      }

      if err = serveTemplate(writer, "/index.html", data); err != nil {
        log.Printf("%s", err)
      }
    } else {
      http.NotFound(writer, request)
    }
  }
}

func login(writer http.ResponseWriter, request *http.Request) {
  switch request.Method {
  case "GET":
    err := serveTemplate(writer, "/login.html", &statusMessage{Status: ""})

    if err != nil {
      error501(writer)
    }

  case "POST":
    request.ParseForm()

    readLoginForm := func(key string, errorStatus string) (string, error) {
      result := request.Form.Get(key)

      if result == "" {
        message := statusMessage{Status: errorStatus}
        err := serveTemplate(writer, "/login.html", &message)

        if err != nil {
          error501(writer)
          return "", err
        }

        return "", errors.New("key not found")
      }

      return result, nil
    }

    handle, err := readLoginForm("handle", "Username required!")

    if err != nil {
      return
    }

    password, err := readLoginForm("password", "Password required!")

    if err != nil {
      return
    }

    account, err := loadUser(handle)

    if err != nil {
      account, err = addUser(handle, password)

      if err != nil {
        error501(writer)
        return
      }
    } else if !account.validate(password) {
      message := statusMessage{Status: "Invalid password"}
      err := serveTemplate(writer, "/login.html", &message)

      if err != nil {
        error501(writer)
      }

      return
    }

    session := http.Cookie{
      Name:  SESSION_NAME,
      Value: account.Handle,
    }

    http.SetCookie(writer, &session)
    http.Redirect(writer, request, "/", http.StatusFound)
  }
}

func logout(writer http.ResponseWriter, request *http.Request) {
  switch request.Method {
  case "GET":
    http.Redirect(writer, request, "/", http.StatusFound)

  case "POST":
    session := http.Cookie{
      Name:    SESSION_NAME,
      Value:   "",
      Expires: time.Now().Add(-time.Minute),
    }

    http.SetCookie(writer, &session)
    http.Redirect(writer, request, "/", http.StatusFound)
  }
}
