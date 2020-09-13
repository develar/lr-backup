// Package obscure contains the Obscure and Reveal commands
package common

import (
  "crypto/rand"
  "encoding/base64"
  "encoding/json"
  "errors"
  "golang.org/x/crypto/nacl/box"
  "os"
)

var InPk string
var InSk string
var OutPk string

var outPublicKey *[32]byte
var inPublicKey *[32]byte
var inPrivateKey *[32]byte

func init() {
  outPublicKey = decodeKey("OUT_PK", OutPk)
  inPublicKey = decodeKey("IN_PK", InPk)
  inPrivateKey = decodeKey("IN_SK", InSk)
}

func setupTestOnly(outPKey *[32]byte, inPKey *[32]byte, inSKey *[32]byte) {
  outPublicKey = outPKey
  inPublicKey = inPKey
  inPrivateKey = inSKey
}

func decodeKey(envName string, base64EncodedValue string) *[32]byte {
  if len(base64EncodedValue) == 0 {
    base64EncodedValue = os.Getenv(envName)
  }

  var err error
  result, err := base64.RawURLEncoding.DecodeString(base64EncodedValue)
  if err != nil {
    panic(err)
  }

  var bytes [32]byte
  copy(bytes[:], result)
  return &bytes
}

func Encrypt(message interface{}) (string, error) {
  result, err := json.Marshal(message)
  if err != nil {
    return "", err
  }

  result, err = box.SealAnonymous(nil, result, outPublicKey, rand.Reader)
  if err != nil {
    return "", err
  }
  return base64.RawURLEncoding.EncodeToString(result), nil
}

func Decrypt(message string, v interface{}) error {
  result, err := base64.RawURLEncoding.DecodeString(message)
  if err != nil {
    return err
  }

  result, ok := box.OpenAnonymous(nil, result, inPublicKey, inPrivateKey)
  if !ok {
    return errors.New("cannot decrypt state")
  }

  return json.Unmarshal(result, v)
}

func DecryptBytes(data []byte) (string, error) {
  result, ok := box.OpenAnonymous(nil, data, inPublicKey, inPrivateKey)
  if !ok {
    return "", errors.New("cannot decrypt state")
  }
  return string(result), nil
}
