package main

import (
  "context"
  "crypto/rand"
  "encoding/base64"
  "fmt"
  "io"
  "log"
  "net/http"
  "strconv"
)

var bindPort = 53672

var oauthStateString string

func main() {
  http.HandleFunc("/", handleMain)
  http.HandleFunc("/login", handleGoogleLogin)
  http.HandleFunc("/callback", handleCallback)
  hostAndPort := "localhost:" + strconv.Itoa(bindPort)
  log.Print("Listen http://" + hostAndPort)
  err := http.ListenAndServe(hostAndPort, nil)
  if err != nil {
    log.Fatal(err)
  }
}

func handleMain(w http.ResponseWriter, _ *http.Request) {
  var htmlIndex = `<html>
<body>
	<a href="/login">Adobe Log In</a>
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
  return base64.URLEncoding.EncodeToString(data)
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
  oauthStateString = tokenGenerator()
  url := adobeOauthConfig.AuthCodeURL(oauthStateString)
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

  token, err := adobeOauthConfig.Exchange(context.Background(), code)
  if err != nil {
    return fmt.Errorf("code exchange failed: %s", err.Error())
  }

  log.Printf(token.AccessToken)
  return nil
}
