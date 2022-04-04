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
	// Send command
	if _, err := dev.writeCommand(0x00, 0xfe, 0x01, 0x01); err == nil {
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
	} else {
		return time.Time{}, fmt.Errorf("send command to server failed: %v", err)
	}
}

// Convert decimal to hex number, ex. 12 decimal convert to 0x12 instead of 0x0C
func dechex(v byte) byte {
	return ((v / 10) << 4) + ((v % 10) & 0x0F)
}

// Set current date and time from machine
func (dev *Rac2000) SetDateTime(t time.Time) (bool, error) {
	year := dechex(byte(t.Year() % 100))
	month := dechex(byte(t.Month()))
	day := dechex(byte(t.Day()))
	week := byte(t.Weekday())
	hour := dechex(byte(t.Hour()))
	minute := dechex(byte(t.Minute()))
	second := dechex(byte(t.Second()))
	command := []byte{0x01, 0xf7, 0x08, 0x01, year, month, day, week, hour, minute, second}
	if _, err := dev.writeCommand(command...); err == nil {
		reply := make([]byte, 14)
		cnt, err := dev.conn.Read(reply)
		if err != nil {
			return false, fmt.Errorf("send command to server failed: %v", err)
		}
		if cnt != 14 {
			return false, fmt.Errorf("invalid server reply. Expected message length is 21 but actual is %v", cnt)
		}
		if reply[7] == 0x00 && reply[8] == 0xfe {
			return true, nil
		}
		return false, fmt.Errorf("invalid command, error code: %v", reply[7])
	} else {
		return false, fmt.Errorf("send command to server failed: %v", err)
	}
}
