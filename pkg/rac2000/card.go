package rac2000

import (
	"fmt"
)

type Card struct {
	FacilityCode uint8  `json:"FacilityCode"`
	Id           uint16 `json:"Id"`
	Password     uint32 `json:"Password"`
	Timezone     byte   `json:"Timezone"`
	Status       byte   `json:"Status"`
}

func (dev *Rac2000) readRamData(address int, length byte) ([]byte, error) {
	// Send command
	// Command format is 17 bytes length in Little Endian byte order
	// byte[0..1]   = start byte: 0x19
	// byte[1]      = unknown, always 0x00
	// byte[2..5]   = machine number, ex. 0x010000
	// byte[6-9]    = fetch command: 0x01fa050f
	// byte[10..12] = RAM address: begin address for registered cards is 0xa00c00
	// byte[13]		= data length: 0x80
	// byte[14..15] = check sum, CRC-16/ARC algorithm
	// byte[16]     = termination byte: 0x03
	if status, _, _, err := dev.sendCommand(0x01, 0x0f, byte(address&0xff), byte((address>>8)&0xff), byte((address>>16)&0xff), length); err == nil {
		if status, _, data, err := dev.sendCommand(0x00, 0x0f); err == nil {
			return data, nil
		} else {
			return []byte{}, fmt.Errorf("send read command to server failed, status code: %v, error: %v", status, err)
		}
	} else {
		return []byte{}, fmt.Errorf("send command to server failed, status code: %v, error: %v", status, err)
	}
}

// Retrieve list of registerd cards in the machine
func (dev *Rac2000) GetCards() ([]Card, error) {
	// Registered card is saved in the RAM at address 0x000ca0
	// Each card information is requires 16 bytes length
	// Card numbuer (including 3 digits facility codes prefix)
	// is stored start from first byte in ASCII terminated by colon (:)
	allcards := make([]Card, 0)
	address := 0x000ca0
	var dataLength byte = 0x80
	for {
		// raed RAM data by 128 bytes chunk
		if data, err := dev.readRamData(address, dataLength); err == nil {
			cards := make([]Card, 0, len(data)/16)
			for i := 0; i < len(data); i += 16 {
				numbers := data[i : i+16]
				if numbers[0] == 0xff {
					break
				}

				num := 0
				numPart := 0
				pass := 0
				for _, digit := range numbers[:14] {
					switch numPart {
					case 0:
						if digit >= 0x30 && digit < 0x3a {
							num *= 10
							num += int(digit & 0x0f)
						}
						if digit == 0x3a {
							numPart = 1
						}
					case 1:
						if digit == 0xff {
							break
						}
						if digit&0x0f == 0x0f {
							pass *= 10
							pass += int(digit >> 4)
						} else {
							pass *= 10
							pass += int(digit >> 4)
							pass *= 10
							pass += int(digit & 0x0f)
						}
					}
				}
				if num > 0 {
					cards = append(cards, Card{
						FacilityCode: uint8(num / 100000),
						Id:           uint16(num % 100000),
						Password:     uint32(pass),
						Timezone:     numbers[14],
						Status:       numbers[15]})
				}
			}
			if len(cards) > 0 {
				allcards = append(allcards, cards...)
			}
			if len(cards) < 8 {
				return allcards, nil
			} else {
				address += int(dataLength)
			}
		} else {
			return []Card{}, err
		}
	}
}

// Register new card to machine
func (dev *Rac2000) AddCard(card Card) (bool, error) {
	// Send command
	// Command format is 29 bytes length in Little Endian byte order
	// byte[0..1]   = start byte: 0x19
	// byte[1]      = unknown, always 0x00
	// byte[2..5]   = machine number, ex. 0x010000
	// byte[6]      = Write command: 0x01
	// byte[7]		= 0xff - byte[8]
	// byte[8]      = parameter length: 0x11
	// byte[9]      = parameter: 0x05
	// byte[10..23] = Card number (including facility code) + password in ASCII format, separated by colon (:)
	// byte[24] 	= Timezone
	// byte[25]		= Status
	// byte[26..27] = check sum, CRC-16/ARC algorithm
	// byte[28]     = termination byte: 0x03
	//command = append(command, []byte{0x01, 0xee, 0x11, 0x05}...)
	command := []byte{0x05, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, card.Timezone, byte(card.Status)}
	var number int
	// format facility code in ASCII
	number = int(card.FacilityCode)
	for i := 0; i < 3; i++ {
		command[3-i] = 0x30 | byte(number%10)
		number /= 10
	}
	// format card id in ASCII
	number = int(card.Id)
	for i := 0; i < 5; i++ {
		command[8-i] = 0x30 | byte(number%10)
		number /= 10
	}
	command[9] = 0x3a // separator
	// format password in binary coded decimal
	numbers := make([]byte, 0)
	number = int(card.Password)
	for number > 0 {
		numbers = append([]byte{0x30 | byte(number%10)}, numbers...)
		number /= 10
	}
	for i, num := range numbers {
		if i%2 == 0 {
			command[10+(i/2)] = command[10+(i/2)] & ((num << 4) | 0x0f)
		} else {
			command[10+(i/2)] = command[10+(i/2)] & (num | 0xf0)
		}
	}
	// send command to machine
	if status, _, _, err := dev.sendCommand(0x01, command...); err == nil {
		return true, nil
	} else {
		return false, fmt.Errorf("send command to server failed, status code: %v, error: %v", status, err)
	}
}

// Register new card to machine
func (dev *Rac2000) DelCard(cardFacilityCode uint8, cardId uint16) (bool, error) {
	// Send command
	// Command format is 29 bytes length in Little Endian byte order
	// byte[0..1]   = start byte: 0x19
	// byte[1]      = unknown, always 0x00
	// byte[2..5]   = machine number, ex. 0x010000
	// byte[6]      = 0x01
	// byte[7]		= 0xff - byte[8]
	// byte[8]      = command parameter length
	// byte[9]      = command code: 0x06
	// byte[10..23] = Card number (including facility code)
	// byte[24] 	= Timezone
	// byte[25]		= Status
	// byte[26..27] = check sum, CRC-16/ARC algorithm
	// byte[28]     = termination byte: 0x03
	command := []byte{0x06, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30}

	var number int
	// format facility code in ASCII
	number = int(cardFacilityCode)
	for i := 0; i < 3; i++ {
		command[3-i] = 0x30 | byte(number%10)
		number /= 10
	}
	// format card id in ASCII
	number = int(cardId)
	for i := 0; i < 5; i++ {
		command[8-i] = 0x30 | byte(number%10)
		number /= 10
	}
	// send command to machine
	if status, _, _, err := dev.sendCommand(0x01, command...); err == nil {
		return true, nil
	} else {
		return false, fmt.Errorf("send command to server failed, status code: %v, error: %v", status, err)
	}
}
