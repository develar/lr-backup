package main

import (
  "github.com/develar/lr-backup/pkg/common"
  "github.com/shibukawa/configdir"
  "log"
  "os"
  "runtime"
)

func readToken(cache *configdir.Config) string {
  data, err := cache.ReadFile("token")
  if err != nil && !os.IsNotExist(err) {
    LogError("cannot read token", err)
    return ""
  }

  if len(data) == 0 {
    return ""
  }

  token, err := common.DecryptBytes(data)
  if err != nil {
    LogError("cannot read token", err)
    return ""
  }

  return token
}

func LogError(reason string, err error) {
	log.Println("Fyne error: ", reason)
	if err != nil {
		log.Println("  Cause:", err)
	}

	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Printf("  At: %s:%d", file, line)
	}
}
