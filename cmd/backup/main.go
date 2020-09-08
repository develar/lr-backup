package main

import (
  "context"
  "crypto/rand"
  "encoding/base64"
  "encoding/json"
  "fmt"
  "github.com/develar/lr-backup/pkg/common"
  "io"
  "log"
  "net"
  "net/http"
)

var oauthStateString string
var port int

func main() {
  http.HandleFunc("/", handleMain)
  http.HandleFunc("/login", handleLogin)
  http.HandleFunc("/callback", handleCallback)

  listener, err := net.Listen("tcp", "127.0.0.1:0")
  if err != nil {
    log.Fatal(err)
  }

  port = listener.Addr().(*net.TCPAddr).Port
  log.Print("Listen http://" + listener.Addr().String())
  err = http.Serve(listener, nil)
  if err != nil {
    log.Fatal(err)
  }
}

func handleMain(w http.ResponseWriter, _ *http.Request) {
  var htmlIndex = `<html>
<body>
	<a href="/login">Authenticate</a>
</body>
</html>`
  _, _ = fmt.Fprintf(w, htmlIndex)
}

func tokenGenerator() string {
  data := make([]byte, 32)
  _, err := io.ReadFull(rand.Reader, data)
  if err != nil {
    log.Fatal(err)
  }
  return base64.RawURLEncoding.EncodeToString(data)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
  stateEncoded, err := json.Marshal(AuthRequest{Port: port, Token: tokenGenerator()})
  if err != nil {
    log.Fatal(err)
  }

  url := common.AdobeOauthConfig.AuthCodeURL(base64.RawURLEncoding.EncodeToString(stateEncoded))
  http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
  err := readCallbackResponse(r.FormValue("state"), r.FormValue("code"))
  if err != nil {
    fmt.Println(err.Error())
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
    return
  }
}

func readCallbackResponse(state string, code string) error {
  if state != oauthStateString {
    return fmt.Errorf("invalid oauth state")
  }

  token, err := common.AdobeOauthConfig.Exchange(context.Background(), code)
  if err != nil {
    return fmt.Errorf("code exchange failed: %s", err.Error())
  }

  log.Printf(token.AccessToken)
  return nil
}

type AuthRequest struct {
  Port  int
  Token string
}
