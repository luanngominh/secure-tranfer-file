package main

import (
	"fmt"
	"net"
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
	if _, err := conn.Write([]byte(clientKey)); err != nil {
		fmt.Println("Send session key error")
		panic(err)
	}

	//Receive session key
	//Protocol use client key is 6 digits, so I use 30 bytes buffer
	sessBuffer := make([]byte, 30)
	n, err = conn.Read(sessBuffer)
	if err != nil {
		fmt.Println("Receive session key error")
		panic(err)
	}
	fmt.Println(string(sessBuffer[:n]))

	//Send file request
	fileRequest := "File: 1.jpg\n"
	if _, err := conn.Write([]byte(fileRequest)); err != nil {
		fmt.Println("Send file request error")
		panic(err)
	}

	// Begin receive file
	fmt.Println("Begin receive file")
}
