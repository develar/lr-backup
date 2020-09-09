package main

import (
  "errors"
  "github.com/aws/aws-lambda-go/events"
  "github.com/aws/aws-lambda-go/lambda"
  "github.com/develar/lr-backup/pkg/common"
  "net/http"
  "net/url"
  "strconv"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
  state := request.QueryStringParameters["state"]
  if len(state) == 0 {
    return nil, errors.New("state is not provided")
  }

  code := request.QueryStringParameters["code"]
  if len(code) == 0 {
    return nil, errors.New("code is not provided")
  }

  var authRequest common.AuthRequest
  err := common.Decrypt(state, &authRequest)
  if err != nil {
    return nil, err
  }

  authResponse, err := common.Encrypt(common.AuthResponse{Code: code, Token: authRequest.Token})
  if err != nil {
    return nil, err
  }

  redirectQuery := url.Values{}
  redirectQuery.Set("r", authResponse)
  redirectUrl := url.URL{
    Scheme:   "http",
    Host:     "127.0.0.1:" + strconv.Itoa(authRequest.Port),
    Path:     "/callback",
    RawQuery: redirectQuery.Encode(),
  }

  return &events.APIGatewayProxyResponse{
    StatusCode: http.StatusTemporaryRedirect,
    Headers:    map[string]string{"Location": redirectUrl.String()},
  }, nil
}

func main() {
  lambda.Start(handler)
}
