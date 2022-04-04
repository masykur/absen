package rac2000

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/sigurn/crc16"
)

type Rac2000 struct {
	conn      *net.TCPConn
	machineId uint16
	hComm     uint16
}

// Open connection and send handshake command to machine
func (dev *Rac2000) Connect(address string, machineId uint16, timeout time.Duration) (bool, error) {
	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.Dial("tcp", address)
	if err != nil {
		return false, err
	}
	var ok bool
	if dev.conn, ok = conn.(*net.TCPConn); ok {
		dev.hComm = 0x0101
		dev.machineId = machineId
		return true, nil
	}
	return false, fmt.Errorf("connection failed")
}

func (dev *Rac2000) Close() {
	dev.conn.Close()
}

// Calculate checksum using CRC-16/ARC algorithm
func calculateChecksum(data []byte) []byte {
	_ = data[0] // early bounds check to guarantee safety of writes below
	table := crc16.MakeTable(crc16.CRC16_ARC)
	checksum := crc16.Checksum(data, table)
	suffix := make([]byte, 3)
	binary.LittleEndian.PutUint16(suffix, checksum)
	suffix[2] = 3 // add data termination
	return append(data, suffix...)
}

func (dev *Rac2000) writeCommand(parameters ...byte) (int, error) {
	buffer := make([]byte, 2+4, 2+4+len(parameters)+3)
	buffer[0] = 0x19
	buffer[1] = 0
	binary.LittleEndian.PutUint16(buffer[2:], dev.machineId)
	buffer = append(buffer, parameters...)
	buffer = calculateChecksum(buffer)
	rv, err := dev.conn.Write(buffer)
	return rv, err
}
