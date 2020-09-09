package main

import (
  "context"
  "log"
  "os"
  "os/signal"
  "time"
)

var oauthStateToken string
var port int

func main() {
  stop := make(chan os.Signal, 1)
  signal.Notify(stop, os.Interrupt, os.Kill)

  server, tokenChannel := startServer()

  select {
  case <-stop:
    log.Println("interrupted")
  case token := <-tokenChannel:
    {
      ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
      defer cancel()
      err := server.Shutdown(ctx)
      if err != nil {
        log.Printf("cannot shutdown server: %v", err)
      }
      log.Printf("got code: %s", token.AccessToken)
    }
  }
}
