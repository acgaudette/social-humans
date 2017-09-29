package main

import (
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
    http.ServeFile(writer, request, ROOT+"/login.html")

  case "POST":
    request.ParseForm()

    handle := request.FormValue("handle")
    account, err := loadUser(handle)

    if err != nil {
      account = &user{
        Handle: handle,
      }

      account.save()
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
