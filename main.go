// use for test algorithm
package main

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
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
	// fmt.Println(pub)
	// fmt.Println(priv)

	privPem, _ := pem.Decode([]byte(priv))
	if privPem.Type != "RSA PRIVATE KEY" {
		fmt.Println("RSA private key error")
		return
	}
	key, _ := x509.ParsePKCS1PrivateKey(privPem.Bytes)

}
