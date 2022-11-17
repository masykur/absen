package cmds

type GetLog struct {
	Command string     `json:"cmd"`
	Data    GetLogData `json:"data"`
}
type GetLogData struct {
	PackageID int    `json:"packageId"`
	NewLog    int    `json:"newLog"`
	BeginTime string `json:"beginTime"`
	EndTime   string `json:"endTime"`
	ClearMark int    `json:"clearMark"`
}
