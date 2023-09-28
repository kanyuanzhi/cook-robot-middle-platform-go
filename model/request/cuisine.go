package request

import "github.com/kanyuanzhi/cook-robot-middle-platform-go/model"

type ListCuisines struct {
}

type UpdateCuisineName struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type UpdateCuisineUnDeletable struct {
	Id          uint `json:"id"`
	UnDeletable bool `json:"unDeletable"`
}

type UpdateCuisinesSorts struct {
	Cuisines []model.CuisineInfo `json:"cuisines"`
}

type DeleteCuisine struct {
	Id uint `json:"id"`
}

type AddCuisine struct {
	Name string `json:"name"`
}
