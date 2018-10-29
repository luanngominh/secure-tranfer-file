package service

// FileSender ...
type FileSender interface {
	Send(conn interface{})
}

// SendFile storage infomation of file and session
type SendFile struct {
	filename string
	key      string
	session  string
}

//Send file to conn
func (s *SendFile) Send(conn interface{}) error {
	return nil
}
