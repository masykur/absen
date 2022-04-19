package rac2000

import (
	"encoding/binary"
	"fmt"
	"time"
)

type Log struct {
	Sensor           byte      `json:"Sensor"`
	Event            byte      `json:"Id"`
	DateTime         time.Time `json:"DateTime"`
	CardFacilityCode uint8     `json:"CardFacilityCode"`
	CardId           uint16    `json:"CardId"`
	DeviceId         int       `json:"DeviceId"`
	ReaderId         int       `json:"ReaderId"`
}

// Fetch chunk of log data from machine
func (dev *Rac2000) fetchLog(previousRecord byte) (byte, []Log, error) {
	// Send command
	// Command format is 13 bytes length in Little Endian byte order
	// byte[0..1]   = start byte: 0x19
	// byte[1]      = unknown, always 0x00
	// byte[2..5]   = machine number, ex. 0x010000
	// byte[6-8]    = fetch command: 0x10fe01
	// byte[9]      = previous record: 0-255
	// byte[10..11] = cyclic redundancy check, CRC-16/ARC algorithm
	// byte[12]     = termination byte: 0x03
	if status, cnt, data, err := dev.sendCommand(0x10, previousRecord); err == nil {
		list := make([]Log, 0, cnt)
		if cnt > 0 {
			// data
			//   data[0]            = record length, known values are 0x08, 0x0e.
			//   data[1..rec_len+1] = record data
			//     record[0]        = sensor
			//                        0x28 = exit button
			//                        0x81 = proximity card reader
			//     record[1]        = event code
			//                        0x00 = granted
			//                        0x06 = denied (password protedted)
			//                        0x14 = denied (unregistered card)
			//     record[2..5]     = 4 bytes (32 bits) of date and time
			//       bit[26..31]    = year
			//       bit[22..25]    = month
			//       bit[18..22]    = day
			//       bit[13..17]    = hour
			//       bit[7..12]     = minute
			//       bit[0..6]      = second
			//     record[6..13]    = 3 digits of facility code + 5 digits of card number in ASCII string encoded

			var ptr byte = 0
			for i := 0; i < int(cnt); i++ {
				recSize := data[ptr]
				rec := data[ptr+1 : ptr+1+recSize]
				dateTimeBits := binary.LittleEndian.Uint32(rec[2:6])
				year := int(dateTimeBits >> 26)
				month := int(0xF & (dateTimeBits >> 22))
				day := int(0x1F & (dateTimeBits >> 17))
				hour := int(0x1F & (dateTimeBits >> 12))
				minute := int(0x3F & (dateTimeBits >> 6))
				second := int(0x3F & dateTimeBits)
				dateTime := time.Date(2000+year, time.Month(month), day, hour, minute, second, 0, time.Local)
				var logRec Log
				if recSize == 14 {
					logRec = Log{
						Sensor:           rec[0],
						Event:            rec[1],
						DateTime:         dateTime,
						CardFacilityCode: uint8(btoi(rec[6:9])),
						CardId:           uint16(btoi(rec[9:])),
						DeviceId:         int(dev.machineId),
						ReaderId:         1} // unknown, temporary hardcoded
					list = append(list, logRec)
				}
				if recSize == 8 {
					logRec = Log{
						Sensor:           rec[0],
						Event:            rec[1],
						DateTime:         dateTime,
						CardFacilityCode: 0,
						CardId:           uint16(rec[7])<<8 + uint16(rec[6]), // unknown, temporary converted from last 2 bytes
						DeviceId:         int(dev.machineId),
						ReaderId:         1} // unknown, temporary hardcoded
					list = append(list, logRec)
				}
				ptr = ptr + recSize + 1
			}
		}
		return cnt, list, nil
	} else {
		return 0, []Log{}, fmt.Errorf("send command to server failed, status code: %v, error: %v", status, err)
	}
}

// Fetch all log data from machine
func (dev *Rac2000) FetchLog() ([]Log, error) {
	logs := make([]Log, 0)
	var previousRecord byte = 0
	for {
		if cnt, list, err := dev.fetchLog(previousRecord); err == nil {
			previousRecord = cnt
			if cnt == 0 {
				break
			}
			logs = append(logs, list...)
		} else {
			return logs, err
		}
	}
	return logs, nil
}
