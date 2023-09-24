package response

type UpdateOfficialDishes struct {
	NewAddedDishesNumber int `json:"newAddedDishesNumber"`
	UpdatesDishesNumber  int `json:"updatesDishesNumber"`
	DeletedDishesNumber  int `json:"deletedDishesNumber"`
}

type SynchronizePersonalDishes struct {
	RemoteNeedAddDishesNumber    int `json:"remoteNeedAddDishesNumber"`
	RemoteNeedUpdateDishesNumber int `json:"remoteNeedUpdateDishesNumber"`
	RemoteNeedDeleteDishesNumber int `json:"remoteNeedDeleteDishesNumber"`
	LocalNeedAddDishesNumber     int `json:"localNeedAddDishesNumber"`
	LocalNeedUpdateDishesNumber  int `json:"localNeedUpdateDishesNumber"`
	LocalNeedDeleteDishesNumber  int `json:"localNeedDeleteDishesNumber"`
}
