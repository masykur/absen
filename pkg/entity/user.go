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
	Id               int    `json:"Id"`
	Level            Level  `json:"Level"`
	Sensor           Sensor `json:"Sensor"`
	CardId           uint16 `json:"CardId"`
	CardFacilityCode uint8  `json:"CardFacilityCode"`
	Fingerprint1     []byte `json:"Fingerprint1,omitempty"`
	Fingerprint2     []byte `json:"Fingerprint2,omitempty"`
}

func (user *User) New(id int) {
	user.Id = id
}
