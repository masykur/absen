package sf3000

import (
	"encoding/binary"
	"fmt"
	"time"
)

// Obtain current date and time from machine
func (dev *Sf3000) GetDateTime() (time.Time, error) {
	// prepare command bytes array
	if response, err := dev.sendCommand(0x011d, 0x0004, 10, 14); err == nil {
		// get date and time data
		num := binary.LittleEndian.Uint32(response)
		// first date is January 1st, 2000
		firstDate := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.Local)
		// device date time is number of seconds after first date
		date := firstDate.Add(time.Second * time.Duration(num))
		return date, nil
	} else {
		return time.Time{}, fmt.Errorf("invalid respond message. %v", err)
	}
}

// Set current date and time from machine
func (dev *Sf3000) SetDateTime(t time.Time) (bool, error) {
	firstDate := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.Local)
	if t.Before(firstDate) {
		return false, fmt.Errorf("invalid date, minimum value is \"2000-01-01 00:00:00\"")
	}
	if _, err := dev.sendCommand(0x011e, 0x0004); err == nil {
		totalSeconds := uint32(t.Sub(firstDate).Seconds())
		if _, err := dev.sendCommandParameter(totalSeconds, 14); err == nil {
			return true, nil
		}
	}
	return false, fmt.Errorf("invalid respond message")
}
