package service

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
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
	clientKeyMess := make([]byte, 1000)
	n, err := c.Read(clientKeyMess)
	if err != nil {
		logger.Log("Error", fmt.Sprintf("Receive client key from %s error", c.RemoteAddr().String()))
		logger.Log("Error", err)
		return
	}

	data, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, config.PrivateKey, clientKeyMess[:n], []byte(""))
	if err != nil {
		logger.Log("Error", fmt.Sprintf("Decrypt key from %s error", c.RemoteAddr().String()))
		return
	}

	clientKeyString := strings.Replace(string(data), "\n", "", 1)

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

	// store session key to session info
	sessInfo.ClientKey = keyParser[1]
	logger.Log("Key", fmt.Sprintf("%s key is %s", c.RemoteAddr().String(), sessInfo.ClientKey))

	// Send session key
	sessionMessage := fmt.Sprintf("Session: %s\n", sessInfo.SessionKey)
	sessionMessageCipher, _ := util.EncryptWithKey([]byte(sessionMessage), []byte(sessInfo.ClientKey))
	_, err = writer.Write(sessionMessageCipher)
	writer.Flush()
	if err != nil {
		logger.Log("Error", fmt.Sprintf("Send session key to %s error", c.RemoteAddr().String()))
		return
	}

	//Receive file request
	//File request same
	//File: meocon.jpg
	logger.Log("Log", fmt.Sprintf("Receive file request from %s", c.RemoteAddr().String()))
	// Use 1000 bytes to store file request
	fileRequestCipher := make([]byte, 1000)
	n, err = c.Read(fileRequestCipher)
	if err != nil {
		logger.Log("Error", "Receive file request error")
		return
	}

	//decrypt filerequest with clientkey
	fileRequestData, err := util.DecryptWithKey(fileRequestCipher[:n], []byte(sessInfo.ClientKey))
	if err != nil {
		logger.Log("Error", "Decrypt file request from %s error", c.RemoteAddr().String())
	}

	fileRequestDataParser := strings.Split(string(fileRequestData), "Session: ")
	if len(fileRequestDataParser) != 2 {
		logger.Log("Error", fmt.Sprintf("%s File request message error", c.RemoteAddr().String()))
		return
	}

	// chekc session key is true
	sessionKey := strings.Replace(fileRequestDataParser[1], "\n", "", 1)
	if sessionKey != sessInfo.SessionKey {
		logger.Log("Error", fmt.Sprintf("%s session key invalid", c.RemoteAddr().String()))
		return
	}

	fileRequestString := strings.Replace(fileRequestDataParser[0], "\n", "", 1)

	//case filename is spam message
	if len(fileRequestString) < 7 {
		logger.Log("Error", fmt.Sprintf("%s File request %s error", c.RemoteAddr().String(), fileRequestString))
		return
	}

	// check file request is true
	if fileRequestString[:6] != "File: " {
		logger.Log("Error", fmt.Sprintf("%s File request %s error", c.RemoteAddr().String(), fileRequestString))
		return
	}
	logger.Log("Log", fmt.Sprintf("%s request %s file", c.RemoteAddr().String(), fileRequestString))

	// check file has existed yet
	fileRequestParses := strings.Split(fileRequestString, " ")
	if len(fileRequestParses) != 2 {
		logger.Log("Error", fmt.Sprintf("%s File request %s error due to parse error", c.RemoteAddr().String(), fileRequestString))
		return
	}

	//1 due to get filename
	// fmt.Println(fmt.Sprintf("%s/%s\n", config.Cfg.StoragePath, fileRequestParses[1]))
	fullFilePath := fmt.Sprintf("%s/%s", config.Cfg.StoragePath, fileRequestParses[1])
	if _, err := os.Stat(fullFilePath); os.IsNotExist(err) {
		logger.Log("Error", fmt.Sprintf("%s request %s file is not exist", c.RemoteAddr().String(), fileRequestParses[1]))
		return
	}
	logger.Log("Log", fmt.Sprintf("%s request %s file is valid, start send file", c.RemoteAddr().String(), fileRequestParses[1]))

	//Send file
	fileSender := &SendFile{
		Filename: fullFilePath,
		Key:      sessInfo.ClientKey,
	}

	//Send with encrypted data
	if err := fileSender.Send(c, fileSender); err != nil {
		logger.Log("Error", fmt.Sprintf("Send file to %s error: %v", c.RemoteAddr().String(), err))
	}
	logger.Log("Success", fmt.Sprintf("Send file to %s success", c.RemoteAddr().String()))

}
