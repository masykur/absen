package rac2000

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"time"
)

type Log struct {
	Event            int       `json:"Id"`
	DateTime         time.Time `json:"DateTime"`
	CardFacilityCode uint8     `json:"CardFacilityCode"`
	CardId           uint16    `json:"CardId"`
	DeviceId         int       `json:"DeviceId"`
	ReaderId         int       `json:"ReaderId"`
}

// Obtain current date and time from machine
func (dev *Rac2000) FetchLog(previousRecord byte) ([]Log, error) {
	// Send command
	// Command format is 13 bytes length in Little Endian byte order
	// byte[0..1]   = start byte: 0x19
	// byte[1]      = unknown, always 0x00
	// byte[2..5]   = machine number, ex. 0x010000
	// byte[6-8]    = fetch command: 0x10fe01
	// byte[9]      = previous record: 0-255
	// byte[10..11] = check sum, CRC-16/ARC algorithm
	// byte[12]     = termination byte: 0x03
	if _, err := dev.writeCommand(0x10, 0xfe, 0x01, previousRecord); err == nil {
		// prepare reply buffer
		reply := make([]byte, 4096)
		// read reply data
		cnt, err := dev.conn.Read(reply)
		if err != nil {
			return []Log{}, fmt.Errorf("send command to server failed: %v", err)
		}
		// reply format. Multiple bytes is in Little Endian byte order
		// byte[0..10]          = message header as below
		//   header[0]          = start byte: 0x91
		//   header[1..2]       = 2 bytes of reply sequence, the value will increase by 1 for every response
		//   header[3-6]        = unknown, the value always 0x00000000
		//   header[7]          = unknown, the value always 0x00000000
		//   header[8..9]       = unknown, random number
		//   header[10]         = record count
		// byte[11..end-3]      = message data as below
		//   data[0]            = record length, known values are 0x08, 0x0e.
		//   data[1..rec_len+1] = record data
		//     record[0]        = unklown, known values are 28, 81
		//     record[1]        = event code
		//     record[2..5]     = 4 bytes (32 bits) of date and time
		//       bit[26..31]    = year
		//       bit[22..25]    = month
		//       bit[18..22]    = day
		//       bit[13..17]    = hour
		//       bit[7..12]     = minute
		//       bit[0..6]      = second
		//     record[6..13]    = 3 digits of facility code + 5 digits of card number in ASCII string encoded
		// byte[end-2..end-1]   = check sum, CRC-16/ARC algorithm

		if cnt > 11 {
			recCount := reply[10]
			list := make([]Log, 0, recCount)
			var ptr byte = 11
			for i := 0; i < int(recCount); i++ {
				recSize := reply[ptr]
				rec := reply[ptr+1 : ptr+1+recSize]
				dateTimeBits := binary.LittleEndian.Uint32(rec[2:6])
				year := int(dateTimeBits >> 26)
				month := int(0xF & (dateTimeBits >> 22))
				day := int(0x1F & (dateTimeBits >> 17))
				hour := int(0x1F & (dateTimeBits >> 12))
				minute := int(0x3F & (dateTimeBits >> 6))
				second := int(0x3F & dateTimeBits)
				dateTime := time.Date(2000+year, time.Month(month), day, hour, minute, second, 0, time.Local)
				var cardFacilityCode uint8 = 0
				var cardNumber int = 0
				var logRec Log
				fmt.Println(recSize, dateTime, rec)
				if recSize == 14 {
					if num, err := strconv.Atoi(b2s(rec[6:9])); err == nil {
						cardFacilityCode = uint8(num)
					}
					if num, err := strconv.Atoi(b2s(rec[9:])); err == nil {
						cardNumber = num
					}
					logRec = Log{
						Event:            int(rec[1]),
						DateTime:         dateTime,
						CardFacilityCode: cardFacilityCode,
						CardId:           uint16(cardNumber),
						DeviceId:         int(dev.machineId),
						ReaderId:         int(rec[0] & 0x0F)}
					list = append(list, logRec)
				}
				if recSize == 8 {
					logRec = Log{
						Event:            int(rec[6])<<8 + int(rec[7]),
						DateTime:         dateTime,
						CardFacilityCode: 0,
						CardId:           0,
						DeviceId:         int(dev.machineId),
						ReaderId:         1}
					list = append(list, logRec)
				}
				ptr = ptr + recSize + 1
			}
			return list, nil
		}
		return []Log{}, nil
	} else {
		return []Log{}, fmt.Errorf("send command to server failed: %v", err)
	}
}
