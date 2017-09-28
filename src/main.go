package main

import (
  "log"
  "net/http"
  "os"
  "os/signal"
)

const (
  ADDRESS = "0.0.0.0"
  PORT    = "5368"
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
        log.Printf("Shutting down...")
        return server.Shutdown(nil)
      }

      return nil

    case <-restart:
      log.Printf("Listening on http://%s:%s...", ADDRESS, PORT)
      go listen(server, restart)

    default:
      if server == nil {
        server = makeServer()
        restart <- true
      }
    }
  }
}

func makeServer() *http.Server {
  server := &http.Server{
    Addr:    ADDRESS + ":" + PORT,
    Handler: http.FileServer(http.Dir("www")),
  }

  return server
}

func listen(server *http.Server, failure chan bool) {
  if err := server.ListenAndServe(); err != nil {
    log.Printf("%s", err)
    failure <- true
  }
}
