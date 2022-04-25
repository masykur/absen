package sf3000

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"unsafe"
)

type Sf3000 struct {
	conn    *net.TCPConn
	command []byte
}

// Convert byte array to string without allocate new memory
func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
func calculateChecksum(data []byte) {
	var checksum uint16 = 0
	for _, v := range data[:len(data)-2] {
		checksum += uint16(v)
	}
	binary.LittleEndian.PutUint16(data[len(data)-2:], checksum)
}
func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func (dev *Sf3000) sendCommandParameter(parameter uint32, responsesLength ...int) ([]byte, error) {
	commandParameter := []byte{0x5a, 0xa5, dev.command[2], dev.command[3], 0, 0, 0, 0, 0, 0}
	binary.LittleEndian.PutUint32(commandParameter[4:8], parameter)
	var checksum uint16 = 0
	for _, v := range commandParameter[0:8] {
		checksum += uint16(v)
	}
	binary.LittleEndian.PutUint16(commandParameter[8:10], checksum)
	buffer := make([]byte, 0x5ff)
	// Send parameter
	dev.conn.Write(commandParameter)
	// read data or termination message
	data := make([]byte, 0)
	for _, length := range responsesLength {
		if cnt, err := dev.conn.Read(buffer[:length]); err == nil {
			if cnt == 14 && bytes.Equal(buffer[:12], []byte{0xaa, 0x55, dev.command[2], dev.command[3], 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00}) {
				break
			} else {
				data = append(data, buffer[4:cnt-2]...)
			}
		}
	}
	return data, nil
}

// send single command
func (dev *Sf3000) sendCommand(command uint16, parameter uint64, responsesLength ...int) ([]byte, error) {
	_ = dev.command[15] // early bounds check to guarantee safety of writes below
	binary.LittleEndian.PutUint16(dev.command[6:8], command)
	binary.LittleEndian.PutUint64(dev.command[8:], parameter)
	var checksum uint16 = 0
	for _, v := range dev.command[0:14] {
		checksum += uint16(v)
	}
	binary.LittleEndian.PutUint16(dev.command[14:16], checksum)

	buffer := make([]byte, 0x5ff)
	// Send authentication command
	dev.conn.Write(dev.command)
	// read response status
	if cnt, err := dev.conn.Read(buffer); err == nil {
		if bytes.Equal(buffer[:cnt-2], []byte{0x5a, 0xa5, dev.command[2], dev.command[3], 0x01, 0x00}) {
			// read data or termination message
			data := make([]byte, 0)
			for _, length := range responsesLength {
				if cnt, err := dev.conn.Read(buffer[:length]); err == nil {
					if cnt == 14 && bytes.Equal(buffer[:12], []byte{0xaa, 0x55, dev.command[2], dev.command[3], 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00}) {
						continue
					} else if len(responsesLength) == 1 && responsesLength[0] == 14 && cnt == 14 && bytes.Equal(buffer[:8], []byte{0xaa, 0x55, dev.command[2], dev.command[3], 0x00, 0x00, 0x01, 0x00}) {
						data = append(data, buffer[8:12]...)
					} else {
						data = append(data, buffer[4:cnt-2]...)
					}
				}
			}
			return data, nil
		}
		return []byte{}, fmt.Errorf("invalid response message")
	} else {
		return []byte{}, fmt.Errorf("unable to read from remote machine. %v", err)
	}
}

func isMessageValid(bytes []byte) bool {
	_ = bytes[1] // early bounds check to guarantee safety of writes below
	var checksum uint16 = 0
	for _, v := range bytes[0 : len(bytes)-2] {
		checksum += uint16(v)
	}
	return checksum == binary.LittleEndian.Uint16(bytes[len(bytes)-2:])
}

// Make connection and authenticate to machine
func (dev *Sf3000) Connect(conn *net.TCPConn, nid uint16, password uint16) (bool, error) {
	dev.command = []byte{0x55, 0xaa, 0x0, 0x0, 0x79, 0x19, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	dev.conn = conn
	binary.LittleEndian.PutUint16(dev.command[2:4], nid)
	// connect command = 0x0052
	if _, err := dev.sendCommand(0x0052, uint64(password), 14); err == nil {
		return true, nil
	}
	return false, fmt.Errorf("connection failed")
}
