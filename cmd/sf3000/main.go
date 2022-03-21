package main

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/masykur/keico/pkg/entity"
	"github.com/masykur/keico/pkg/machines"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var (
	host                *string
	port                *int
	nid                 *int
	password            *int
	getProductFlag      *bool
	getSerialNumberFlag *bool
	getTimeFlag         *bool
	getUserCountFlag    *bool
	getUsersFlag        *bool
	getUser             *int
)

func init() {
	host = flag.String("host", "127.0.0.1", "Specifies the host name or IP address of the remote machine to connect to")
	port = flag.Int("port", 5005, "Specify a port number")
	nid = flag.Int("nid", 1, "Specify machine number")
	password = flag.Int("pass", 0, "Specify device password")
	getProductFlag = flag.Bool("get-product", false, "Obtain machine product code")
	getSerialNumberFlag = flag.Bool("get-serial-number", false, "Obtain machine serial number")
	getTimeFlag = flag.Bool("get-time", false, "Get machine date and time")
	getUserCountFlag = flag.Bool("get-user-count", false, "Get number of user registered into machine")
	getUsersFlag = flag.Bool("get-users", false, "Get number of user registered in the machine")
	getUser = flag.Int("get-user", 0, "Get enroll data of user")
}

func main() {
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}
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
	device := new(machines.Sf3000)
	connected, err := device.Connect(conn, uint16(*nid), uint16(*password))
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	//connected := connect(conn, uint16(*nid), uint16(*password))
	if connected {
		if *getProductFlag {
			var productCode string
			productCode, err = device.GetProductCode()
			if err == nil {
				fmt.Println(productCode)
			} else {
				log.Fatalln(err)
			}
		}
		if *getSerialNumberFlag {
			var serialNumber string
			serialNumber, err = device.GetSerialNumber()
			if err == nil {
				fmt.Println(serialNumber)
			} else {
				log.Fatalln(err)
			}
		}
		if *getTimeFlag {
			var dateTime time.Time
			dateTime, err = device.GetDateTime()
			if err == nil {
				fmt.Println(dateTime)
			} else {
				log.Fatalln(err)
			}
		}
		if *getUserCountFlag {
			var userCount int
			userCount, err = device.GetUserCount()
			if err == nil {
				fmt.Println(userCount)
			} else {
				log.Fatalln(err)
			}
		}
		if *getUsersFlag {
			var users []entity.User
			users, err = device.GetUsers()
			if err == nil {
				fmt.Println("  No    User ID    Priv  Sensor Card ID")
				fmt.Println("---- ---------- ------- ------- -------")
				for i, user := range users {
					fmt.Printf("%4s %10s\t%7s\t%7s\t%7s\n", strconv.Itoa(i+1), strconv.Itoa(user.Id), strconv.Itoa(int(user.Level)), strconv.Itoa(int(user.Sensor)), strconv.Itoa(user.CardId))
				}
			} else {
				log.Fatalln(err)
			}
		}
		if *getUser != 0 {
			var user entity.User
			user, err = device.GetEnrollData(*getUser)
			if err == nil {
				fmt.Println("   User ID      Card ID Finger Print Data")
				fmt.Println("---------- ------------ -----------------")
				fmt.Printf("%10s\t%7s\t%v %v\n", strconv.Itoa(user.Id), strconv.Itoa(int(user.CardId)), hex.EncodeToString(user.FingerPrint1), hex.EncodeToString(user.FingerPrint2))
				f, _ := os.Create("data1.bmp")
				w := bufio.NewWriter(f)
				w.Write(user.FingerPrint1[12 : len(user.FingerPrint1)-4])
				w.Flush()
				f.Close()
				s := 0
				for _, v := range user.FingerPrint2 {
					s += int(v)
				}
				b := make([]byte, 4)
				binary.LittleEndian.PutUint32(b, uint32(s))
				fmt.Println(s, hex.EncodeToString(b))
			} else {
				log.Fatalln(err)
			}
		}
	}
}
