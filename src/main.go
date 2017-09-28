package main

import (
  "context"
  "log"
  "net/http"
  "os"
  "os/signal"
  "time"
)

const (
  ADDRESS      = "0.0.0.0"
  PORT         = "5368"
  ROOT         = "www"
  EXIT_TIMEOUT = 4
)

func main() {
  if err := run(); err != nil {
    log.Fatal(err)
  }
}

func run() error {
  interrupt := make(chan os.Signal, 1)
  signal.Notify(interrupt, os.Interrupt)

  restart := make(chan bool, 1)
  var server *http.Server

  for {
    select {
    case <-interrupt:
      if server != nil {
        return shutdown(server)
      }

      return nil

    case <-restart:
      go listen(server, restart)

    default:
      if server == nil {
        server = newServer()
        restart <- true
      }
    }
  }
}

func newServer() *http.Server {
  server := &http.Server{
    Addr:    ADDRESS + ":" + PORT,
    Handler: newHandler(),
  }

  return server
}

func listen(server *http.Server, failure chan bool) {
  log.Printf("Listening on http://%s", server.Addr)

  if err := server.ListenAndServe(); err != nil {
    log.Printf("%s", err)
    failure <- true
  }
}

func shutdown(server *http.Server) error {
  log.Printf("Shutting down...")

  background, _ := context.WithTimeout(
    context.Background(), EXIT_TIMEOUT*time.Second,
  )

  return server.Shutdown(background)
}

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
