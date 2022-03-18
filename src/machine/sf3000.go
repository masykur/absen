package machine

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
	"unsafe"
)

type Sf3000 struct {
	conn    *net.TCPConn
	command []byte
}

// Convert byte array to string without allocate new memory
func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func (dev *Sf3000) prepareCommand(command uint16, parameters ...uint16) {
	_ = dev.command[15] // early bounds check to guarantee safety of writes below
	binary.LittleEndian.PutUint16(dev.command[6:8], command)
	for i := 0; i < 4; i++ {
		if len(parameters) > i {
			binary.LittleEndian.PutUint16(dev.command[8+i*2:10+i*2], parameters[i])
		} else {
			binary.LittleEndian.PutUint16(dev.command[8+i*2:10+i*2], 0)
		}
	}
	var checksum uint16 = 0
	for _, v := range dev.command[0:14] {
		checksum += uint16(v)
	}
	binary.LittleEndian.PutUint16(dev.command[14:16], checksum)
}
func isMessageValid(bytes []byte) bool {
	_ = bytes[1] // early bounds check to guarantee safety of writes below
	var checksum uint16 = 0
	for _, v := range bytes[0 : len(bytes)-2] {
		checksum += uint16(v)
	}
	return checksum == binary.LittleEndian.Uint16(bytes[len(bytes)-2:])
}

// Make connection and authenticate to machine
func (dev *Sf3000) Connect(conn *net.TCPConn, nid uint16, password uint16) (bool, error) {
	dev.command = []byte{0x55, 0xaa, 0x0, 0x0, 0x79, 0x19, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	binary.LittleEndian.PutUint16(dev.command[2:4], nid)
	// connect command = 0x0052
	dev.prepareCommand(0x0052, password)
	//const command uint16 = 0x0052
	dev.conn = conn

	reply := make([]byte, 14)
	// Send authentication command
	dev.conn.Write(dev.command)
	cnt, err := dev.conn.Read(reply)
	if err != nil {
		return false, err
	}
	if cnt != 8 {
		return false, fmt.Errorf("invalid server reply. Expected message length is 8 but actual is %v", cnt)
	}
	replyStatus := binary.LittleEndian.Uint16(reply[4:6])
	if replyStatus == 1 && isMessageValid(reply[:cnt]) {
		cnt, err = conn.Read(reply)
		if err != nil {
			return false, err
		}
		if cnt != 14 {
			return false, fmt.Errorf("invalid server reply. Expected message length is 14 but actual is %v", cnt)
		}
		replyStatus = binary.LittleEndian.Uint16(reply[6:8])
		if replyStatus == 1 && isMessageValid(reply[:cnt]) {
			return true, nil
		}
	}
	return false, fmt.Errorf("connection failed")
}

// Obtain product code from machine
func (dev *Sf3000) GetProductCode() (string, error) {
	// prepare command bytes array
	dev.prepareCommand(0x0114, 0)
	// Send command
	dev.conn.Write(dev.command)
	reply1 := make([]byte, 8)
	cnt, err := dev.conn.Read(reply1)
	if err != nil {
		return "", fmt.Errorf("send command to server failed: %v", err)
	}
	if cnt != 8 {
		return "", fmt.Errorf("invalid server reply. Expected message length is 8 but actual is %v", cnt)
	}
	replyStatus := binary.LittleEndian.Uint16(reply1[4:6])
	// verify reply status
	if replyStatus == 1 && isMessageValid(reply1) {
		reply2 := make([]byte, 14)
		// read second message
		cnt, err = dev.conn.Read(reply2)
		if err != nil {
			return "", fmt.Errorf("read server reply failed: %v", err)
		}
		if cnt != 14 {
			return "", fmt.Errorf("invalid server reply. Expected message length is 14 but actual is %v", cnt)
		}
		replyStatus = binary.LittleEndian.Uint16(reply2[6:8])
		// verify second reply status
		if replyStatus == 1 && isMessageValid(reply2) {
			reply3 := make([]byte, 38)
			// read actual response contain product code
			cnt, err = dev.conn.Read(reply3)
			if err != nil {
				return "", fmt.Errorf("read server reply failed: %v", err)
			}
			if cnt != 38 {
				return "", fmt.Errorf("invalid server reply. Expected message length is 38 but actual is %v", cnt)
			}
			// parse and return product code info
			result := reply3[4 : cnt-2]
			return b2s(result), nil
		}
	}
	return "", fmt.Errorf("invalid respond message")
}

// Obtain machine serial number
func (dev *Sf3000) GetSerialNumber() (string, error) {
	// prepare command bytes array
	dev.prepareCommand(0x0115, 0)
	// Send command
	dev.conn.Write(dev.command)
	reply1 := make([]byte, 8)
	cnt, err := dev.conn.Read(reply1)
	if err != nil {
		return "", fmt.Errorf("send command to server failed: %v", err)
	}
	if cnt != 8 {
		return "", fmt.Errorf("invalid server reply. Expected message length is 8 but actual is %v", cnt)
	}
	replyStatus := binary.LittleEndian.Uint16(reply1[4:6])
	// verify reply status
	if replyStatus == 1 && isMessageValid(reply1) {
		reply2 := make([]byte, 14)
		// read second message
		cnt, err = dev.conn.Read(reply2)
		if err != nil {
			return "", fmt.Errorf("read server reply failed: %v", err)
		}
		if cnt != 14 {
			return "", fmt.Errorf("invalid server reply. Expected message length is 14 but actual is %v", cnt)
		}
		replyStatus = binary.LittleEndian.Uint16(reply2[6:8])
		// verify second reply status
		if replyStatus == 1 && isMessageValid(reply2) {
			reply3 := make([]byte, 38)
			// read actual response contain product code
			cnt, err = dev.conn.Read(reply3)
			if err != nil {
				return "", fmt.Errorf("read server reply failed: %v", err)
			}
			if cnt != 38 {
				return "", fmt.Errorf("invalid server reply. Expected message length is 38 but actual is %v", cnt)
			}
			// parse and return product code info
			result := reply3[4 : cnt-2]
			return b2s(result), nil
		}
	}
	return "", fmt.Errorf("invalid respond message")
}

// Obtain current date and time from machine
func (dev *Sf3000) GetDateTime() (time.Time, error) {
	// prepare command bytes array
	dev.prepareCommand(0x011D, 0x0004)
	// Send command
	dev.conn.Write(dev.command)
	reply1 := make([]byte, 8)
	cnt, err := dev.conn.Read(reply1)
	if err != nil {
		return time.Time{}, fmt.Errorf("send command to server failed: %v", err)
	}
	if cnt != 8 {
		return time.Time{}, fmt.Errorf("invalid server reply. Expected message length is 8 but actual is %v", cnt)
	}
	replyStatus := binary.LittleEndian.Uint16(reply1[4:6])
	// verify reply status
	if replyStatus == 1 && isMessageValid(reply1) {
		reply2 := make([]byte, 10)
		// read second message
		cnt, err = dev.conn.Read(reply2)
		if err != nil {
			return time.Time{}, fmt.Errorf("read server reply failed: %v", err)
		}
		if cnt != 10 {
			return time.Time{}, fmt.Errorf("invalid server reply. Expected message length is 10 but actual is %v", cnt)
		}
		if isMessageValid(reply2) {
			// get date and time data
			num := binary.LittleEndian.Uint32(reply2[4:8])
			// first date is January 1st, 2000
			firstDate := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.Local)
			// device date time is number of seconds after first date
			date := firstDate.Add(time.Second * time.Duration(num))

			// verify respond status
			reply3 := make([]byte, 14)
			// read actual response contain product code
			cnt, err = dev.conn.Read(reply3)
			if err != nil {
				return time.Time{}, fmt.Errorf("read server reply failed: %v", err)
			}
			if cnt != 14 {
				return time.Time{}, fmt.Errorf("invalid server reply. Expected message length is 14 but actual is %v", cnt)
			}
			replyStatus = binary.LittleEndian.Uint16(reply3[6:8])
			if replyStatus == 1 && isMessageValid(reply3) {
				return date, nil
			}
		}
	}
	return time.Time{}, fmt.Errorf("invalid respond message")
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
					CardId: int(uCardId)})
			}
			return users, nil
		}
	}
	return []User{}, fmt.Errorf("invalid respond message")
}

func (dev *Sf3000) GetEnrollData(userId int) (User, error) {
	// prepare command bytes array
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
		reply2 := make([]byte, 1026) // max package size is 1026 bytes
		data := make([]byte, 0, 2864)
		for {
			cnt, err = dev.conn.Read(reply2)
			if err != nil {
				return User{}, fmt.Errorf("read server reply failed: %v", err)
			}
			data = append(data, reply2[4:cnt-2]...)
			if cnt < 1026 || len(data) == cap(data) {
				break
			}
		}
		user := User{
			Id:     userId,
			Level:  Level(0),
			Sensor: Sensor(0),
			CardId: int(binary.LittleEndian.Uint16(data[24:26])),
			Data:   data[32:]}
		return user, nil
	}
	return User{}, fmt.Errorf("invalid respond message")

}
