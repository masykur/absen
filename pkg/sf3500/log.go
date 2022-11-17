package sf3500

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/masykur/absen/pkg/sf3500/cmds"

	"github.com/masykur/absen/pkg/sf3500/models"
)

func (dev *Sf3500) GetLog(packageId int, newLog int, beginTime time.Time, endTime time.Time, clearMark int) (int, []models.LogData, error) {
	command := cmds.GetLog{Command: "GetLogData", Data: cmds.GetLogData{PackageID: packageId, NewLog: newLog, BeginTime: beginTime.Format("20060102"), EndTime: endTime.Format("20060102"), ClearMark: clearMark}}
	if commandBytes, err := json.Marshal(command); err == nil {
		if response, err := dev.sendCommand(commandBytes); err == nil {
			var logList models.LogResponse
			if err := json.Unmarshal(response, &logList); err == nil {
				if logList.ResultCode == 0 {
					return logList.ResultData.PackageID, logList.ResultData.Logs, nil
				} else {
					return 0, nil, fmt.Errorf("error code %d", logList.ResultCode)
				}
			} else {
				return 0, nil, err
			}
		} else {
			return 0, nil, err
		}
	} else {
		return 0, nil, err
	}
}
