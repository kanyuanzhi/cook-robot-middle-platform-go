package request

type CountDishes struct {
	Filter              string `json:"filter" form:"filter"`
	EnableCuisineFilter bool   `json:"enableCuisineFilter" form:"enableCuisineFilter"`
	CuisineFilter       string `json:"cuisineFilter" form:"cuisineFilter"`
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
