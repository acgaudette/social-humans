package main

import (
  "html/template"
  "net/http"
)

type router struct {
  mux *http.ServeMux
}

func newRouter() *router {
  this := &router{}
  this.mux = http.NewServeMux()

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
  writer http.ResponseWriter, path string, data *user,
) error {
  t, err := template.ParseFiles(ROOT + path)

  if err != nil {
    return err
  }

  err = t.Execute(writer, data)

  return err
}
