package sf3500

import (
	"encoding/json"
	"fmt"

	"github.com/masykur/absen/pkg/sf3500/cmds"
	"github.com/masykur/absen/pkg/sf3500/models"
)

func (dev *Sf3500) GetUserList(packageId int) (int, []models.User, error) {
	command := cmds.GetUserList{Command: "GetUserIdList", Data: cmds.GetUserListData{PackageID: packageId}}
	commandBytes, _ := json.Marshal(command)
	if response, err := dev.sendCommand(commandBytes); err == nil {
		var userList models.UserResponse
		if err := json.Unmarshal(response, &userList); err == nil {
			if userList.ResultCode == 0 {
				return userList.ResultData.PackageID, userList.ResultData.Users, nil
			} else {
				return 0, nil, fmt.Errorf("error code %d", userList.ResultCode)
			}
		} else {
			return 0, nil, err
		}
	} else {
		return 0, nil, err
	}
}

func (dev *Sf3500) GetUserInfo(packageId int, userIds ...string) (int, []models.User, error) {
	command := cmds.GetUserInfo{Command: "GetUserInfo", Data: cmds.GetUserInfoData{PackageID: packageId, UsersId: userIds}}
	if commandBytes, err := json.Marshal(command); err == nil {
		if response, err := dev.sendCommand(commandBytes); err == nil {
			var userList models.UserResponse
			if err := json.Unmarshal(response, &userList); err == nil {
				if userList.ResultCode == 0 {
					return userList.ResultData.PackageID, userList.ResultData.Users, nil
				} else {
					return 0, nil, fmt.Errorf("error code %d", userList.ResultCode)
				}
			} else {
				fmt.Println("JSON error", err)
				return 0, nil, err
			}
		} else {
			return 0, nil, err
		}
	} else {
		return 0, nil, err
	}
}
