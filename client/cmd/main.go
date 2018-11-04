package main

import (
	"fmt"
	"net"

	"github.com/luanngominh/secure-tranfer-file/util"
)

func main() {
	// TODO: client get ip_server:port_server [tên file]
	conn, err := net.Dial("tcp", "localhost:1212")
	if err != nil {
		fmt.Println("Dial error")
	}

	//Receive pubkey
	//Protocol use RSA 2048 bits
	pubBuffer := make([]byte, 3000)
	n, err := conn.Read(pubBuffer)
	if err != nil {
		fmt.Println("Receive public key error")
		panic(err)
	}

	pubKey := string(pubBuffer[:n])
	fmt.Printf("Pubkey is: \n%s\n", pubKey)

	//Send client key
	//TODO: random key
	clientKey := "Key: meoconxinhxinh\n"

	//TODO: send key with public key encrypt
	if _, err := conn.Write([]byte(clientKey)); err != nil {
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
	session, _ := util.DecryptWithKey(sessBuffer[:n], []byte("meoconxinhxinh"))
	fmt.Println(string(session))

	//Send file request
	fileRequest := "File: 1.jpg\n"
	if _, err := conn.Write([]byte(fileRequest)); err != nil {
		fmt.Println("Send file request error")
		panic(err)
	}

	// Begin receive file
	fmt.Println("Begin receive file")
}
