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
      log.Printf("Listening on http://%s", server.Addr)
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

func shutdown(server *http.Server) error {
  log.Printf("Shutting down...")

  background, _ := context.WithTimeout(
    context.Background(), EXIT_TIMEOUT*time.Second,
  )

  return server.Shutdown(background)
}
