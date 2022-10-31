package models

type DeviceResponse struct {
	Command    string     `json:"cmd"`
	ResultCode int        `json:"result_code"`
	ResultData DeviceInfo `json:"result_data"`
}

type DeviceInfo struct {
	Name                string `json:"name"`
	DeviceID            string `json:"deviceId"`
	Firmware            string `json:"firmware"`
	FingerprintVersion  string `json:"fpVer"`
	FaceVersion         string `json:"faceVer"`
	PalmprintVersion    string `json:"pvVer"`
	MaximumBufferLength int    `json:"maxBufferLen"`
	UserLimit           int    `json:"userLimit"`
	FingerprintLimit    int    `json:"fpLimit"`
	FaceLimit           int    `json:"faceLimit"`
	PasswordLimit       int    `json:"pwdLimit"`
	CardLimit           int    `json:"cardLimit"`
	LogLimit            int    `json:"logLimit"`
	UserCount           int    `json:"userCount"`
	ManagerCount        int    `json:"managerCount"`
	FingerprintCount    int    `json:"fpCount"`
	FaceCount           int    `json:"faceCount"`
	PasswordCount       int    `json:"pwdCount"`
	CardCount           int    `json:"cardCount"`
	LogCount            int    `json:"logCount"`
	AllLogCount         int    `json:"allLogCount"`
}
