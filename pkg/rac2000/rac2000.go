package rac2000

import (
	"fmt"
	"net"
	"time"
)

type Rac2000 struct {
	conn *net.TCPConn
}

// Open connection and send handshake command to machine
func (dev *Rac2000) Connect(address string, timeout time.Duration) (bool, error) {
	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.Dial("tcp", address)
	if err != nil {
		return false, err
	}
	var ok bool
	if dev.conn, ok = conn.(*net.TCPConn); ok {
		return true, nil
	}
	return false, fmt.Errorf("connection failed")
}

func (dev *Rac2000) Close() {
	dev.conn.Close()
}
