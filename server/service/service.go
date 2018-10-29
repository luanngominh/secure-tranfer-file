package service

import (
	"bufio"
	"fmt"
	"net"

	"github.com/go-kit/kit/log"

	"github.com/luanngominh/secure-tranfer-file/server/config"
	"github.com/luanngominh/secure-tranfer-file/util"
)

//SessionInfo ...
type SessionInfo struct {
	ClientKey   string
	FileRequest string
	SessionKey  string
}

// ConnectionHandle handle connection client and server
func ConnectionHandle(c net.Conn, logger log.Logger) {
	defer logger.Log("Log", fmt.Sprintf("%s close connection", c.RemoteAddr().String()))
	defer c.Close()

	reader := bufio.NewReader(c)
	writer := bufio.NewWriter(c)

	logger.Log("Log", fmt.Sprintf("Begin transfer file in secure channel to %s", c.RemoteAddr().String()))

	sessInfo := &SessionInfo{
		SessionKey: util.GenerateSessionKey(),
	}

	// server send public key
	logger.Log("Log", fmt.Sprintf("Send public key to %s", c.RemoteAddr().String()))
	_, err := writer.WriteString(fmt.Sprintf("PublicKey: %s", config.Cfg.PublicKey))
	writer.Flush()
	if err != nil {
		logger.Log("Error", "Send public key to %s error", c.RemoteAddr().String())
		logger.Log("Error", err)
		return
	}

	// receive client key
	// Example receive key
	// Key: 123456\n
	logger.Log("Log", fmt.Sprintf("Receive client key from %s", c.RemoteAddr().String()))
	clientKey, err := reader.ReadString('\n')
	if err != nil {
		logger.Log("Error", fmt.Sprintf("Receive client key from %s error", c.RemoteAddr().String()))
		logger.Log("Error", err)
		return
	}
	//TODO: check valid format key
	logger.Log("Key", fmt.Sprintf("%s key is %s", c.RemoteAddr().String(), clientKey))

	// Send session key
	_, err = writer.WriteString(fmt.Sprintf("Session: %s", sessInfo.SessionKey))
	writer.Flush()
	if err != nil {
		logger.Log("Error", fmt.Sprintf("Send session key to %s error", c.RemoteAddr().String()))
		return
	}

	// Receive file request
	logger.Log("Log", fmt.Sprintf("Receive file request from %s", c.RemoteAddr().String()))
	fileRequest, err := reader.ReadString('\r')
	fmt.Printf(fileRequest)

	fmt.Println("#####")
	// if err != nil {
	// 	logger.Log("Error", fmt.Sprintf("Receive file request from %s error", c.RemoteAddr().String()))
	// 	logger.Log("Error", err)
	// 	return
	// }
	// logger.Log("File Request", fmt.Sprintf("%s want receive %s", c.RemoteAddr().String(), fileRequest))

	for {
	}
	// transfer file

}
