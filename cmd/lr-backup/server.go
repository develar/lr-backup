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
  "net/url"
)

var authResponseTemplate = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{{- if .Error -}}Something went wrong{{- else -}}You're all set{{- end -}}</title>
     <script defer src="https://use.fontawesome.com/releases/v5.3.1/js/all.js"></script>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma@0.9.0/css/bulma.min.css" integrity="sha256-aPeK/N8IHpHsvPBCf49iVKMdusfobKo2oxF8lRruWJg=" crossorigin="anonymous">
  </head>
  <body>
  <section class="section">
    <div class="container">
      <h1 class="title">
        {{- if .Error -}}
          Something went wrong
        {{- else -}}
          You're all set
        {{- end -}}
      </h1>
      <p class="subtitle">
        {{- if .Error -}}
          An unknown error has occurred. Here's the error message if you want to file issue: {{ .Error }}
        {{-  else -}}
          You’ve successfully signed in. Feel free to close this browser tab and return to where you previously left off.
        {{- end }}
      </p>
    </div>
  </section>
  </body>
</html>
`

func startServer() (*http.Server, chan *oauth2.Token, error) {
  tokenChannel := make(chan *oauth2.Token)

  authResponseTemplate, err := template.New("authResponse").Parse(authResponseTemplate)
  if err != nil {
    return nil, nil, err
  }

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
      writer.Header().Set("Content-Type", "text/html; charset=utf-8")
      writer.WriteHeader(http.StatusBadRequest)
      err2 := authResponseTemplate.Execute(writer, struct {
        Error string
      }{
        Error: "state is not equal to expected",
      })
      if err2 != nil {
        log.Print(err)
      }
      return
    }

    token, err := exchangeCodeAndWriteResponse(authResponse, writer, authResponseTemplate)
    if err != nil {
      tokenChannel <- nil
      log.Print(err)
      http.Error(writer, err.Error(), http.StatusInternalServerError)
      return
    }

    tokenChannel <- token
  })

  listener, err := net.Listen("tcp", "127.0.0.1:0")
  if err != nil {
    return nil, nil, err
  }

  port = listener.Addr().(*net.TCPAddr).Port
  go func() {
    log.Print("Listen http://" + listener.Addr().String())
    err = server.Serve(listener)
    if err != nil {
      log.Fatal(err)
    }
  }()
  return server, tokenChannel, nil
}

func handleMain(w http.ResponseWriter, _ *http.Request) {
  var htmlIndex = `<html>
<body>
	<a href="/login">Authenticate</a>
</body>
</html>`
  _, _ = fmt.Fprint(w, htmlIndex)
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
  authCodeUrl, err := generateSignInUrl()
  if err != nil {
    log.Print(err)
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  http.Redirect(w, r, authCodeUrl.String(), http.StatusTemporaryRedirect)
}

func generateSignInUrl() (*url.URL, error) {
  oauthStateToken = generateToken()
  stateEncoded, err := common.Encrypt(common.AuthRequest{Port: port, Token: oauthStateToken})
  if err != nil {
    return nil, err
  }

  //res.redirect(`https://ims-na1.adobelogin.com/ims/authorize?&scope=openid,creative_sdk&response_type=code&redirect_uri=https://localhost:8000/callback`)
  //})
  redirectQuery := url.Values{}
  redirectQuery.Set("client_id", common.AdobeOauthConfig.ClientID)
  redirectQuery.Set("scope", "AdobeID,openid,lr_partner_apis")
  redirectQuery.Set("state", stateEncoded)
  return &url.URL{
    Scheme:   "https",
    Host:     "ims-na1.adobelogin.com/ims/authorize",
    Path:     "/",
    RawQuery: redirectQuery.Encode(),
  }, nil
}

func exchangeCodeAndWriteResponse(response common.AuthResponse, w http.ResponseWriter, authResponseTemplate *template.Template) (*oauth2.Token, error) {
  data := struct {
    Error string
  }{
  }

  token, err := common.AdobeOauthConfig.Exchange(context.Background(), response.Code)
  if err != nil {
    data.Error = err.Error()
    token = nil
  }

  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  err = authResponseTemplate.Execute(w, data)
  if err != nil {
    return nil, fmt.Errorf("could not execute template for web response: %w", err)
  }

  return token, nil
}
