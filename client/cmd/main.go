package main

import (
	"fmt"
	"net"
)

func main() {
	// TODO: client get ip_server:port_server [tên file]
	conn, err := net.Dial("tcp", "localhost:60000")
	if err != nil {
		fmt.Println("Dial error")
	}

	fmt.Fprintf(conn, "Key: dfdf\n")
	fmt.Fprintf(conn, "File: nancy.jp\n")

}
