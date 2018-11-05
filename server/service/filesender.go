package service

import (
	"io/ioutil"
	"net"

	"github.com/luanngominh/secure-tranfer-file/util"
)

// FileSender ...
type FileSender interface {
	Send(c net.Conn, fileInfo *SendFile) error
}

// SendFile storage infomation of file and session
type SendFile struct {
	Filename string
	Key      string
}

//Send file to conn
func (s *SendFile) Send(c net.Conn, fileInfo *SendFile) error {
	fileContent, err := ioutil.ReadFile(fileInfo.Filename)
	if err != nil {
		return err
	}

	cipherContent, err := util.EncryptWithKey(fileContent, []byte(fileInfo.Key))

	_, err = c.Write(cipherContent)
	if err != nil {
		return err
	}

	return nil
}
