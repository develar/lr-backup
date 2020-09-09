package common

import (
  "crypto/rand"
  "golang.org/x/crypto/nacl/box"
  "testing"
)

func TestExtractDuplicate(t *testing.T) {
  senderPublicKey, senderPrivateKey, err := box.GenerateKey(rand.Reader)
  if err != nil {
    t.Error(err)
  }

  receiverPublicKey, receiverPrivateKey, err := box.GenerateKey(rand.Reader)
  if err != nil {
    t.Error(err)
  }

  setupTestOnly(receiverPublicKey, senderPublicKey, senderPrivateKey)
  message := AuthRequest{Port: 42, Token: "hello"}
  encryptedString, err := Encrypt(message)
  if err != nil {
    t.Error(err)
  }

  setupTestOnly(senderPublicKey, receiverPublicKey, receiverPrivateKey)
  var decodedMessage AuthRequest
  err = Decrypt(encryptedString, &decodedMessage)
  if err != nil {
    t.Error(err)
  }

  if decodedMessage.Token != "hello" {
    t.Error("hello !=" + decodedMessage.Token)
  }
}
