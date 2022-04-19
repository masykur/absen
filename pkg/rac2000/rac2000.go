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

const (
	beginCommand  byte = 0x19
	beginResponse byte = 0x91
	endResponse   byte = 0x03
)

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

// Send command to machine and read response message from machine
func (dev *Rac2000) sendCommand(method byte, parameters ...byte) (byte, byte, []byte, error) {
	buffer := make([]byte, 2+4+3, 2+4+4+len(parameters)+3)
	buffer[0] = 0x19
	buffer[1] = 0
	binary.LittleEndian.PutUint16(buffer[2:], dev.machineId)
	buffer[6] = byte(method)
	paramLength := byte(len(parameters))
	buffer[7] = 0xff - paramLength
	buffer[8] = paramLength
	buffer = append(buffer, parameters...)
	buffer = calculateChecksum(buffer)
	if _, err := dev.conn.Write(buffer); err == nil {
		reply := make([]byte, 0x1ff)
		// reading reply data
		// byte[0]          = start byte: 0x91
		// byte[1]          = reply sequence, the value will increase by 1 for every response
		// byte[2..5]       = machine number
		// byte[6]			= unknown, the value is always 0x00
		// byte[7]			= status code, 0x00 = success
		// byte[8]			= 0xff - byte[9]
		// byte[9]			= data length: 0x08
		// byte[10]			= unknown, some cases are equal to 1st parameter, some cases are data count
		if cnt, err := dev.conn.Read(reply); err == nil && cnt > 0 {
			message := reply[:cnt-1]
			table := crc16.MakeTable(crc16.CRC16_ARC)
			if crc16.Checksum(message, table) == 0 {
				// verify begin and termination bytes, and check CRC response
				if reply[0] == beginResponse && reply[cnt-1] == endResponse {
					status := reply[7]
					count := reply[10]
					data := reply[11 : 11+int(reply[9]-1)]
					return status, count, data, nil
				} else {
					return 0xf7, 0x00, []byte{}, fmt.Errorf("invalid server response or CRC")
				}
			} else {
				return 0xf6, 0x00, []byte{}, fmt.Errorf("invalid cyclic redudancy check")
			}
		} else {
			return 0xf5, 0x00, []byte{}, fmt.Errorf("error when reading response from server: %v", err)
		}
	} else {
		return 0xf4, 0x00, []byte{}, fmt.Errorf("send command to server failed: %v", err)
	}
}

// Convert byte array to string without allocate new memory
// func b2s(b []byte) string {
// 	return *(*string)(unsafe.Pointer(&b))
// }

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

// Convert string in byte array format to integer
func btoi(data []byte) int {
	num := 0
	for _, digit := range data {
		if digit >= 0x30 && digit < 0x3a {
			num *= 10
			num += int(digit & 0x0f)
		}
	}
	return num
}
