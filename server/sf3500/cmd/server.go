package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type Device struct {
	ID        string `gorm:"size:20;primaryKey"`
	Name      string `gorm:"size:255"`
	Model     string `gorm:"size:255"`
	Mode      string `gorm:"size:10"` // push or pull
	IsOnline  bool
	IPAddress string `gorm:"size:16"`
	Port      uint16
	Password  string `gorm:"size:255"`
	Logs      []Log
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
	ID           uint   `gorm:"primaryKey;autoIncrement:false"`
	Number       string `gorm:"size:20"`
	Name         string `gorm:"size:200;not null"`
	Privilage    uint16 `gorm:"not null"`
	Photo        string
	Card         string `gorm:"size:10;index"`
	Fingerprints string
	Face         string
	Password     string `gorm:"size:255"`
	ValidStart   string `gorm:"size:8"`
	ValidEnd     string `gorm:"size:8"`
	TimeGroups   string
	Logs         []Log
	ActivatedAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// type UserTimezone struct {
// 	UserID    uint   `gorm:"primaryKey"`
// 	Day       string `gorm:"primaryKey"`
// 	Number    uint16 `gorm:"primaryKey"`
// 	TimeStart time.Time
// 	TimeEnd   time.Time
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// }

type Log struct {
	ID         uint
	DeviceID   string    `gorm:"size:20"`
	UserID     uint      `gorm:"not null"`
	Time       time.Time `gorm:"not null"`
	VerifyMode string    `gorm:"size:20"`
	IoMode     uint16    `gorm:"size:20"`
	InOut      string    `gorm:"size:5;not null"`
	DoorMode   string    `gorm:"size:20"`
	LogPhoto   LogPhoto
	CreatedAt  time.Time
}

type LogPhoto struct {
	gorm.Model
	LogID     uint
	Photo     string
	CreatedAt time.Time
}

type LogJson struct {
	UserID     string `json:"userId"`
	Time       string `json:"time"`
	VerifyMode string `json:"verifyMode"`
	IoMode     uint16 `json:"ioMode"`
	InOut      string `json:"inOut"`
	DoorMode   string `json:"doorMode"`
	LogPhoto   string `json:"logPhoto"`
}

type EnrollJson struct {
	UserID       string   `json:"userId"`
	UserNumber   string   `json:"userNo"`
	Name         string   `json:"name"`
	Privilage    int      `json:"privilege"`
	Photo        string   `json:"photo"`
	Card         string   `json:"card"`
	Fingerprints []string `json:"fps"`
	Face         string   `json:"face"`
	Password     string   `json:"pwd"`
	ValidStart   string   `json:"vaildStart"`
	ValidEnd     string   `json:"vaildEnd"`
	TimeGroups   TimeGroups
}

type TimeGroups struct {
	Sunday   []string `json:"sun"`
	Monday   []string `json:"mon"`
	Tuesday  []string `json:"tue"`
	Wedneday []string `json:"wed"`
	Thursday []string `json:"thu"`
	Friday   []string `json:"fri"`
	Saturday []string `json:"sat"`
}

var serverCommand = &cobra.Command{
	Use:   "server",
	Short: "Run http server",
	Long:  `this command runs server and it won't turn off till you manually do it'`,
	Run:   runServer}

func init() {
	RootCmd.AddCommand(serverCommand)
}

func runServer(cmd *cobra.Command, args []string) {
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Device{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Log{})
	db.AutoMigrate(&LogPhoto{})

	serverAddress := fmt.Sprintf("%v:%v", host, port)
	fmt.Println("Keico SF3500 http server")
	fmt.Println("Created by Ahmad Masykur 2022")
	fmt.Printf("Starting server, listening on \"%v\"\n", serverAddress)

	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(serverAddress, nil)
	// return nil
}

func handleRequest(writer http.ResponseWriter, request *http.Request) {
	if body, err := io.ReadAll(request.Body); err == nil {
		bodyData := []byte(body)
		db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		deviceId := request.Header.Get("dev_id")
		requestCode := request.Header.Get("request_code")

		var device Device
		result := db.First(&device, "id = ?", deviceId)
		addresses := strings.Split(request.RemoteAddr, ":")
		model := request.Header.Get("dev_model")
		device = Device{ID: deviceId, Model: model, IPAddress: addresses[0], Port: port, IsOnline: true, Mode: "push"}
		if result.RowsAffected == 0 {
			if trx := db.Create(&device); trx.Error != nil {
				log.Fatalln("error", trx.Error)
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			if device.Mode != model || device.IPAddress != addresses[0] || device.Port != port || !device.IsOnline || device.Mode != "push" {
				if trx := db.Model(&device).Updates(Device{Model: model, IPAddress: addresses[0], IsOnline: true, Mode: "push"}); trx.Error != nil {
					log.Fatalln("error", trx.Error)
					writer.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
		}

		switch requestCode {
		case "realtime_enroll_data":
			var enroll EnrollJson
			if err = json.Unmarshal(bodyData, &enroll); err == nil {
				if userId, err := strconv.Atoi(enroll.UserID); err == nil {
					var user User
					var fingerprints string
					var timeGroups string
					if fps, err := json.Marshal(enroll.Fingerprints); err == nil && len(fps) > 0 {
						fingerprints = string(fps)
					}
					if tgrs, err := json.Marshal(enroll.Fingerprints); err == nil && len(tgrs) > 0 {
						timeGroups = string(tgrs)
					}
					result := db.First(&user, userId)
					if result.RowsAffected == 0 {
						user = User{ID: uint(userId),
							Number:       enroll.UserNumber,
							Name:         enroll.Name,
							Privilage:    uint16(enroll.Privilage),
							Photo:        enroll.Photo,
							Card:         enroll.Card,
							Fingerprints: fingerprints,
							Face:         enroll.Face,
							Password:     enroll.Password,
							ValidStart:   enroll.ValidStart,
							ValidEnd:     enroll.ValidEnd,
							TimeGroups:   timeGroups}
						if trx := db.Create(&user); trx.Error != nil {
							log.Fatalln("error", trx.Error)
							writer.WriteHeader(http.StatusInternalServerError)
							return
						}
					} else {
						if trx := db.Model(&user).Updates(User{
							Number:       enroll.UserNumber,
							Name:         enroll.Name,
							Privilage:    uint16(enroll.Privilage),
							Photo:        enroll.Photo,
							Card:         enroll.Card,
							Fingerprints: fingerprints,
							Face:         enroll.Face,
							Password:     enroll.Password,
							ValidStart:   enroll.ValidStart,
							ValidEnd:     enroll.ValidEnd,
							TimeGroups:   timeGroups}); trx.Error != nil {
							log.Fatalln("error", trx.Error)
							writer.WriteHeader(http.StatusInternalServerError)
							return
						}
					}
				} else {
					log.Fatalln("error", err)
					writer.WriteHeader(http.StatusInternalServerError)
					return
				}
			} else {
				log.Fatalln("error", err)
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
		case "realtime_glog":
			var glog LogJson
			if err = json.Unmarshal(bodyData, &glog); err == nil {
				var logRecord Log
				var logPhoto LogPhoto
				if datetime, err := time.ParseInLocation("20060102150405", glog.Time, time.Local); err == nil {
					logRecord.Time = datetime
					if userId, err := strconv.Atoi(glog.UserID); err == nil {
						var user User
						result := db.First(&user, userId)
						if result.RowsAffected == 0 {
							user = User{ID: uint(userId), Name: strconv.Itoa(userId), Privilage: 0}
							if trx := db.Create(&user); trx.Error != nil {
								log.Fatalln("error", trx.Error)
								writer.WriteHeader(http.StatusInternalServerError)
								return
							}
						}

						logRecord = Log{Time: datetime, UserID: uint(userId), DeviceID: deviceId, DoorMode: glog.DoorMode, InOut: glog.InOut, IoMode: glog.IoMode, VerifyMode: glog.VerifyMode}
						if trx := db.Create(&logRecord); trx.Error != nil {
							log.Fatalln("error", trx.Error)
							writer.WriteHeader(http.StatusInternalServerError)
							return
						}
						logPhoto = LogPhoto{LogID: logRecord.ID, Photo: glog.LogPhoto}
						if trx := db.Create(&logPhoto); trx.Error != nil {
							log.Fatalln("error", trx.Error)
							writer.WriteHeader(http.StatusInternalServerError)
							return
						}
					}
				}
			} else {
				log.Fatalln("error", err)
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
	writer.Header()["response_code"] = []string{"OK"}
	writer.Header()["trans_id"] = []string{"100"}
	writer.WriteHeader(http.StatusOK)
}
