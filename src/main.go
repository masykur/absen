package main

import (
	"flag"
	"log"
	"net"
	"os"
	"strconv"
)

const (
	start_dword int = 0xaa55
	end_dword   int = 0x1979
)

var (
	ip       *string
	port     *int
	nid      *int
	password *int
)

func init() {
	ip = flag.String("ip", "127.0.0.1", "Specifies the host name or IP address of the remote device to connect to")
	port = flag.Int("port", 5005, "Specify a port number")
	nid = flag.Int("nid", 1, "Specify NID of device")
	password = flag.Int("pass", 0, "Specify device password")
}

func main() {
	flag.Parse()
	servAddr := *ip + ":" + strconv.Itoa(*port)
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
	log.Printf("Connect to %v\n", tcpAddr)
	connect(conn, *nid, *password)
	productCode := getProductCode(conn, *nid, *password)
	log.Printf("Product code: %v\n", productCode)
}

func calculateChecksum(buffer []byte) uint16 {
	var checksum int = 0
	if len(buffer) > 2 {
		for i := 0; i < len(buffer)-2; i++ {
			checksum += int(buffer[i])
		}
		return uint16(checksum)
	}
	return 0
}

func connect(conn *net.TCPConn, nid int, password int) bool {
	buffer := make([]byte, 16)
	const command int = 0x0052
	buffer[0] = byte(start_dword & 0xff)
	buffer[1] = byte(start_dword >> 8)

	buffer[2] = byte(nid & 0xff)
	buffer[3] = byte(nid >> 8)

	buffer[4] = byte(end_dword & 0xff)
	buffer[5] = byte(end_dword >> 8)

	buffer[6] = byte(command & 0xff)
	buffer[7] = byte(command >> 8)

	buffer[8] = byte(password & 0xff)
	buffer[9] = byte(password >> 8)
	checksum := calculateChecksum(buffer)
	buffer[14] = byte(checksum & 0xff)
	buffer[15] = byte(checksum >> 8)

	reply := make([]byte, 14)
	// log.Println("Send handshake command")
	conn.Write(buffer)
	cnt, err := conn.Read(reply)
	if err != nil {
		log.Fatalf("Send handshake to server failed: %v\n", err)
		os.Exit(1)
	}
	if cnt != 8 {
		log.Fatalf("Invalid server reply: length=%v, data=%v\n", cnt, reply)
		os.Exit(1)
	}
	replyNid := int(reply[2]) + (int(reply[3]) << 8)
	replyStatus := int16(reply[4]) + (int16(reply[5]) << 8)
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
		replyNid = int(reply[2]) + (int(reply[3]) << 8)
		replyStatus = int16(reply[6]) + (int16(reply[7]) << 8)
		if replyStatus == 1 && replyNid == nid {
			log.Println("Connection established")
			return true
		}
	}
	log.Fatalln("Connection failed")
	return false
}
func getProductCode(conn *net.TCPConn, nid int, password int) string {
	const command int = 0x0114
	buffer := make([]byte, 16)

	buffer[0] = byte(start_dword & 0xff)
	buffer[1] = byte(start_dword >> 8)

	buffer[2] = byte(nid & 0xff)
	buffer[3] = byte(nid >> 8)

	buffer[4] = byte(end_dword & 0xff)
	buffer[5] = byte(end_dword >> 8)

	buffer[6] = byte(command & 0xff)
	buffer[7] = byte(command >> 8)

	buffer[8] = byte(password & 0xff)
	buffer[9] = byte(password >> 8)

	checksum := calculateChecksum(buffer)
	buffer[14] = byte(checksum & 0xff)
	buffer[15] = byte(checksum >> 8)

	reply := make([]byte, 54)
	// log.Println("Send handshake command")
	conn.Write(buffer)
	cnt, err := conn.Read(reply)
	if err != nil {
		log.Fatalf("Send handshake to server failed: %v\n", err)
		os.Exit(1)
	}
	if cnt != 8 {
		log.Fatalf("Invalid server reply #1: length=%v, data=%v\n", cnt, reply)
		os.Exit(1)
	}
	replyNid := int(reply[2]) + (int(reply[3]) << 8)
	replyStatus := int16(reply[4]) + (int16(reply[5]) << 8)
	if replyStatus == 1 && replyNid == nid {
		cnt, err = conn.Read(reply)
		if err != nil {
			log.Fatalf("Read server reply failed: %v\n", err)
			os.Exit(1)
		}
		if cnt != 14 {
			log.Fatalf("Invalid server reply #2: length=%v, data=%v\n", cnt, reply)
			os.Exit(1)
		}
		replyNid = int(reply[2]) + (int(reply[3]) << 8)
		replyStatus = int16(reply[6]) + (int16(reply[7]) << 8)
		if replyStatus == 1 && replyNid == nid {
			cnt, err = conn.Read(reply)
			if err != nil {
				log.Fatalf("Read server reply failed: %v\n", err)
				os.Exit(1)
			}
			if cnt != 38 {
				log.Fatalf("Invalid server reply #3: length=%v, data=%v\n", cnt, reply)
				os.Exit(1)
			}
			result := reply[4 : cnt-2]
			return string(result)
		}
	}
	return ""
}
