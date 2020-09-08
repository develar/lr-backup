package main

import (
  "encoding/base64"
  "encoding/json"
  "errors"
  "github.com/aws/aws-lambda-go/events"
  "github.com/aws/aws-lambda-go/lambda"
  "github.com/develar/lr-backup/pkg/common"
  "net/http"
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
  decoded, err := base64.RawURLEncoding.DecodeString(state)
  err = json.Unmarshal(decoded, &authRequest)
  if err != nil {
    return nil, err
  }

  return &events.APIGatewayProxyResponse{
    StatusCode: http.StatusTemporaryRedirect,
    Headers:    map[string]string{"Location": "http://127.0.0.1:" + strconv.Itoa(authRequest.Port) + "?code=" + code},
  }, nil
}

func main() {
  lambda.Start(handler)
}
