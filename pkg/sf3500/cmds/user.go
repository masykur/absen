package cmds

type GetUserList struct {
	Command string          `json:"cmd"`
	Data    GetUserListData `json:"data"`
}
type GetUserListData struct {
	PackageID int `json:"packageId"`
}

type GetUserInfo struct {
	Command string          `json:"cmd"`
	Data    GetUserInfoData `json:"data"`
}
type GetUserInfoData struct {
	PackageID int      `json:"packageId"`
	UsersId   []string `json:"usersId"`
}
