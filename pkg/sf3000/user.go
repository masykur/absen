package sf3000

import (
	"encoding/binary"
	"fmt"
)

type Level int16

// const (
// 	userLevel   level = 0
// 	masterLevel level = 1
// )

type Sensor int16

// const (
// 	fingerPrint1 sensor = 1
// 	fingerPrint2 sensor = 2
// 	card         sensor = 8
// )

type User struct {
	Id               int    `json:"Id"`
	Level            Level  `json:"Level"`
	Sensor           Sensor `json:"Sensor"`
	CardFacilityCode uint8  `json:"CardFacilityCode"`
	CardId           uint16 `json:"CardId"`
	Fingerprint1     []byte `json:"Fingerprint1,omitempty"` // Template format=SmartBio, Size=1404 bytes little endian, Image Dimension=256x256 (403dpi)
	Fingerprint2     []byte `json:"Fingerprint2,omitempty"` // Template format=SmartBio, Size=1404 bytes little endian, Image Dimension=256x256 (403dpi)
}

// Obtain number of user registered to machine
func (dev *Sf3000) GetUserCount() (int, error) {
	// prepare command bytes array
	dev.prepareCommand(0x0116, 0x0000, 0x0000, 0x0001)
	// Send command
	dev.conn.Write(dev.command)
	reply1 := make([]byte, 8)
	cnt, err := dev.conn.Read(reply1)
	if err != nil {
		return 0, fmt.Errorf("send command to server failed: %v", err)
	}
	if cnt != 8 {
		return 0, fmt.Errorf("invalid server reply. Expected message length is 8 but actual is %v", cnt)
	}
	replyStatus := binary.LittleEndian.Uint16(reply1[4:6])
	// verify reply status
	if replyStatus == 1 && isMessageValid(reply1[:cnt]) {
		reply2 := make([]byte, 14)
		// read second message
		cnt, err = dev.conn.Read(reply2)
		if err != nil {
			return 0, fmt.Errorf("read server reply failed: %v", err)
		}
		if cnt != 14 {
			return 0, fmt.Errorf("invalid server reply. Expected message length is 14 but actual is %v", cnt)
		}
		replyStatus = binary.LittleEndian.Uint16(reply2[6:8])
		// verify second reply status
		if replyStatus == 1 && isMessageValid(reply2[:cnt]) {
			// get user count data
			num := binary.LittleEndian.Uint32(reply2[8:12])
			// first date is January 1st, 2000
			return int(num), nil
		}
	}
	return 0, fmt.Errorf("invalid respond message")
}

// Get list of users from machine
func (dev *Sf3000) GetUsers() ([]User, error) {
	// prepare command bytes array
	dev.prepareCommand(0x0109, 0x0)
	// Send command
	dev.conn.Write(dev.command)
	reply1 := make([]byte, 8)
	cnt, err := dev.conn.Read(reply1)
	if err != nil {
		return []User{}, fmt.Errorf("send command to server failed: %v", err)
	}
	if cnt != 8 {
		return []User{}, fmt.Errorf("invalid server reply. Expected message length is 8 but actual is %v", cnt)
	}
	replyStatus := binary.LittleEndian.Uint16(reply1[4:6])
	// verify reply status
	if replyStatus == 1 && isMessageValid(reply1[:cnt]) {
		// read second message
		reply2 := make([]byte, 14)
		cnt, err = dev.conn.Read(reply2)
		if err != nil {
			return []User{}, fmt.Errorf("read server reply failed: %v", err)
		}
		if cnt != 14 {
			return []User{}, fmt.Errorf("invalid server reply. Expected message length is 14 but actual is %v", cnt)
		}
		replyStatus = binary.LittleEndian.Uint16(reply2[6:8])
		// verify second reply status
		if replyStatus == 1 && isMessageValid(reply2[:cnt]) {
			// get user count data
			dataLength := binary.LittleEndian.Uint32(reply2[8:12])
			reply3 := make([]byte, 1026) // max package size is 1026 bytes
			data := make([]byte, 0, dataLength*8)
			for {
				cnt, err = dev.conn.Read(reply3)
				if err != nil {
					return []User{}, fmt.Errorf("read server reply failed: %v", err)
				}
				data = append(data, reply3[4:cnt-2]...)
				if cnt == 0 || len(data) == cap(data) {
					break
				}
			}

			users := make([]User, 0, dataLength)
			for i := uint32(0); i < dataLength; i++ {
				uId := binary.LittleEndian.Uint32(data[i*8 : 4+i*8])
				uLevel := data[4+i*8]
				uSensor := data[5+i*8]
				uCardId := binary.LittleEndian.Uint16(data[6+i*8 : 8+i*8])
				users = append(users, User{
					Id:     int(uId),
					Level:  Level(uLevel),
					Sensor: Sensor(uSensor),
					CardId: uCardId})
			}
			return users, nil
		}
	}
	return []User{}, fmt.Errorf("invalid respond message")
}

func (dev *Sf3000) GetEnrollData(userId int) (User, error) {
	// prepare command bytes array
	const (
		fingerPrintSize int = 1404 + 12
		enrollDataSize  int = (4*8 + fingerPrintSize*2)
		bufferLength    int = 1460 + 1422
	)
	userIds := make([]byte, 4)
	binary.LittleEndian.PutUint32(userIds, uint32(userId))
	dev.prepareCommand(0x0103, uint16(userId&0xffff), uint16(userId>>16))
	// Send command
	dev.conn.Write(dev.command)
	reply1 := make([]byte, 8)
	cnt, err := dev.conn.Read(reply1)
	if err != nil {
		return User{}, fmt.Errorf("send command to server failed: %v", err)
	}
	if cnt != 8 {
		return User{}, fmt.Errorf("invalid server reply. Expected message length is 8 but actual is %v", cnt)
	}
	// verify reply status
	if isMessageValid(reply1[:cnt]) {
		// read second message
		buffer := make([]byte, bufferLength)
		data := make([]byte, 0, enrollDataSize)
		index := 0
		for index < bufferLength {
			if cnt, err = dev.conn.Read(buffer[index:]); err == nil {
				index += cnt
			}
		}
		for i := 0; i < bufferLength; i += 1026 {
			chunk := buffer[i+4 : min(bufferLength, i+1026)]
			data = append(data, chunk...)
		}
		cardStatus := binary.LittleEndian.Uint32(data[4:8])
		fingerprint1Status := binary.LittleEndian.Uint32(data[8:12])
		fingerprint2Status := binary.LittleEndian.Uint32(data[12:16])
		cardId := uint16(0)
		cardFacilityCode := uint8(0)
		if cardStatus == 1 {
			cardId = binary.LittleEndian.Uint16(data[24:26])
			cardFacilityCode = data[26]
		}
		fingerprint1 := []byte{}
		if fingerprint1Status == 1 {
			fingerprint1 = data[32 : 32+fingerPrintSize]
		}
		fingerprint2 := []byte{}
		if fingerprint2Status == 1 {
			fingerprint2 = data[32+fingerPrintSize:]
		}
		user := User{
			Id:               userId,
			Level:            Level(0),
			Sensor:           Sensor(0),
			CardFacilityCode: cardFacilityCode,
			CardId:           cardId,
			Fingerprint1:     fingerprint1,
			Fingerprint2:     fingerprint2}
		return user, nil
	}
	return User{}, fmt.Errorf("invalid respond message")
}

func (dev *Sf3000) SetEnrollData(user User) (bool, error) {
	// prepare command bytes array
	const (
		fingerPrintSize int = 1404 + 12
		enrollDataSize  int = (4*8 + fingerPrintSize*2)
	)
	userIds := make([]byte, 4)
	binary.LittleEndian.PutUint32(userIds, uint32(user.Id))
	dev.prepareCommand(0x0104, uint16(user.Id&0xffff), uint16(user.Id>>16))
	// Send command
	dev.conn.Write(dev.command)
	reply1 := make([]byte, 8)
	cnt, err := dev.conn.Read(reply1)
	if err != nil {
		return false, fmt.Errorf("send command to server failed: %v", err)
	}
	if cnt != 8 {
		return false, fmt.Errorf("invalid server reply. Expected message length is 8 but actual is %v", cnt)
	}
	// verify reply status
	if isMessageValid(reply1[:cnt]) {
		// prepare enroll data
		data := make([]byte, 0, enrollDataSize)
		// 1st 4 bytes
		data = append(data, []byte{0, 0, 0, 0}...)
		// 2nd 4 bytes, 1 if card is registered
		if user.CardId > 0 {
			data = append(data, []byte{1, 0, 0, 0}...)
		} else {
			data = append(data, []byte{0, 0, 0, 0}...)
		}
		// 3rd 4 bytes, 1 if fingerprint1 is registered
		if len(user.Fingerprint1) > 0 {
			data = append(data, []byte{1, 0, 0, 0}...)
		} else {
			data = append(data, []byte{0, 0, 0, 0}...)
		}
		// 4th 4 bytes, 1 if fingerprint2 is registered
		if len(user.Fingerprint2) > 0 {
			data = append(data, []byte{1, 0, 0, 0}...)
		} else {
			data = append(data, []byte{0, 0, 0, 0}...)
		}
		// 5th & 6th 4 bytes, reserved
		data = append(data, []byte{0, 0, 0, 0, 0, 0, 0, 0}...)
		// 7th 4 bytes, store card id and facility code
		if user.CardId > 0 {
			data = append(data, []byte{byte(user.CardId & 0xff), byte(user.CardId >> 8), user.CardFacilityCode, 0}...)
		} else {
			data = append(data, []byte{0, 0, 0, 0}...)
		}
		// 8th 4 bytes, 1 if any card or finger registered
		if user.CardId > 0 || len(user.Fingerprint1) > 0 || len(user.Fingerprint2) > 0 {
			data = append(data, []byte{1, 0, 0, 0}...)
		} else {
			data = append(data, []byte{0, 0, 0, 0}...)
		}
		// 9th 1404+12 bytes of fingerprint1 data
		if len(user.Fingerprint1) > 0 {
			data = append(data, user.Fingerprint1...)
		} else {
			data = append(data, make([]byte, fingerPrintSize)...)
		}
		// 10th 1404+12 bytes of fingerprint1 data
		if len(user.Fingerprint2) > 0 {
			data = append(data, user.Fingerprint2...)
		} else {
			data = append(data, make([]byte, fingerPrintSize)...)
		}
		chunk := 0
		buffer := make([]byte, 0)
		for i := 0; i < enrollDataSize; i += 1020 {
			buffer = append(buffer, []byte{0x5a, 0xa5, 0x55, 0x01}...)
			if i+1020 < enrollDataSize {
				buffer = append(buffer, data[i:i+1020]...)
				buffer = append(buffer, []byte{0, 0}...)
				calculateChecksum(buffer[chunk : chunk+1026])
			} else {
				buffer = append(buffer, data[i:]...)
				buffer = append(buffer, []byte{0, 0}...)
				calculateChecksum(buffer[chunk:])
			}
			chunk += 1026
		}
		dev.conn.Write(buffer[:1026])
		dev.conn.Write(buffer[1026:2486])
		dev.conn.Write(buffer[2486:])
		reply2 := make([]byte, 14)
		cnt, err = dev.conn.Read(reply2)
		if err != nil {
			return false, fmt.Errorf("read server reply failed: %v", err)
		}
		if cnt != 14 {
			return false, fmt.Errorf("invalid server reply. Expected message length is 14 but actual is %v", cnt)
		}
		replyStatus := binary.LittleEndian.Uint16(reply2[6:8])
		if replyStatus == 1 && isMessageValid(reply2) {
			return true, nil
		}
	}
	return false, fmt.Errorf("invalid respond message")

}
