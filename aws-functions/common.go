package main

import "golang.org/x/oauth2"

//goland:noinspection GoVarAndConstTypeMayBeOmitted
var ClientId string = "5af3a11792fa45d3ab6575d93e510785"
//goland:noinspection GoVarAndConstTypeMayBeOmitted
var ClientSecret string = "39cedc11-ff67-4945-9f9c-a79b87a5c7fe"

var adobeOauthConfig *oauth2.Config = &oauth2.Config{
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
