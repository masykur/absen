package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

type level int16

// const (
// 	userLevel   level = 0
// 	masterLevel level = 1
// )

type sensor int16

// const (
// 	fingerPrint1 sensor = 1
// 	fingerPrint2 sensor = 2
// 	card         sensor = 8
// )

type user struct {
	id     int
	level  level
	sensor sensor
	cardId int
}

const (
	start_dword int = 0xaa55
	end_dword   int = 0x1979
)

var (
	host             *string
	port             *int
	nid              *int
	password         *int
	getProductFlag   *bool
	getTimeFlag      *bool
	getUserCountFlag *bool
	getUsersFlag     *bool
)

func init() {
	host = flag.String("host", "127.0.0.1", "Specifies the host name or IP address of the remote machine to connect to")
	port = flag.Int("port", 5005, "Specify a port number")
	nid = flag.Int("nid", 1, "Specify NID of device")
	password = flag.Int("pass", 0, "Specify device password")
	getProductFlag = flag.Bool("get-product", false, "Get product code")
	getTimeFlag = flag.Bool("get-time", false, "Get machine date and time")
	getUserCountFlag = flag.Bool("get-user-count", false, "Get number of user registered into machine")
	getUsersFlag = flag.Bool("get-users", false, "Get number of user registered into machine")
}

func main() {
	flag.Parse()
	servAddr := *host + ":" + strconv.Itoa(*port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	log.Printf("Connect to %v\n", tcpAddr)
	connected := connect(conn, uint16(*nid), uint16(*password))
	if connected {
		if *getProductFlag {
			productCode := getProductCode(conn, uint16(*nid))
			log.Println(productCode)
		}
		if *getTimeFlag {
			dateTime := getDateTime(conn, uint16(*nid))
			log.Println(dateTime)
		}
		if *getUserCountFlag {
			userCount := getUserCount(conn, uint16(*nid))
			log.Println(userCount)
		}
		if *getUsersFlag {
			users := getUsers(conn, uint16(*nid))
			fmt.Println("  No    User ID    Priv  Sensor Card ID")
			fmt.Println("---- ---------- ------- ------- -------")
			for i, user := range users {
				fmt.Printf("%4s %10s\t%7s\t%7s\t%7s\n", strconv.Itoa(i+1), strconv.Itoa(user.id), strconv.Itoa(int(user.level)), strconv.Itoa(int(user.sensor)), strconv.Itoa(user.cardId))
			}
		}
	}
}

func calculateChecksum(buffer []byte) {
	_ = buffer[1] // early bounds check to guarantee safety of writes below
	var checksum uint16 = 0
	buffLen := len(buffer)
	for _, v := range buffer[0 : buffLen-2] {
		checksum += uint16(v)
	}
	binary.LittleEndian.PutUint16(buffer[buffLen-2:buffLen], checksum)
}

func connect(conn *net.TCPConn, nid uint16, password uint16) bool {
	const command uint16 = 0x0052
	buffer := make([]byte, 16)
	binary.LittleEndian.PutUint16(buffer[0:2], uint16(start_dword))
	binary.LittleEndian.PutUint16(buffer[2:4], uint16(nid))
	binary.LittleEndian.PutUint16(buffer[4:6], uint16(end_dword))
	binary.LittleEndian.PutUint16(buffer[6:8], command)
	binary.LittleEndian.PutUint16(buffer[8:10], password)
	calculateChecksum(buffer)

	reply := make([]byte, 14)
	// Send authentication command
	conn.Write(buffer)
	cnt, err := conn.Read(reply)
	if err != nil {
		log.Fatalf("Send authentication to server failed: %v\n", err)
		os.Exit(1)
	}
	if cnt != 8 {
		log.Fatalf("Invalid server reply: length=%v, data=%v\n", cnt, reply)
		os.Exit(1)
	}
	replyNid := binary.LittleEndian.Uint16(reply[2:4])
	replyStatus := binary.LittleEndian.Uint16(reply[4:6])
	if replyStatus == 1 && replyNid == nid {
		cnt, err = conn.Read(reply)
		if err != nil {
			log.Fatalf("Read server reply failed: %v\n", err)
			os.Exit(1)
		}
		if cnt != 14 {
			log.Fatalf("Invalid server reply: length=%v, data=%v\n", cnt, reply)
			os.Exit(1)
		}
		replyNid = binary.LittleEndian.Uint16(reply[2:4])
		replyStatus = binary.LittleEndian.Uint16(reply[6:8])
		if replyStatus == 1 && replyNid == nid {
			log.Println("Connection established")
			return true
		}
	}
	log.Fatalln("Connection failed")
	return false
}

func getProductCode(conn *net.TCPConn, nid uint16) string {
	const command uint32 = 0x00000114
	buffer := make([]byte, 16)
	binary.LittleEndian.PutUint16(buffer[0:2], uint16(start_dword))
	binary.LittleEndian.PutUint16(buffer[2:4], uint16(nid))
	binary.LittleEndian.PutUint16(buffer[4:6], uint16(end_dword))
	binary.LittleEndian.PutUint32(buffer[6:10], command)
	calculateChecksum(buffer)
	// Send command
	conn.Write(buffer)
	reply1 := make([]byte, 8)
	cnt, err := conn.Read(reply1)
	if err != nil {
		log.Fatalf("Send command to server failed: %v\n", err)
		os.Exit(1)
	}
	if cnt != 8 {
		log.Fatalf("Invalid server reply #1: length=%v, data=%v\n", cnt, reply1)
		os.Exit(1)
	}
	replyNid := binary.LittleEndian.Uint16(reply1[2:4])
	replyStatus := binary.LittleEndian.Uint16(reply1[4:6])
	// verify reply status
	if replyStatus == 1 && replyNid == nid {
		reply2 := make([]byte, 14)
		// read second message
		cnt, err = conn.Read(reply2)
		if err != nil {
			log.Fatalf("Read server reply failed: %v\n", err)
			os.Exit(1)
		}
		if cnt != 14 {
			log.Fatalf("Invalid server reply #2: length=%v, data=%v\n", cnt, reply2)
			os.Exit(1)
		}
		replyNid = binary.LittleEndian.Uint16(reply2[2:4])
		replyStatus = binary.LittleEndian.Uint16(reply2[6:8])
		// verify second reply status
		if replyStatus == 1 && replyNid == nid {
			reply3 := make([]byte, 38)
			// read actual response contain product code
			cnt, err = conn.Read(reply3)
			if err != nil {
				log.Fatalf("Read server reply failed: %v\n", err)
				os.Exit(1)
			}
			if cnt != 38 {
				log.Fatalf("Invalid server reply #3: length=%v, data=%v\n", cnt, reply3)
				os.Exit(1)
			}
			// parse and return product code info
			result := reply3[4 : cnt-2]
			return string(result)
		}
	}
	return ""
}

func getDateTime(conn *net.TCPConn, nid uint16) time.Time {
	const command uint32 = 0x0004011D
	buffer := make([]byte, 16)
	binary.LittleEndian.PutUint16(buffer[0:2], uint16(start_dword))
	binary.LittleEndian.PutUint16(buffer[2:4], uint16(nid))
	binary.LittleEndian.PutUint16(buffer[4:6], uint16(end_dword))
	binary.LittleEndian.PutUint32(buffer[6:10], command)
	calculateChecksum(buffer)

	// Send command
	conn.Write(buffer)
	reply1 := make([]byte, 8)
	cnt, err := conn.Read(reply1)
	if err != nil {
		log.Fatalf("Send command to server failed: %v\n", err)
		os.Exit(1)
	}
	if cnt != 8 {
		log.Fatalf("Invalid server reply #1: length=%v, data=%v\n", cnt, reply1)
		os.Exit(1)
	}
	replyNid := binary.LittleEndian.Uint16(reply1[2:4])
	replyStatus := binary.LittleEndian.Uint16(reply1[4:6])
	// verify reply status
	if replyStatus == 1 && replyNid == nid {
		reply2 := make([]byte, 10)
		// read second message
		cnt, err = conn.Read(reply2)
		if err != nil {
			log.Fatalf("Read server reply failed: %v\n", err)
			os.Exit(1)
		}
		if cnt != 10 {
			log.Fatalf("Invalid server reply #2: length=%v, data=%v\n", cnt, reply2)
			os.Exit(1)
		}
		replyNid = binary.LittleEndian.Uint16(reply2[2:4])
		// verify second reply status
		if replyNid == nid {
			// get date and time data
			num := binary.LittleEndian.Uint32(reply2[4:8])
			// first date is January 1st, 2000
			firstDate := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.Local)
			// device date time is number of seconds after first date
			date := firstDate.Add(time.Second * time.Duration(num))

			// verify respond status
			reply3 := make([]byte, 14)
			// read actual response contain product code
			cnt, err = conn.Read(reply3)
			if err != nil {
				log.Fatalf("Read server reply failed: %v\n", err)
				os.Exit(1)
			}
			if cnt != 14 {
				log.Fatalf("Invalid server reply #3: length=%v, data=%v\n", cnt, reply3)
				os.Exit(1)
			}
			replyNid = binary.LittleEndian.Uint16(reply3[2:4])
			replyStatus = binary.LittleEndian.Uint16(reply3[6:8])
			if replyStatus == 1 && replyNid == nid {
				return date
			}
		}
	}
	return time.Time{}
}

func getUserCount(conn *net.TCPConn, nid uint16) int {
	const command uint32 = 0x00000116
	buffer := make([]byte, 16)
	binary.LittleEndian.PutUint16(buffer[0:2], uint16(start_dword))
	binary.LittleEndian.PutUint16(buffer[2:4], uint16(nid))
	binary.LittleEndian.PutUint16(buffer[4:6], uint16(end_dword))
	binary.LittleEndian.PutUint32(buffer[6:10], command)
	binary.LittleEndian.PutUint32(buffer[10:14], 0x00010000)
	calculateChecksum(buffer)

	// Send command
	conn.Write(buffer)
	reply1 := make([]byte, 8)
	cnt, err := conn.Read(reply1)
	if err != nil {
		log.Fatalf("Send command to server failed: %v\n", err)
		os.Exit(1)
	}
	if cnt != 8 {
		log.Fatalf("Invalid server reply #1: length=%v, data=%v\n", cnt, reply1)
		os.Exit(1)
	}
	replyNid := binary.LittleEndian.Uint16(reply1[2:4])
	replyStatus := binary.LittleEndian.Uint16(reply1[4:6])
	// verify reply status
	if replyStatus == 1 && replyNid == nid {
		reply2 := make([]byte, 14)
		// read second message
		cnt, err = conn.Read(reply2)
		if err != nil {
			log.Fatalf("Read server reply failed: %v\n", err)
			os.Exit(1)
		}
		if cnt != 14 {
			log.Fatalf("Invalid server reply #2: length=%v, data=%v\n", cnt, reply2)
			os.Exit(1)
		}
		replyNid = binary.LittleEndian.Uint16(reply2[2:4])
		replyStatus = binary.LittleEndian.Uint16(reply2[6:8])
		// verify second reply status
		if replyStatus == 1 && replyNid == nid {
			// get user count data
			num := binary.LittleEndian.Uint32(reply2[8:12])
			// first date is January 1st, 2000
			return int(num)
		}
	}
	return 0
}

func getUsers(conn *net.TCPConn, nid uint16) []user {
	const command uint32 = 0x00000109
	buffer := make([]byte, 16)
	binary.LittleEndian.PutUint16(buffer[0:2], uint16(start_dword))
	binary.LittleEndian.PutUint16(buffer[2:4], uint16(nid))
	binary.LittleEndian.PutUint16(buffer[4:6], uint16(end_dword))
	binary.LittleEndian.PutUint32(buffer[6:10], command)
	calculateChecksum(buffer)

	// Send command
	conn.Write(buffer)
	reply1 := make([]byte, 8)
	cnt, err := conn.Read(reply1)
	if err != nil {
		log.Fatalf("Send command to server failed: %v\n", err)
		os.Exit(1)
	}
	if cnt != 8 {
		log.Fatalf("Invalid server reply #1: length=%v, data=%v\n", cnt, reply1)
		os.Exit(1)
	}
	replyNid := binary.LittleEndian.Uint16(reply1[2:4])
	replyStatus := binary.LittleEndian.Uint16(reply1[4:6])
	// verify reply status
	if replyStatus == 1 && replyNid == nid {
		// read second message
		reply2 := make([]byte, 14)
		cnt, err = conn.Read(reply2)
		if err != nil {
			log.Fatalf("Read server reply failed: %v\n", err)
			os.Exit(1)
		}
		if cnt != 14 {
			log.Fatalf("Invalid server reply #2: length=%v, data=%v\n", cnt, reply2)
			os.Exit(1)
		}
		replyNid = binary.LittleEndian.Uint16(reply2[2:4])
		replyStatus = binary.LittleEndian.Uint16(reply2[6:8])
		// verify second reply status
		if replyStatus == 1 && replyNid == nid {
			// get user count data
			dataLength := binary.LittleEndian.Uint32(reply2[8:12])
			//log.Printf("User count = %v | %v\n", userCount, reply2)
			reply3 := make([]byte, 1026) // max package size is 2026 bytes
			data := make([]byte, 0, dataLength*8)
			for {
				cnt, err = conn.Read(reply3)
				if err != nil {
					log.Fatalf("Read server reply failed: %v\n", err)
					os.Exit(1)
				}
				data = append(data, reply3[4:cnt-2]...)
				if cnt == 0 || len(data) == cap(data) {
					break
				}
			}

			users := make([]user, 0, dataLength)
			for i := uint32(0); i < dataLength; i++ {
				uId := binary.LittleEndian.Uint32(data[i*8 : 4+i*8])
				uLevel := data[4+i*8]
				uSensor := data[5+i*8]
				uCardId := binary.LittleEndian.Uint16(data[6+i*8 : 8+i*8])
				users = append(users, user{int(uId), level(uLevel), sensor(uSensor), int(uCardId)})
			}
			return users
		}
	}
	return []user{}
}
