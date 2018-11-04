package service

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

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

	// receive client key with 1000 bytes buffer
	//TODO: Receive with cirpher text, decrypt it
	clientKey := make([]byte, 1000)
	n, err := c.Read(clientKey)
	if err != nil {
		logger.Log("Error", fmt.Sprintf("Receive client key from %s error", c.RemoteAddr().String()))
		logger.Log("Error", err)
		return
	}

	//n-1 to remove \n character
	clientKeyString := string(clientKey[:n-1])

	//5 to get string "Key: ", so I use "< 6" to except case lester than 6
	//Case spam case
	if len(clientKeyString) < 6 {
		logger.Log("KeyInvalid", fmt.Sprintf("%s key invalid is %s", c.RemoteAddr().String(), clientKeyString))
		return
	}

	//Check start with "Key: "
	if clientKeyString[:5] != "Key: " {
		logger.Log("KeyInvalid", fmt.Sprintf("%s key invalid is %s", c.RemoteAddr().String(), clientKeyString))
		return
	}

	//Get client key from key response
	keyParser := strings.Split(clientKeyString, " ")
	if len(keyParser) != 2 {
		logger.Log("KeyInvalid", fmt.Sprintf("%s key invalid is %s", c.RemoteAddr().String(), clientKeyString))
		return
	}

	logger.Log("Key", fmt.Sprintf("%s key is %s", c.RemoteAddr().String(), clientKeyString))

	// store session key to session info
	sessInfo.ClientKey = keyParser[1]
	fmt.Println(keyParser[1])

	// TODO: send session key with encrypt with client key
	// Send session key
	sessionMessage := fmt.Sprintf("Session: %s\n", sessInfo.SessionKey)
	sessionMessageCipher, _ := util.EncryptWithKey([]byte(sessionMessage), []byte(sessInfo.ClientKey))
	fmt.Println(sessionMessageCipher)
	_, err = writer.Write(sessionMessageCipher)
	writer.Flush()
	if err != nil {
		logger.Log("Error", fmt.Sprintf("Send session key to %s error", c.RemoteAddr().String()))
		return
	}

	// Receive file request
	//File request same
	//File: meocon.jpg
	logger.Log("Log", fmt.Sprintf("Receive file request from %s", c.RemoteAddr().String()))
	// Use 1000 bytes to store file request
	fileRequest := make([]byte, 1000)
	n, err = c.Read(fileRequest)
	if err != nil {
		logger.Log("Key", "Receive file request error")
		return
	}
	// TODO: decrypt filerequest with clientkey

	//n - 1 because remove \n character
	fileRequestString := string(fileRequest[:n-1])

	//case filename is spam message
	if len(fileRequestString) < 7 {
		logger.Log("Log", fmt.Sprintf("%s File request %s error", c.RemoteAddr().String(), fileRequestString))
		return
	}

	// check file request is true
	if fileRequestString[:6] != "File: " {
		logger.Log("Log", fmt.Sprintf("%s File request %s error", c.RemoteAddr().String(), fileRequestString))
		return
	}
	logger.Log("Log", fmt.Sprintf("%s request %s file", c.RemoteAddr().String(), fileRequestString))

	// check file has existed yet
	fileRequestParses := strings.Split(fileRequestString, " ")
	if len(fileRequestParses) != 2 {
		logger.Log("Log", fmt.Sprintf("%s File request %s error due to parse error", c.RemoteAddr().String(), fileRequestString))
		return
	}

	//1 due to get filename
	fmt.Println(fmt.Sprintf("%s/%s\n", config.Cfg.StoragePath, fileRequestParses[1]))
	if _, err := os.Stat(fmt.Sprintf("%s/%s", config.Cfg.StoragePath, fileRequestParses[1])); os.IsNotExist(err) {
		logger.Log("Log", fmt.Sprintf("%s request %s file is not exist", c.RemoteAddr().String(), fileRequestParses[1]))
		return
	}
	logger.Log("Log", fmt.Sprintf("%s request %s file is valid, start send file", c.RemoteAddr().String(), fileRequestParses[1]))

}
