// use for test algorithm
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

var (
	pub  string
	priv string
)

func init() {
	data, err := base64.StdEncoding.DecodeString(os.Getenv("PRIVATE"))
	if err != nil {
		panic(err)
	}
	priv = string(data)

	data, err = base64.StdEncoding.DecodeString(os.Getenv("PUBLIC"))
	if err != nil {
		panic(err)
	}
	pub = string(data)
}

func main() {
	fmt.Println(pub)
	fmt.Println(priv)

	block, _ := pem.Decode([]byte(pub))
	if block == nil || block.Type != "PUBLIC KEY" {
		log.Fatal("fail to decode PEM block")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	mess := "meo con xinh xinh"
	cirpher, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey.(*rsa.PublicKey), []byte(mess), []byte(""))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(cirpher))

	block, _ = pem.Decode([]byte(priv))
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	data, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, cirpher, []byte(""))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}
