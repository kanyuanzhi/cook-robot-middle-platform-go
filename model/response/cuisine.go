package response

import "github.com/kanyuanzhi/cook-robot-middle-platform-go/model"

type ListCuisines struct {
	Cuisines []model.CuisineInfo `json:"cuisines"`
}

type AddCuisine struct {
	Cuisine model.CuisineInfo `json:"cuisine"`
}
