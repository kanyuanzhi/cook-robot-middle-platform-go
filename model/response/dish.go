package response

import (
	"github.com/google/uuid"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/model"
)

type CountDishes struct {
	Count int64 `json:"count"`
}

type ListDishes struct {
	Dishes []model.DishInfo `json:"dishes"`
}

type UpdateDishImage struct {
	Image string `json:"image"`
}

type AddDish struct {
	Dish model.DishInfo `json:"dish"`
}

type UpdateDishWithSteps struct {
	Dish model.DishInfo `json:"dish"`
}

type GetDish struct {
	Dish model.DishInfo `json:"dish"`
}

type AddDishCustomSteps struct {
	UUID        uuid.UUID                `json:"uuid"`
	CustomSteps []map[string]interface{} `json:"customSteps"`
}
