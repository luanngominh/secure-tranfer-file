package main

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net"
	"os"

	"github.com/go-kit/kit/log"

	"github.com/luanngominh/secure-tranfer-file/server/config"
	"github.com/luanngominh/secure-tranfer-file/server/service"
)

var (
	errs = make(chan error)
)

func init() {
	config.Cfg.Port = os.Getenv("PORT")
	config.Cfg.Address = os.Getenv("ADDR")

	data, err := base64.StdEncoding.DecodeString(os.Getenv("PRIVATE"))
	if err != nil {
		panic(err)
	}
	config.Cfg.PrivateKey = string(data)

	data, err = base64.StdEncoding.DecodeString(os.Getenv("PUBLIC"))
	if err != nil {
		panic(err)
	}
	config.Cfg.PublicKey = string(data)

	config.Cfg.StoragePath = os.Getenv("FILE_STORAGE")

	// Create rsa.Privatekey
	block, _ := pem.Decode([]byte(config.Cfg.PrivateKey))
	config.PrivateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
}

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	listenAddr := fmt.Sprintf("%s:%s", config.Cfg.Address, config.Cfg.Port)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		panic(err)
	}
	logger.Log("Listen", fmt.Sprintf("Server listening on port %s", listenAddr))

	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Log("Connection Error", err)
			continue
		}
		logger.Log("Log", fmt.Sprintf("%s connected", conn.RemoteAddr().String()))

		// transfer file in secure channel
		go service.ConnectionHandle(conn, logger)
	}

}
