package sf3500

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

const (
	PROTOCOL_KEY        uint32 = 404232216
	HEADER_SIZE         int    = 32
	RECEIVE_BUFFER_SIZE int    = 409600
)

type Sf3500 struct {
	conn *net.TCPConn
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

// send single command
func (dev *Sf3500) sendCommand(command []byte) ([]byte, error) {
	commandLength := len(command)
	buffer := make([]byte, HEADER_SIZE, HEADER_SIZE+commandLength)
	binary.LittleEndian.PutUint32(buffer[0:4], uint32(commandLength))
	binary.LittleEndian.PutUint32(buffer[4:8], PROTOCOL_KEY)
	buffer = append(buffer, command...)
	dev.conn.Write(buffer)
	responseHeader := make([]byte, HEADER_SIZE)
	var responseLength int
	if cnt, err := dev.conn.Read(responseHeader); err == nil {
		if cnt != HEADER_SIZE {
			return nil, fmt.Errorf("invalid response")
		}
		responseLength = int(binary.LittleEndian.Uint32(responseHeader))
	}
	if responseLength == 0 {
		return nil, nil
	}
	responseBuffer := make([]byte, min(responseLength, 1460))
	var stringBuffer bytes.Buffer
	counting := 0
	for {
		counting++
		if cnt, err := dev.conn.Read(responseBuffer); err == nil {
			if cnt == 0 {
				break
			}
			stringBuffer.Write(responseBuffer[:cnt])
			if cnt == 1 && responseBuffer[0] == 0 {
				break
			}
			if responseBuffer[cnt-2] == 10 && responseBuffer[cnt-1] == 0 {
				break
			}
		} else {
			return nil, nil
		}
	}
	// fmt.Println(stringBuffer.Len())
	// fmt.Println(string(stringBuffer.Bytes()[:min(80, stringBuffer.Len())]))
	// fmt.Println(string(stringBuffer.Bytes()[stringBuffer.Len()-min(80, stringBuffer.Len()):]))
	return stringBuffer.Bytes()[:stringBuffer.Len()-2], nil
}

// Open connection
func (dev *Sf3500) Connect(address string, timeout time.Duration) (bool, error) {
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

func (dev *Sf3500) Close() {
	dev.conn.Close()
}
