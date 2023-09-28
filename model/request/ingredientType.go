package request

import "github.com/kanyuanzhi/cook-robot-middle-platform-go/model"

type AddIngredientType struct {
	Name string `json:"name"`
}

type UpdateIngredientType struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	UnDeletable bool   `json:"unDeletable"`
}

type UpdateIngredientTypeSorts struct {
	IngredientTypes []model.IngredientTypeInfo `json:"ingredientTypes"`
}

type DeleteIngredientType struct {
	Id uint `json:"id"`
}
