package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"github.com/luanngominh/secure-tranfer-file/util"
)

func main() {
	// TODO: client get ip_server:port_server [tên file]
	conn, err := net.Dial("tcp", "localhost:1212")
	if err != nil {
		fmt.Println("Dial error")
	}

	defer conn.Close()

	//Receive pubkey
	//Protocol use RSA 2048 bits
	pubBuffer := make([]byte, 3000)
	n, err := conn.Read(pubBuffer)
	if err != nil {
		fmt.Println("Receive public key error")
		panic(err)
	}

	pubKeyMess := string(pubBuffer[:n])
	pubKey := strings.Split(pubKeyMess, ": ")[1]

	//create public key encrypt
	block, _ := pem.Decode([]byte(pubKey))
	if block == nil || block.Type != "PUBLIC KEY" {
		log.Fatal("fail to decode PEM block")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	//TODO: random key
	clientKey := "Key: meoconxinhxinh\n"
	cirpher, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey.(*rsa.PublicKey), []byte(clientKey), []byte(""))

	//Send client key encrypted
	if _, err := conn.Write([]byte(cirpher)); err != nil {
		fmt.Println("Send session key error")
		panic(err)
	}

	//Receive session key
	sessBuffer := make([]byte, 1000)
	n, err = conn.Read(sessBuffer)
	if err != nil {
		fmt.Println("Receive session key error")
		panic(err)
	}

	session, err := util.DecryptWithKey(sessBuffer[:n], []byte("meoconxinhxinh"))
	if err != nil {
		fmt.Println("Dycrypt data error")
		panic(err)
	}
	fmt.Println(string(session))

	//Send file request
	fileRequestMess := []byte(fmt.Sprintf("File: 1.jpg\n%s", session))
	fileRequestCipher, err := util.EncryptWithKey(fileRequestMess, []byte("meoconxinhxinh"))
	if _, err := conn.Write([]byte(fileRequestCipher)); err != nil {
		fmt.Println("Send file request error")
		panic(err)
	}

	// Begin receive file
	fmt.Println("Begin receive file")

	fo, err := os.Create("1.jpg")
	if err != nil {
		fmt.Println("Create file error")
		panic(err)
	}

	defer fo.Close()

	resultWithCirpherData, err := ioutil.ReadAll(conn)
	if err != nil {
		fmt.Println("Receive file error")
		panic(err)
	}

	data, err := util.DecryptWithKey(resultWithCirpherData, []byte("meoconxinhxinh"))
	if err != nil {
		fmt.Println("Decrypt key error")
		panic(err)
	}

	io.Copy(fo, bytes.NewReader(data))

	fmt.Println("Receive file complete")
}
