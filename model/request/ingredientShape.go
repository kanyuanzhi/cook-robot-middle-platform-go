package request

import "github.com/kanyuanzhi/cook-robot-middle-platform-go/model"

type AddIngredientShape struct {
	Name string `json:"name"`
}

type UpdateIngredientShape struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type UpdateIngredientShapeSorts struct {
	IngredientShapes []model.IngredientShapeInfo `json:"ingredientShapes"`
}

type DeleteIngredientShape struct {
	Id uint `json:"id"`
}
