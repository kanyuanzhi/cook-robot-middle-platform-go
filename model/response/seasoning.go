package response

import "github.com/kanyuanzhi/middle-platform/model"

type ListSeasonings struct {
	Seasonings []model.SeasoningInfo `json:"seasonings"`
}

type AddSeasoning struct {
	Seasoning model.SeasoningInfo `json:"seasoning"`
}
