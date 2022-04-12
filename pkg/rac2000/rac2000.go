package rac2000

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
	"unsafe"

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

// Convert byte array to string without allocate new memory
func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// Binary Coded Decimal (BCD) to decimal converter.
// Binary coded decimal is not the same as hexadecimal.
// Whereas a 4-bit hexadecimal number is valid up to F16 representing binary 11112, (decimal 15),
// binary coded decimal numbers stop at 9 binary 10012.
// This means that although 16 numbers (24) can be represented using four binary digits,
// in the BCD numbering system the six binary code combinations of:
// 1010 (decimal 10), 1011 (decimal 11), 1100 (decimal 12), 1101 (decimal 13), 1110 (decimal 14),
// and 1111 (decimal 15) are classed as forbidden numbers and can not be used.
func bcd2dec(b byte) int {
	// convert using very simple method without error checking
	return int(b - 6*(b>>4))
}

// Decimal (base10) to Binary CodedDecimal (BCD) converter.
// Binary coded decimal is not the same as hexadecimal.
// Whereas a 4-bit hexadecimal number is valid up to F16 representing binary 11112, (decimal 15),
// binary coded decimal numbers stop at 9 binary 10012.
// This means that although 16 numbers (24) can be represented using four binary digits,
// in the BCD numbering system the six binary code combinations of:
// 1010 (decimal 10), 1011 (decimal 11), 1100 (decimal 12), 1101 (decimal 13), 1110 (decimal 14),
// and 1111 (decimal 15) are classed as forbidden numbers and can not be used.
func dec2bcd(d int) byte {
	// convert using very simple method without error checking
	return byte(d + (d / 10 * 6))
}
