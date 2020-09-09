package main

import (
  "context"
  "crypto/rand"
  "encoding/base64"
  "fmt"
  "github.com/develar/lr-backup/pkg/common"
  "golang.org/x/oauth2"
  "html/template"
  "io"
  "log"
  "net"
  "net/http"
)

var authResponseTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>{{ if .OK }}Success!{{ else }}Failure!{{ end }}</title>
</head>
<body>
<h1>{{ if .OK }}Success!{{ else }}Failure!{{ end }}</h1>
<hr>
<pre style="width: 750px; white-space: pre-wrap;">
{{ if eq .OK false }}
  {{ if .Error }}Error: {{ .Error }}<br>{{ end }}
{{ else }}
  All done. Please go back to lr-backup.
{{ end }}
</pre>
</body>
</html>
`

func startServer() (*http.Server, chan *oauth2.Token) {
  tokenChannel := make(chan *oauth2.Token)

  serveMux := http.NewServeMux()
  var server = &http.Server{Handler: serveMux}
  serveMux.HandleFunc("/", handleMain)
  serveMux.HandleFunc("/login", handleLogin)
  serveMux.HandleFunc("/callback", func(writer http.ResponseWriter, request *http.Request) {
    var authResponse common.AuthResponse
    err := common.Decrypt(request.URL.Query().Get("r"), &authResponse)
    if err != nil {
      log.Print(err)
      http.Error(writer, err.Error(), http.StatusBadRequest)
      return
    }

    if authResponse.Token != oauthStateToken {
      http.Error(writer, "state is not equal to expected", http.StatusBadRequest)
      return
    }

    token, err := exchangeCodeAndWriteResponse(authResponse, writer)
    if err != nil {
      log.Print(err)
      http.Error(writer, err.Error(), http.StatusInternalServerError)
      return
    }

    tokenChannel <- token
  })

  listener, err := net.Listen("tcp", "127.0.0.1:0")
  if err != nil {
    log.Fatal(err)
  }

  port = listener.Addr().(*net.TCPAddr).Port
  go func() {
    log.Print("Listen http://" + listener.Addr().String())
    err = server.Serve(listener)
    if err != nil {
      log.Fatal(err)
    }
  }()
  return server, tokenChannel
}

func handleMain(w http.ResponseWriter, _ *http.Request) {
  var htmlIndex = `<html>
<body>
	<a href="/login">Authenticate</a>
</body>
</html>`
  _, _ = fmt.Fprintf(w, htmlIndex)
}

func generateToken() string {
  data := make([]byte, 16)
  _, err := io.ReadFull(rand.Reader, data)
  if err != nil {
    log.Fatal(err)
  }
  return base64.RawURLEncoding.EncodeToString(data)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
  oauthStateToken = generateToken()
  stateEncoded, err := common.Encrypt(common.AuthRequest{Port: port, Token: oauthStateToken})
  if err != nil {
    log.Print(err)
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  url := common.AdobeOauthConfig.AuthCodeURL(stateEncoded)
  http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func exchangeCodeAndWriteResponse(response common.AuthResponse, w http.ResponseWriter) (*oauth2.Token, error) {
  token, err := common.AdobeOauthConfig.Exchange(context.Background(), response.Code)
  if err != nil {
    return nil, fmt.Errorf("code exchange failed: %w", err)
  }

  log.Printf(token.AccessToken)

  data := struct {
    OK    bool
    Error string
  }{
    OK: true,
  }

  t, err := template.New("authResponse").Parse(authResponseTemplate)
  if err != nil {
    return nil, fmt.Errorf("could not parse template for web response: %w", err)
  }

  err = t.Execute(w, data)
  if err != nil {
    return nil, fmt.Errorf("could not execute template for web response: %w", err)
  }

  return token, nil
}
