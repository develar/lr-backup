package common

import "golang.org/x/oauth2"

var ClientId string
var ClientSecret string

// https://console.adobe.io/
var AdobeOauthConfig = &oauth2.Config{
  //RedirectURL:  "http://localhost:" + strconv.Itoa(bindPort) + "/callback",
  RedirectURL:  "http://localhost:53672/callback",
  ClientID:     ClientId,
  ClientSecret: ClientSecret,
  Scopes:       []string{"AdobeID", "creative_sdk"},
  Endpoint: oauth2.Endpoint{
    AuthURL:  "https://ims-na1.adobelogin.com/ims/authorize/v1",
    TokenURL: "https://ims-na1.adobelogin.com/ims/token/v2",
  },
}
