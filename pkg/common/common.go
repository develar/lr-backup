package common

import (
  "golang.org/x/oauth2"
  "os"
)

var ClientId string
var ClientSecret string

// https://console.adobe.io/
var AdobeOauthConfig *oauth2.Config

func init() {
  if len(ClientId) == 0 {
    ClientId = os.Getenv("CLIENT_ID")
  }
  if len(ClientSecret) == 0 {
    ClientSecret = os.Getenv("CLIENT_SECRET")
  }

  AdobeOauthConfig = &oauth2.Config{
    ClientID:     ClientId,
    ClientSecret: ClientSecret,
    Scopes:       []string{"AdobeID", "openid", "lr_partner_apis"},
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
