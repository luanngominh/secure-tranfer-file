package service

import (
	"io/ioutil"
	"net"
)

// FileSender ...
type FileSender interface {
	Send(c net.Conn)
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

	_, err = c.Write(fileContent)
	if err != nil {
		return err
	}

	return nil
}
