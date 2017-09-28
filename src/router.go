package main

import (
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
