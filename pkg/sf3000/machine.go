package sf3000

import (
	"encoding/binary"
	"fmt"
)

// Obtain product code from machine
func (dev *Sf3000) GetProductCode() (string, error) {
	// prepare command bytes array
	dev.prepareCommand(0x0114, 0)
	// Send command
	dev.conn.Write(dev.command)
	reply1 := make([]byte, 8)
	cnt, err := dev.conn.Read(reply1)
	if err != nil {
		return "", fmt.Errorf("send command to server failed: %v", err)
	}
	if cnt != 8 {
		return "", fmt.Errorf("invalid server reply. Expected message length is 8 but actual is %v", cnt)
	}
	replyStatus := binary.LittleEndian.Uint16(reply1[4:6])
	// verify reply status
	if replyStatus == 1 && isMessageValid(reply1) {
		reply2 := make([]byte, 14)
		// read second message
		cnt, err = dev.conn.Read(reply2)
		if err != nil {
			return "", fmt.Errorf("read server reply failed: %v", err)
		}
		if cnt != 14 {
			return "", fmt.Errorf("invalid server reply. Expected message length is 14 but actual is %v", cnt)
		}
		replyStatus = binary.LittleEndian.Uint16(reply2[6:8])
		// verify second reply status
		if replyStatus == 1 && isMessageValid(reply2) {
			reply3 := make([]byte, 38)
			// read actual response contain product code
			cnt, err = dev.conn.Read(reply3)
			if err != nil {
				return "", fmt.Errorf("read server reply failed: %v", err)
			}
			if cnt != 38 {
				return "", fmt.Errorf("invalid server reply. Expected message length is 38 but actual is %v", cnt)
			}
			// parse and return product code info
			result := reply3[4 : cnt-2]
			return b2s(result), nil
		}
	}
	return "", fmt.Errorf("invalid respond message")
}

// Obtain machine serial number
func (dev *Sf3000) GetSerialNumber() (string, error) {
	// prepare command bytes array
	dev.prepareCommand(0x0115, 0)
	// Send command
	dev.conn.Write(dev.command)
	reply1 := make([]byte, 8)
	cnt, err := dev.conn.Read(reply1)
	if err != nil {
		return "", fmt.Errorf("send command to server failed: %v", err)
	}
	if cnt != 8 {
		return "", fmt.Errorf("invalid server reply. Expected message length is 8 but actual is %v", cnt)
	}
	replyStatus := binary.LittleEndian.Uint16(reply1[4:6])
	// verify reply status
	if replyStatus == 1 && isMessageValid(reply1) {
		reply2 := make([]byte, 14)
		// read second message
		cnt, err = dev.conn.Read(reply2)
		if err != nil {
			return "", fmt.Errorf("read server reply failed: %v", err)
		}
		if cnt != 14 {
			return "", fmt.Errorf("invalid server reply. Expected message length is 14 but actual is %v", cnt)
		}
		replyStatus = binary.LittleEndian.Uint16(reply2[6:8])
		// verify second reply status
		if replyStatus == 1 && isMessageValid(reply2) {
			reply3 := make([]byte, 38)
			// read actual response contain product code
			cnt, err = dev.conn.Read(reply3)
			if err != nil {
				return "", fmt.Errorf("read server reply failed: %v", err)
			}
			if cnt != 38 {
				return "", fmt.Errorf("invalid server reply. Expected message length is 38 but actual is %v", cnt)
			}
			// parse and return product code info
			result := reply3[4 : cnt-2]
			return b2s(result), nil
		}
	}
	return "", fmt.Errorf("invalid respond message")
}
