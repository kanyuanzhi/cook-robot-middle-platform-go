package request

type CountIngredients struct {
	EnableTypeFilter bool   `json:"enableTypeFilter" form:"enableTypeFilter"`
	TypeFilter       string `json:"typeFilter" form:"typeFilter"`
}

type ListIngredients struct {
	PageIndex int `json:"pageIndex" form:"pageIndex"`
	PageSize  int `json:"pageSize" form:"pageSize"`
	CountIngredients
}

type AddIngredient struct {
	Name string `json:"name"`
}

type UpdateIngredient struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
	Type uint   `json:"type"`
}

type DeleteIngredient struct {
	Id uint `json:"id"`
}
