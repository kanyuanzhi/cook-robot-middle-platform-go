package response

import "github.com/kanyuanzhi/cook-robot-middle-platform-go/model"

type ListSeasonings struct {
	Seasonings []model.SeasoningInfo `json:"seasonings"`
}

type AddSeasoning struct {
	Seasoning model.SeasoningInfo `json:"seasoning"`
}
