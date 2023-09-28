package request

import "github.com/kanyuanzhi/cook-robot-middle-platform-go/model"

type ListSeasonings struct {
}

type UpdateSeasoning struct {
	model.SeasoningInfo
}

type UpdateSeasoningsSorts struct {
	Seasonings []model.SeasoningInfo `json:"seasonings"`
}

type DeleteSeasoning struct {
	Id uint `json:"id"`
}

type AddSeasoning struct {
	Name  string `json:"name"`
	Pump  uint32 `json:"pump"`
	Ratio uint32 `json:"ratio"`
}

type UpdateSeasoningPumpRatios struct {
	Seasonings []model.SeasoningInfo `json:"seasonings"`
}
