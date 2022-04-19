package rac2000

import (
	"fmt"
	"time"
)

// Obtain current date and time from machine
func (dev *Rac2000) GetDateTime() (time.Time, error) {
	// Command format is 13 bytes length in Little Endian byte order
	// byte[0]   	= start byte: 0x19
	// byte[1]      = unknown, always 0x00
	// byte[2..5]   = machine number, ex. 0x010000
	// byte[6]      = Read command: 0x00
	// byte[7]		= 0xff - byte[8]
	// byte[8]      = parameter length: 0x08
	// byte[9]      = parameter: 0x01
	// byte[10..11] = check sum, CRC-16/ARC algorithm
	// byte[12]		= termination byte, the value is always 0x03
	if status, _, data, err := dev.sendCommand(0x00, 0x01); err == nil {
		// 21 bytes reply in Little Endian byte order
		// byte[0]          = start byte: 0x91
		// byte[1]          = reply sequence, the value will increase by 1 for every response
		// byte[2..5]       = machine number
		// byte[6..7]       = unknown, the value is always 0x0000
		// byte[8]			= 0xff - byte[9]
		// byte[9]			= data length: 0x08
		// byte[10]			= data #1 parameter: 0x01
		// byte[11]			= 2 digits of year in binary coded decimal
		// byte[12]			= month number (1-12) in binary coded decimal
		// byte[13]			= day of month binary coded decimal
		// byte[14]			= week number in binary coded decimal
		// byte[15]			= hour (24h) in binary coded decimal
		// byte[16]			= minute in binary coded decimal
		// byte[17]			= second binary coded decimal
		// byte[18..19]   	= 2 bytes of check sum using CRC-16/ARC algorithm
		// byte[20]			= termination byte, the value is always 0x03
		year := 2000 + bcd2dec(data[0])
		month := time.Month(bcd2dec(data[1]))
		day := bcd2dec(data[2])
		hour := bcd2dec(data[4])
		minute := bcd2dec(data[5])
		second := bcd2dec(data[6])
		dateTime := time.Date(year, month, day, hour, minute, second, 0, time.Local)
		return dateTime, nil
	} else {
		return time.Time{}, fmt.Errorf("send command to server failed, status code: %v, error: %v", status, err)
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
	// Command format is 20 bytes length in Little Endian byte order
	// byte[0]   	= start byte: 0x19
	// byte[1]      = unknown, always 0x00
	// byte[2..5]   = machine number, ex. 0x010000
	// byte[6]      = Write command: 0x01
	// byte[7]		= 0xff - byte[8]
	// byte[8]      = parameter length: 0x08
	// byte[9]      = parameter: 0x01
	// byte[10]		= 2 digits of year in binary coded decimal
	// byte[11]		= month number (1-12) in binary coded decimal
	// byte[12]		= day of month binary coded decimal
	// byte[13]		= week number in binary coded decimal
	// byte[14]		= hour (24h) in binary coded decimal
	// byte[15]		= minute in binary coded decimal
	// byte[16]		= second binary coded decimal
	// byte[17..18] = check sum, CRC-16/ARC algorithm
	// byte[19]		= termination byte, the value is always 0x03
	command := []byte{0x01, year, month, day, week, hour, minute, second}
	if status, _, _, err := dev.sendCommand(0x01, command...); err == nil {
		// 14 bytes reply in Little Endian byte order
		// byte[0]          = start byte: 0x91
		// byte[1]          = reply sequence, the value will increase by 1 for every response
		// byte[2..5]       = machine number
		// byte[6..7]       = unknown, the value is always 0x0000
		// byte[8]			= 0xff - byte[9]
		// byte[9]			= data length: 0x01
		// byte[10]			= data #1 parameter: 0x01
		// byte[11..12]   	= 2 bytes of check sum using CRC-16/ARC algorithm
		// byte[13]			= termination byte, the value is always 0x03
		return true, nil
	} else {
		return false, fmt.Errorf("send command to server failed, status code: %v, error: %v", status, err)
	}
}
