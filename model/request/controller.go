package request

type CommandData struct {
	UUID            string  `json:"uuid"`
	CustomStepsUUID string  `json:"customStepsUuid"`
	Temperature     float64 `json:"temperature"`
}

type ExecuteCommandRequest struct {
	CommandType string      `json:"commandType" form:"commandType"`
	CommandName string      `json:"commandName" form:"commandName"`
	CommandData CommandData `json:"commandData" form:"commandData"`
}
