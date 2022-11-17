package models

type UserResponse struct {
	Command    string   `json:"cmd"`
	ResultCode int      `json:"result_code"`
	ResultData UserInfo `json:"result_data"`
}

type UserInfo struct {
	PackageID int    `json:"packageId"`
	UserCount int    `json:"userCount"`
	Users     []User `json:"users"`
}

type User struct {
	UserID       string     `json:"userId"`
	Name         string     `json:"name"`
	Privilage    int        `json:"privilage"`
	Photo        string     `json:"photo"`
	Card         string     `json:"card"`
	Fingerprints []string   `json:"fps"`
	Face         string     `json:"face"`
	Palmprint    string     `json:"palm"`
	Password     string     `json:"pwd"`
	ValidStart   string     `json:"vaildStart"`
	ValidEnd     string     `json:"vaildEnd"`
	TimeGroups   TimeGroups `json:"timeGroups,omitempty"`
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
