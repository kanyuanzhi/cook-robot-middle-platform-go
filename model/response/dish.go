package response

import "github.com/kanyuanzhi/middle-platform/model"

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
