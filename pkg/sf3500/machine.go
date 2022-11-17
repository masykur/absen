package sf3500

import (
	"encoding/json"

	"github.com/masykur/absen/pkg/sf3500/models"
)

func (dev *Sf3500) GetDeviceInfo() (models.DeviceInfo, error) {
	command := []byte("{\"cmd\": \"GetDeviceInfo\"}")
	if response, err := dev.sendCommand(command); err == nil {
		var device models.DeviceResponse
		data := []byte(response)
		if err := json.Unmarshal(data, &device); err == nil {
			return device.ResultData, nil
		} else {
			return models.DeviceInfo{}, err
		}
	} else {
		return models.DeviceInfo{}, err
	}
}
