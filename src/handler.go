package main

import (
  "net/http"
)

type handler struct {
  fileServer http.Handler
}

func newHandler() *handler {
  this := &handler{}
  this.fileServer = http.FileServer(http.Dir(ROOT))
  return this
}

func (this *handler) ServeHTTP(
  writer http.ResponseWriter, request *http.Request,
) {
  this.fileServer.ServeHTTP(writer, request)
}
