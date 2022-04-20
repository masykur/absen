package sf3000

import (
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

func (dev *Sf3000) prepareCommand(command uint16, parameters ...uint16) {
	_ = dev.command[15] // early bounds check to guarantee safety of writes below
	binary.LittleEndian.PutUint16(dev.command[6:8], command)
	for i := 0; i < 4; i++ {
		if len(parameters) > i {
			binary.LittleEndian.PutUint16(dev.command[8+i*2:10+i*2], parameters[i])
		} else {
			binary.LittleEndian.PutUint16(dev.command[8+i*2:10+i*2], 0)
		}
	}
	calculateChecksum(dev.command)
	// var checksum uint16 = 0
	// for _, v := range dev.command[0:14] {
	// 	checksum += uint16(v)
	// }
	// binary.LittleEndian.PutUint16(dev.command[14:16], checksum)
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
	binary.LittleEndian.PutUint16(dev.command[2:4], nid)
	// connect command = 0x0052
	dev.prepareCommand(0x0052, password)
	//const command uint16 = 0x0052
	dev.conn = conn

	reply := make([]byte, 14)
	// Send authentication command
	dev.conn.Write(dev.command)
	cnt, err := dev.conn.Read(reply[:8])
	if err != nil {
		return false, err
	}
	if cnt != 8 {
		return false, fmt.Errorf("invalid server reply. Expected message length is 8 but actual is %v", cnt)
	}
	replyStatus := binary.LittleEndian.Uint16(reply[4:6])
	if replyStatus == 1 && isMessageValid(reply[:cnt]) {
		cnt, err = conn.Read(reply)
		if err != nil {
			return false, err
		}
		if cnt != 14 {
			return false, fmt.Errorf("invalid server reply. Expected message length is 14 but actual is %v", cnt)
		}
		replyStatus = binary.LittleEndian.Uint16(reply[6:8])
		if replyStatus == 1 && isMessageValid(reply[:cnt]) {
			return true, nil
		}
	}
	return false, fmt.Errorf("connection failed")
}
