package response

import "github.com/kanyuanzhi/cook-robot-middle-platform-go/model"

type CountIngredientTypes struct {
	Count int64 `json:"count"`
}

type ListIngredientTypes struct {
	IngredientTypes []model.IngredientTypeInfo `json:"ingredientTypes"`
}

type AddIngredientType struct {
	IngredientType model.IngredientTypeInfo `json:"ingredientType"`
}
