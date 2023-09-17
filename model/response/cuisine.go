package response

import "github.com/kanyuanzhi/middle-platform/model"

type ListCuisines struct {
	Cuisines []model.CuisineInfo `json:"cuisines"`
}

type AddCuisine struct {
	Cuisine model.CuisineInfo `json:"cuisine"`
}
