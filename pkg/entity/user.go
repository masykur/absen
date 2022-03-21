package entity

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
	Id           int
	Level        Level
	Sensor       Sensor
	CardId       int
	FingerPrint1 []byte
	FingerPrint2 []byte
}

func (user *User) New(id int) {
	user.Id = id
}
