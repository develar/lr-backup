package common

import (
  "golang.org/x/oauth2"
  "os"
)

var clientId string
var clientSecret string

// https://console.adobe.io/
var AdobeOauthConfig *oauth2.Config

func init() {
  if len(clientId) == 0 {
    clientId = os.Getenv("CLIENT_ID")
    if len(clientId) == 0 {
      panic("clientId is not set")
    }
  }
  if len(clientSecret) == 0 {
    clientSecret = os.Getenv("CLIENT_SECRET")
    if len(clientSecret) == 0 {
      panic("clientSecret is not set")
    }
  }

  AdobeOauthConfig = &oauth2.Config{
    ClientID:     clientId,
    ClientSecret: clientSecret,
    Scopes:       []string{"AdobeID,openid,lr_partner_apis"},
    Endpoint: oauth2.Endpoint{
      AuthURL:  "https://ims-na1.adobelogin.com/ims/authorize/v1",
      TokenURL: "https://ims-na1.adobelogin.com/ims/token/v2",
    },
  }
}

type AuthRequest struct {
  Port  int
  Token string
}

type AuthResponse struct {
  Code  string
  Token string
}
