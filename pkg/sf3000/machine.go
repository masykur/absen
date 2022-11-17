package sf3000

import (
	"fmt"
)

// Obtain product code from machine
func (dev *Sf3000) GetProductCode() (string, error) {
	if response, err := dev.sendCommand(0x0114, 0x0, 14, 38); err == nil {
		return b2s(response), nil
	}
	return "", fmt.Errorf("invalid respond message")
}

// Obtain machine serial number
func (dev *Sf3000) GetSerialNumber() (string, error) {
	if response, err := dev.sendCommand(0x0115, 0x0, 14, 38); err == nil {
		return b2s(response), nil
	}
	return "", fmt.Errorf("invalid respond message")
}
