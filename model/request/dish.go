package request

import "github.com/google/uuid"

type CountDishes struct {
	Filter              string `json:"filter" form:"filter"`
	EnableCuisineFilter bool   `json:"enableCuisineFilter" form:"enableCuisineFilter"`
	CuisineFilter       string `json:"cuisineFilter" form:"cuisineFilter"`
	IsOfficial          bool   `json:"isOfficial" form:"isOfficial"`
}

type ListDishes struct {
	PageIndex int `json:"pageIndex" form:"pageIndex"`
	PageSize  int `json:"pageSize" form:"pageSize"`
	CountDishes
}

type UpdateDish struct {
	Id      uint   `json:"id"`
	Name    string `json:"name"`
	Cuisine uint   `json:"cuisine"`
}

type UpdateDishMark struct {
	Id   uint `json:"id"`
	Mark bool `json:"mark"`
}

type UpdateDishWithSteps struct {
	UpdateDish
	Steps []map[string]interface{} `json:"steps"`
}

type DeleteDishes struct {
	Ids []uint `json:"ids"`
}

type AddDish struct {
	Name    string                   `json:"name"`
	Cuisine uint                     `json:"cuisine"`
	Steps   []map[string]interface{} `json:"steps"`
}

type GetDish struct {
	UUID string `json:"uuid" form:"uuid"`
}

type UpdateDishCustomSteps struct {
	Id              uint                                `json:"id" form:"id"`
	CustomStepsList map[string][]map[string]interface{} `json:"customStepsList" form:"customStepsList"`
}

type AddDishCustomSteps struct {
	Id uint `json:"id" form:"id"`
}

type DeleteDishCustomSteps struct {
	Id   uint      `json:"id" form:"id"`
	UUID uuid.UUID `json:"uuid" form:"uuid"`
}

type AddDishToPersonals struct {
	Id uint `json:"id" form:"id"`
}
