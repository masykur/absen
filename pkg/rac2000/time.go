package rac2000

import (
	"encoding/hex"
	"fmt"
	"time"
)

// NOTE: this isn't multi-Unicode-codepoint aware, like specifying skintone or
//       gender of an emoji: https://unicode.org/emoji/charts/full-emoji-modifiers.html
func substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

// Obtain current date and time from machine
func (dev *Rac2000) GetDateTime() (time.Time, error) {
	// prepare command bytes array
	var command = []byte{0x19, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0xfe, 0x01, 0x01, 0xb1, 0xd6, 0x03}
	// Send command
	dev.conn.Write(command)
	reply := make([]byte, 21)
	cnt, err := dev.conn.Read(reply)
	if err != nil {
		return time.Time{}, fmt.Errorf("send command to server failed: %v", err)
	}
	if cnt != 21 {
		return time.Time{}, fmt.Errorf("invalid server reply. Expected message length is 21 but actual is %v", cnt)
	}
	replyText := hex.EncodeToString(reply)
	datePart := substr(replyText, 22, 6)
	timePart := substr(replyText, 30, 6)
	if dateTime, err := time.ParseInLocation("20060102150405", "20"+datePart+timePart, time.Local); err == nil {
		return dateTime, nil
	}
	return time.Time{}, fmt.Errorf("invalid respond message")
}

// // Set current date and time from machine
// func (dev *Rac2000) SetDateTime(t time.Time) (bool, error) {
// 	firstDate := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.Local)
// 	if t.Before(firstDate) {
// 		return false, fmt.Errorf("Invalid date, minimum value is \"2000-01-01 00:00:00\"")
// 	}
// 	// prepare command bytes array
// 	dev.prepareCommand(0x011e, 0x0004)
// 	// Send command
// 	dev.conn.Write(dev.command)
// 	reply := make([]byte, 14)
// 	cnt, err := dev.conn.Read(reply[:8])
// 	if err != nil {
// 		return false, fmt.Errorf("send command to server failed: %v", err)
// 	}
// 	if cnt != 8 {
// 		return false, fmt.Errorf("invalid server reply. Expected message length is 8 but actual is %v", cnt)
// 	}
// 	replyStatus := binary.LittleEndian.Uint16(reply[4:6])
// 	// verify reply status
// 	if replyStatus == 1 && isMessageValid(reply[:8]) {
// 		totalSeconds := uint32(t.Sub(firstDate).Seconds())
// 		buffer := []byte{0x5a, 0xa5, dev.command[2], dev.command[3], 0, 0, 0, 0, 0, 0}
// 		binary.LittleEndian.PutUint32(buffer[4:8], totalSeconds)
// 		calculateChecksum(buffer)
// 		dev.conn.Write(buffer)

// 		// read second message
// 		cnt, err = dev.conn.Read(reply)
// 		if err != nil {
// 			return false, fmt.Errorf("read server reply failed: %v", err)
// 		}
// 		if cnt != 14 {
// 			return false, fmt.Errorf("invalid server reply. Expected message length is 14 but actual is %v", cnt)
// 		}
// 		if isMessageValid(reply) {
// 			return true, nil
// 		}
// 	}
// 	return false, fmt.Errorf("invalid respond message")
// }
