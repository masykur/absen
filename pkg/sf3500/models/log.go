package models

import (
	"fmt"
	"strings"
	"time"
)

type LogResponse struct {
	Command    string  `json:"cmd"`
	ResultCode int     `json:"result_code"`
	ResultData LogInfo `json:"result_data"`
}

type LogInfo struct {
	PackageID   int       `json:"packageId"`
	LogCount    int       `json:"logCount"`
	AllLogCount int       `json:"allLogCount"`
	Logs        []LogData `json:"logs"`
}

type LogData struct {
	UserID     string     `json:"userId"`
	Time       CustomTime `json:"time"`
	VerifyMode string     `json:"verifyMode"`
	IoMode     int        `json:"ioMode"`
	InOut      string     `json:"inOut"`
	DoorMode   string     `json:"doorMode"`
	LogPhoto   string     `json:"logPhoto"`
}

// CustomTime provides an example of how to declare a new time Type with a custom formatter.
// Note that time.Time methods are not available, if needed you can add and cast like the String method does
// Otherwise, only use in the json struct at marshal/unmarshal time.
type CustomTime time.Time

const ctLayout = "20060102150405"

// UnmarshalJSON Parses the json string in the custom format
func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	nt, err := time.Parse(ctLayout, s)
	*ct = CustomTime(nt)
	return
}

// MarshalJSON writes a quoted string in the custom format
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(ct.String()), nil
}

// String returns the time in the custom format
func (ct *CustomTime) String() string {
	t := time.Time(*ct)
	return fmt.Sprintf("%q", t.Format(ctLayout))
}

// Returns the time object
func (ct *CustomTime) Time() time.Time {
	t := time.Time(*ct)
	return t
}
