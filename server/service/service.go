package service

import (
	"fmt"
	"net"

	"github.com/go-kit/kit/log"
)

type Transfer interface {
}

func ConnectionHandle(c net.Conn, logger log.Logger) error {
	defer logger.Log("Log", fmt.Sprintf("%s close connection", c.RemoteAddr().String()))
	defer c.Close()
	return nil
}
