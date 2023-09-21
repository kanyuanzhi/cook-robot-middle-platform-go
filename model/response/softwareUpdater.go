package response

type CheckUpdatePermission struct {
	IsRunning   bool `json:"isRunning"`
	IsUpdating  bool `json:"isUpdating"`
	IsPermitted bool `json:"isPermitted"`
}

type CheckUpdateInfo struct {
	IsLatest      bool   `json:"isLatest"`
	LatestVersion string `json:"latestVersion"`
	HasFile       bool   `json:"hasFile"`
}
