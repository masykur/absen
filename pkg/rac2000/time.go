package rac2000

import (
	"fmt"
	"time"
)

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
		// decode date and time data from Binary Coded Decimal (BCD) to datetime
		year := 2000 + bcd2dec(reply[11])
		month := time.Month(bcd2dec(reply[12]))
		day := bcd2dec(reply[13])
		hour := bcd2dec(reply[15])
		minute := bcd2dec(reply[16])
		second := bcd2dec(reply[17])
		dateTime := time.Date(year, month, day, hour, minute, second, 0, time.Local)
		return dateTime, nil
	} else {
		return time.Time{}, fmt.Errorf("send command to server failed: %v", err)
	}
}

// Set date and time to machine
func (dev *Rac2000) SetDateTime(t time.Time) (bool, error) {
	// convert decimal of each datetime parts to binary coded decimal
	year := dec2bcd(t.Year() % 100)
	month := dec2bcd(int(t.Month()))
	day := dec2bcd(t.Day())
	week := byte(t.Weekday())
	hour := dec2bcd(t.Hour())
	minute := dec2bcd(t.Minute())
	second := dec2bcd(t.Second())
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
