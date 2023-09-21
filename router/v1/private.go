package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kanyuanzhi/middle-platform/api/v1/private"
)

func InitPrivateRouter(routerGroup *gin.RouterGroup) {
	cuisineApi := &private.CuisineApi{}
	routerGroup.GET("cuisine/list", cuisineApi.List)
	//routerGroup.PUT("cuisine/update-name", cuisineApi.UpdateName)
	//routerGroup.PUT("cuisine/update-unDeletable", cuisineApi.UpdateUnDeletable)
	//routerGroup.PUT("cuisine/update-sorts", cuisineApi.UpdateSorts)
	//routerGroup.DELETE("cuisine/delete", cuisineApi.Delete)
	//routerGroup.POST("cuisine/add", cuisineApi.Add)

	dishApi := &private.DishApi{}
	routerGroup.GET("dish/list", dishApi.List)   //
	routerGroup.GET("dish/count", dishApi.Count) //
	routerGroup.GET("dish/get", dishApi.Get)     //
	//routerGroup.PUT("dish/update", dishApi.Update)
	routerGroup.DELETE("dish/delete", dishApi.Delete) //
	//routerGroup.POST("dish/update-image", dishApi.UpdateImage)
	routerGroup.POST("dish/add", dishApi.Add)                                //
	routerGroup.PUT("dish/update-with-steps", dishApi.UpdateWithSteps)       //
	routerGroup.PUT("dish/update-customSteps", dishApi.UpdateCustomSteps)    //
	routerGroup.PUT("dish/add-customSteps", dishApi.AddCustomSteps)          //
	routerGroup.DELETE("dish/delete-customSteps", dishApi.DeleteCustomSteps) //

	seasoningApi := &private.SeasoningApi{}
	routerGroup.GET("seasoning/list", seasoningApi.List) //
	//routerGroup.PUT("seasoning/update", seasoningApi.Update)
	//routerGroup.DELETE("seasoning/delete", seasoningApi.Delete)
	//routerGroup.POST("seasoning/add", seasoningApi.Add)
	routerGroup.PUT("seasoning/update-pumpRatios", seasoningApi.UpdatePumpRatios) //

	ingredientApi := &private.IngredientApi{}
	routerGroup.GET("ingredient/list", ingredientApi.List) //
	//routerGroup.GET("ingredient/count", ingredientApi.Count)
	//routerGroup.POST("ingredient/add", ingredientApi.Add)
	//routerGroup.PUT("ingredient/update", ingredientApi.Update)
	//routerGroup.DELETE("ingredient/delete", ingredientApi.Delete)

	ingredientTypeApi := &private.IngredientTypeApi{}
	routerGroup.GET("ingredient-type/list", ingredientTypeApi.List) //
	//routerGroup.GET("ingredient-type/count", ingredientTypeApi.Count)
	//routerGroup.POST("ingredient-type/add", ingredientTypeApi.Add)
	//routerGroup.PUT("ingredient-type/update", ingredientTypeApi.Update)
	//routerGroup.PUT("ingredient-type/update-sorts", ingredientTypeApi.UpdateSorts)
	//routerGroup.DELETE("ingredient-type/delete", ingredientTypeApi.Delete)

	ingredientShapeApi := &private.IngredientShapeApi{}
	routerGroup.GET("ingredient-shape/list", ingredientShapeApi.List) //
	//routerGroup.GET("ingredient-shape/count", ingredientShapeApi.Count)
	//routerGroup.POST("ingredient-shape/add", ingredientShapeApi.Add)
	//routerGroup.PUT("ingredient-shape/update", ingredientShapeApi.Update)
	//routerGroup.PUT("ingredient-shape/update-sorts", ingredientShapeApi.UpdateSorts)
	//routerGroup.DELETE("ingredient-shape/delete", ingredientShapeApi.Delete)

	controllerApi := &private.ControllerApi{}
	routerGroup.POST("controller/execute", controllerApi.Execute)
	routerGroup.GET("controller/fetch-status", controllerApi.FetchStatus)

	systemApi := &private.SystemApi{}
	routerGroup.GET("system/get-qrCode", systemApi.GetQrCode)
	routerGroup.GET("system/shutdown", systemApi.Shutdown)

	softwareUpdaterApi := &private.SoftwareUpdaterApi{
		IsUpdating:     false,
		LatestVersion:  "",
		UpdateFilePath: []string{},
	}
	routerGroup.GET("softwareUpdater/get-softwareInfo", softwareUpdaterApi.GetSoftwareInfo)
	routerGroup.GET("softwareUpdater/check-updateInfo", softwareUpdaterApi.CheckUpdateInfo)
	routerGroup.GET("softwareUpdater/check-updatePermission", softwareUpdaterApi.CheckUpdatePermission)
	routerGroup.GET("softwareUpdater/update", softwareUpdaterApi.Update)
}
