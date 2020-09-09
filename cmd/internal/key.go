package main

import (
  "crypto/rand"
  "encoding/base64"
  "fmt"
  "github.com/develar/lr-backup/pkg/common"
  "golang.org/x/crypto/nacl/box"
)

// https://libsodium.gitbook.io/doc/public-key_cryptography/sealed_boxes
// https://crypto.stackexchange.com/questions/72610/nacl-allows-decryption-with-the-same-public-key-as-for-encryption
// "As a mitigation, the sender can create an ephemeral key pair, deleted right after the message has been encrypted."
// That's why sealed box is used instead of just crypto box.
func main() {
  publicKey, privateKey, err := box.GenerateKey(rand.Reader)
  if err != nil {
    panic(err)
  }

  println("public key: " + base64.RawURLEncoding.EncodeToString(publicKey[:]))
  println("private key: " + base64.RawURLEncoding.EncodeToString(privateKey[:]))

  msg := []byte("Hello")
  // This encrypts msg and appends the result to the nonce.
  encrypted, err := box.SealAnonymous(nil, msg, publicKey, rand.Reader)
  if err != nil {
    panic(err)
  }

  decrypted, ok := box.OpenAnonymous(nil, encrypted, publicKey, privateKey)
  if !ok {
    panic("decryption error")
  }
  fmt.Println(string(decrypted))
}
