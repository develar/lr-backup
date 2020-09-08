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
Error: {{ .Name }}<br>
{{ if .Description }}Description: {{ .Description }}<br>{{ end }}
{{ if .Code }}Code: {{ .Code }}<br>{{ end }}
{{ if .HelpURL }}Look here for help: <a href="{{ .HelpURL }}">{{ .HelpURL }}</a><br>{{ end }}
{{ else }}
All done. Please go back to lr-backup.
{{ end }}
</pre>
</body>
</html>
`

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
  var t = template.Must(template.New("authResponse").Parse(authResponseTemplate))

  var responseBody strings.Builder
  var err = t.Execute(&responseBody, struct {
      OK            bool
      numberOfCarts int
    }{
      OK: true,
    })
  if err != nil {
    return nil, fmt.Errorf("could not execute template for web response: %w", err)
  }

  return &events.APIGatewayProxyResponse{
    StatusCode: 200,
    Headers:    map[string]string{"Content-Type": "text/html"},
    Body:       responseBody.String(),
  }, nil
}

func main() {
  // Make the handler available for Remote Procedure Call by AWS Lambda
  lambda.Start(handler)
}
