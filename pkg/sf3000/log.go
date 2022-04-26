package sf3000

import (
	"encoding/binary"
	"fmt"
	"time"
)

type SensorType uint8
type Mode uint8
type FunctionKey uint8

const (
	Keypad      SensorType = 1 << iota // 1
	Card                               // 2
	Fingerprint                        // 4
)

const (
	AnyMode                                  Mode = iota // 0
	FingerPrintMode                                      // 1
	CardOrFingerPrintMode                                // 2
	IdAndFingerPrintOrCardMode                           // 3
	IdAndFingerPrintOrIdAndCardMode                      // 4
	IdAndFingerPrintOrCardAndFingerPrintMode             // 5
	OpenMode                                             // 6
	CloseMode                                            // 7
	CardMode                                             // 8
	IdOrFingerPrintMode                                  // 9
	IdOrCardMode                                         // 10
	IdAndCardMode                                        // 11
	CardAndFingerPrintMode                               // 12
	IdAndFingerPrintMode                                 // 13
	IdAndCardAndFingerPrintMode                          // 14
)
const (
	F1    FunctionKey = iota // 0
	F2                       // 1
	F3                       // 2
	F4                       // 3
	NoKey                    // 4
)

func (e SensorType) String() string {
	switch e {
	case Keypad:
		return "Keypad"
	case Card:
		return "Card"
	case Fingerprint:
		return "Fingerprint"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}

func (e Mode) String() string {
	switch e {
	case AnyMode:
		return "Any"
	case FingerPrintMode:
		return "Fingerprint"
	case CardOrFingerPrintMode:
		return "Card or fingerprint"
	case IdAndFingerPrintOrCardMode:
		return "ID, Card or fingerprint"
	case IdAndFingerPrintOrIdAndCardMode:
		return "ID and fingerprint or ID and card"
	case IdAndFingerPrintOrCardAndFingerPrintMode:
		return "ID and fingerprint or Card and fingerprint"
	case OpenMode:
		return "Open"
	case CloseMode:
		return "Close"
	case CardMode:
		return "Card"
	case IdOrFingerPrintMode:
		return "ID or fingerprint"
	case IdOrCardMode:
		return "ID or card"
	case IdAndCardMode:
		return "ID and card"
	case CardAndFingerPrintMode:
		return "card and fingerprint"
	case IdAndFingerPrintMode:
		return "ID and fingerprint"
	case IdAndCardAndFingerPrintMode:
		return "ID and card and fingerprint"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}

func (e FunctionKey) String() string {
	switch e {
	case F1:
		return "F1"
	case F2:
		return "F2"
	case F3:
		return "F3"
	case F4:
		return "F4"
	case NoKey:
		return "No key"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}

type Log struct {
	UserID         int32       `json:"UserID"`
	Event          byte        `json:"Event"`
	DateTime       time.Time   `json:"DateTime"`
	UserType       uint8       `json:"UserType"`
	SensorType     SensorType  `json:"SensorType"`
	Mode           Mode        `json:"Mode"`
	FunctionKey    FunctionKey `json:"FunctionKey"`
	FunctionNumber uint8       `json:"FunctionNumber"`
}

// Fetch log data from machine
func (dev *Sf3000) FetchAllLogs() (int, []Log, error) {
	if _, err := dev.sendCommand(0x0111, 0x00, 14); err == nil {
		if response, err := dev.sendCommand(0x010f, 0x00, 14); err == nil {
			count := binary.LittleEndian.Uint32(response)
			size := int(count) * 12
			size += int((count+84)/85) * 6 // roundup
			raw := make([]byte, 0, size)
			for len(raw) < size {
				buffer := make([]byte, 1460)
				if cnt, err := dev.conn.Read(buffer); err == nil {
					raw = append(raw, buffer[:cnt]...)
				}
			}
			logs := make([]Log, 0, count)
			for i := 0; i < size; i += 4 + 85*12 + 2 {
				rec := raw[4+i : min(i+(4+85*12+2), size)-2]
				for i := 0; i < len(rec); i += 12 {
					totalSeconds := binary.LittleEndian.Uint32(rec[i:])
					firstDate := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.Local)
					date := firstDate.Add(time.Second * time.Duration(totalSeconds))
					userID := int32(binary.LittleEndian.Uint32(rec[4+i:]))
					prop := binary.LittleEndian.Uint32(rec[8+i:])
					priv := (prop & 0x01)
					sensor := (prop >> 1) & 0xff
					mode := (prop >> 9) & 0x0f
					fk := (prop >> 13) & 0xff
					logRecord := Log{
						UserID:         userID,
						Event:          0,
						DateTime:       date,
						UserType:       uint8(priv),
						SensorType:     SensorType(sensor),
						Mode:           Mode(mode),
						FunctionKey:    FunctionKey(fk / 10),
						FunctionNumber: uint8(fk % 10)}
					logs = append(logs, logRecord)

					//break
				}
			}
			return int(count), logs, nil
		} else {
			return 0, []Log{}, fmt.Errorf("invalid respond message. %v", err)
		}
	} else {
		return 0, []Log{}, fmt.Errorf("invalid respond message. %v", err)
	}
}
