package main

import (
  "fmt"
  "github.com/aws/aws-lambda-go/events"
  "github.com/aws/aws-lambda-go/lambda"
  "html/template"
  "strings"
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

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
  var responseBody strings.Builder
  data := struct {
    OK          bool
    Error string
  }{
    OK: true,
  }

  state := request.QueryStringParameters["state"]
  if len(state) == 0 {
    data.Error = "state is not provided"
    data.OK = false
  }

  code := request.QueryStringParameters["code"]
  if len(code) == 0 {
    data.Error = "code is not provided"
    data.OK = false
  }

  t, err := template.New("authResponse").Parse(authResponseTemplate)
  if err != nil {
    return nil, fmt.Errorf("could not parse template for web response: %w", err)
  }

  err = t.Execute(&responseBody, data)
  if err != nil {
    return nil, fmt.Errorf("could not execute template for web response: %w", err)
  }

  responseBody.WriteString(state)
  responseBody.WriteString("<br>")
  responseBody.WriteString(code)

  return &events.APIGatewayProxyResponse{
    StatusCode: 200,
    Headers:    map[string]string{"Content-Type": "text/html"},
    Body:       responseBody.String(),
  }, nil
}

func main() {
  lambda.Start(handler)
}
