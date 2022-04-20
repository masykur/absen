package sf3000

import (
	"encoding/binary"
	"fmt"
	"time"
)

// Obtain current date and time from machine
func (dev *Sf3000) GetDateTime() (time.Time, error) {
	// prepare command bytes array
	dev.prepareCommand(0x011D, 0x0004)
	// Send command
	dev.conn.Write(dev.command)
	reply1 := make([]byte, 8)
	cnt, err := dev.conn.Read(reply1)
	if err != nil {
		return time.Time{}, fmt.Errorf("send command to server failed: %v", err)
	}
	if cnt != 8 {
		return time.Time{}, fmt.Errorf("invalid server reply. Expected message length is 8 but actual is %v", cnt)
	}
	replyStatus := binary.LittleEndian.Uint16(reply1[4:6])
	// verify reply status
	if replyStatus == 1 && isMessageValid(reply1) {
		reply2 := make([]byte, 10)
		// read second message
		cnt, err = dev.conn.Read(reply2)
		if err != nil {
			return time.Time{}, fmt.Errorf("read server reply failed: %v", err)
		}
		if cnt != 10 {
			return time.Time{}, fmt.Errorf("invalid server reply. Expected message length is 10 but actual is %v", cnt)
		}
		if isMessageValid(reply2) {
			// get date and time data
			num := binary.LittleEndian.Uint32(reply2[4:8])
			// first date is January 1st, 2000
			firstDate := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.Local)
			// device date time is number of seconds after first date
			date := firstDate.Add(time.Second * time.Duration(num))

			// verify respond status
			reply3 := make([]byte, 14)
			// read actual response contain product code
			cnt, err = dev.conn.Read(reply3)
			if err != nil {
				return time.Time{}, fmt.Errorf("read server reply failed: %v", err)
			}
			if cnt != 14 {
				return time.Time{}, fmt.Errorf("invalid server reply. Expected message length is 14 but actual is %v", cnt)
			}
			replyStatus = binary.LittleEndian.Uint16(reply3[6:8])
			if replyStatus == 1 && isMessageValid(reply3) {
				return date, nil
			}
		}
	}
	return time.Time{}, fmt.Errorf("invalid respond message")
}

// Set current date and time from machine
func (dev *Sf3000) SetDateTime(t time.Time) (bool, error) {
	firstDate := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.Local)
	if t.Before(firstDate) {
		return false, fmt.Errorf("invalid date, minimum value is \"2000-01-01 00:00:00\"")
	}
	// prepare command bytes array
	dev.prepareCommand(0x011e, 0x0004)
	// Send command
	dev.conn.Write(dev.command)
	reply := make([]byte, 14)
	cnt, err := dev.conn.Read(reply[:8])
	if err != nil {
		return false, fmt.Errorf("send command to server failed: %v", err)
	}
	if cnt != 8 {
		return false, fmt.Errorf("invalid server reply. Expected message length is 8 but actual is %v", cnt)
	}
	replyStatus := binary.LittleEndian.Uint16(reply[4:6])
	// verify reply status
	if replyStatus == 1 && isMessageValid(reply[:8]) {
		totalSeconds := uint32(t.Sub(firstDate).Seconds())
		buffer := []byte{0x5a, 0xa5, dev.command[2], dev.command[3], 0, 0, 0, 0, 0, 0}
		binary.LittleEndian.PutUint32(buffer[4:8], totalSeconds)
		calculateChecksum(buffer)
		dev.conn.Write(buffer)

		// read second message
		cnt, err = dev.conn.Read(reply)
		if err != nil {
			return false, fmt.Errorf("read server reply failed: %v", err)
		}
		if cnt != 14 {
			return false, fmt.Errorf("invalid server reply. Expected message length is 14 but actual is %v", cnt)
		}
		if isMessageValid(reply) {
			return true, nil
		}
	}
	return false, fmt.Errorf("invalid respond message")
}
