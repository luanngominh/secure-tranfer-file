package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"

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
