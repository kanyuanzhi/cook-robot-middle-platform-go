package private

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/global"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/model"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/model/request"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/model/response"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/utils"
	"log"
	"strconv"
	"strings"
)

type DishApi struct{}

func (api *DishApi) Count(c *gin.Context) {
	var countDishesRequest request.CountDishes
	if err := request.ShouldBindQuery(c, &countDishesRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	filterDb, err := request.GenerateDishQueryCondition(countDishesRequest.Filter, countDishesRequest.EnableCuisineFilter,
		strings.Split(countDishesRequest.CuisineFilter, ","), countDishesRequest.IsOfficial)
	if err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	var count int64
	if err := filterDb.Count(&count).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	countDishesResponse := response.CountDishes{
		Count: count,
	}

	response.SuccessData(c, countDishesResponse)
}

func (api *DishApi) List(c *gin.Context) {
	var listDishesRequest request.ListDishes
	if err := request.ShouldBindQuery(c, &listDishesRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	filterDb, err := request.GenerateDishQueryCondition(listDishesRequest.Filter, listDishesRequest.EnableCuisineFilter,
		strings.Split(listDishesRequest.CuisineFilter, ","), listDishesRequest.IsOfficial)
	if err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	var dishes []model.SysDish
	if err := filterDb.Limit(listDishesRequest.PageSize).Offset((listDishesRequest.PageIndex - 1) * listDishesRequest.PageSize).
		Order("id").Find(&dishes).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	dishesInfo := []model.DishInfo{}
	for _, dish := range dishes {
		dishesInfo = append(dishesInfo, model.DishInfo{
			Id:              dish.Id,
			Image:           "data:image/png;base64," + base64.StdEncoding.EncodeToString(dish.Image),
			Name:            dish.Name,
			UUID:            dish.UUID,
			Steps:           dish.Steps,
			CustomStepsList: dish.CustomStepsList,
			Cuisine:         dish.Cuisine,
		})
	}

	listDishesResponse := response.ListDishes{
		Dishes: dishesInfo,
	}

	response.SuccessData(c, listDishesResponse)
}

// 仅更新名称和所属菜系
func (api *DishApi) Update(c *gin.Context) {
	var updateDishRequest request.UpdateDish
	if err := request.ShouldBindJSON(c, &updateDishRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	dish := model.SysDish{
		FXModel: global.FXModel{
			Id: updateDishRequest.Id,
		},
		Name:    updateDishRequest.Name,
		Cuisine: updateDishRequest.Cuisine,
	}

	if err := global.FXDb.Model(&dish).Select("name", "cuisine").Updates(dish).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	response.SuccessMessage(c, "更新成功")
}

func (api *DishApi) UpdateMark(c *gin.Context) {
	var updateDishMarkRequest request.UpdateDishMark
	if err := request.ShouldBindJSON(c, &updateDishMarkRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	dish := model.SysDish{
		FXModel: global.FXModel{
			Id: updateDishMarkRequest.Id,
		},
		IsMarked: updateDishMarkRequest.Mark,
	}

	if err := global.FXDb.Model(&dish).Select("is_marked").Updates(dish).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	if updateDishMarkRequest.Mark {
		response.SuccessMessage(c, "已添加我的菜品")
	} else {
		response.SuccessMessage(c, "已移除我的菜品")
	}
}

func (api *DishApi) UpdateWithSteps(c *gin.Context) {
	var updateDishWithStepsRequest request.UpdateDishWithSteps
	if err := request.ShouldBindJSON(c, &updateDishWithStepsRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	dish := model.SysDish{
		FXModel: global.FXModel{
			Id: updateDishWithStepsRequest.Id,
		},
		Name:    updateDishWithStepsRequest.Name,
		Cuisine: updateDishWithStepsRequest.Cuisine,
		Steps:   updateDishWithStepsRequest.Steps,
		CustomStepsList: map[string][]map[string]interface{}{
			uuid.New().String(): updateDishWithStepsRequest.Steps,
			uuid.New().String(): updateDishWithStepsRequest.Steps,
			uuid.New().String(): updateDishWithStepsRequest.Steps,
		},
	}

	if err := global.FXDb.Model(&dish).Select("name", "cuisine", "steps", "custom_steps_list").Updates(dish).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	updateDishWithStepsResponse := response.UpdateDishWithSteps{
		Dish: model.DishInfo{
			Id:              dish.Id,
			Image:           "data:image/png;base64," + base64.StdEncoding.EncodeToString(dish.Image),
			Name:            dish.Name,
			UUID:            dish.UUID,
			Steps:           dish.Steps,
			CustomStepsList: dish.CustomStepsList,
			Cuisine:         dish.Cuisine,
		},
	}

	response.SuccessMessageData(c, updateDishWithStepsResponse, "更新成功")
}

func (api *DishApi) Delete(c *gin.Context) {
	var deleteDishesRequest request.DeleteDishes
	if err := request.ShouldBindJSON(c, &deleteDishesRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	var dishes []model.SysDish
	for _, id := range deleteDishesRequest.Ids {
		dishes = append(dishes, model.SysDish{
			FXModel: global.FXModel{
				Id: id,
			},
		})
	}

	if err := global.FXDb.Where("id in ?", deleteDishesRequest.Ids).Select("uuid", "owner").Find(&dishes).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	log.Println(dishes)

	var userDeletedDishes []model.SysUserDeletedDish
	for _, dish := range dishes {
		userDeletedDishes = append(userDeletedDishes, model.SysUserDeletedDish{
			UUID:  dish.UUID,
			Owner: dish.Owner,
		})
	}

	tx := global.FXDb.Begin()

	if err := tx.Create(&userDeletedDishes).Error; err != nil {
		tx.Rollback()
		response.ErrorMessage(c, err.Error())
		return
	}

	if err := tx.Where("id in ?", deleteDishesRequest.Ids).Delete(&dishes).Error; err != nil {
		tx.Rollback()
		response.ErrorMessage(c, err.Error())
		return
	}

	if err := tx.Commit().Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	response.SuccessMessage(c, "删除成功")
}

func (api *DishApi) Add(c *gin.Context) {
	var addDishRequest request.AddDish
	if err := request.ShouldBindJSON(c, &addDishRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	imageData, err := utils.LoadLocalImage("./assets/default_dish_image.png")
	if err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	dish := model.SysDish{
		Name:    addDishRequest.Name,
		Cuisine: addDishRequest.Cuisine,
		UUID:    uuid.New(),
		Steps:   addDishRequest.Steps,
		CustomStepsList: map[string][]map[string]interface{}{
			uuid.New().String(): addDishRequest.Steps,
			uuid.New().String(): addDishRequest.Steps,
			uuid.New().String(): addDishRequest.Steps,
		},
		Image:      imageData,
		IsOfficial: false,
		IsShared:   false,
		Owner:      global.FXSoftwareInfo.SerialNumber,
	}

	if err := global.FXDb.Create(&dish).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	addDishResponse := response.AddDish{
		Dish: model.DishInfo{
			Id:              dish.Id,
			Image:           "data:image/png;base64," + base64.StdEncoding.EncodeToString(dish.Image),
			Name:            dish.Name,
			UUID:            dish.UUID,
			Steps:           dish.Steps,
			CustomStepsList: dish.CustomStepsList,
			Cuisine:         dish.Cuisine,
			IsOfficial:      dish.IsOfficial,
			IsShared:        dish.IsShared,
			Owner:           dish.Owner,
		},
	}

	response.SuccessMessageData(c, addDishResponse, "创建菜品"+dish.Name+"成功")
}

func (api *DishApi) UpdateImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}
	defer src.Close()

	// Read the file content into a byte slice
	fileData := make([]byte, file.Size)
	_, err = src.Read(fileData)
	if err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}
	idStr := c.PostForm("id")

	id, _ := strconv.Atoi(idStr)
	dish := model.SysDish{
		FXModel: global.FXModel{Id: uint(id)},
		Image:   fileData,
	}

	if err := global.FXDb.Model(&dish).Update("image", dish.Image).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	updatedImage := "data:image/png;base64," + base64.StdEncoding.EncodeToString(dish.Image)

	response.SuccessMessageData(c, updatedImage, "更新菜品图片成功")
}

func (api *DishApi) Get(c *gin.Context) {
	var getDishRequest request.GetDish
	if err := request.ShouldBindQuery(c, &getDishRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	var dish model.SysDish
	if err := global.FXDb.Where("uuid = ?", getDishRequest.UUID).First(&dish).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	getDishResponse := response.GetDish{
		Dish: model.DishInfo{
			Id:              dish.Id,
			Image:           "data:image/png;base64," + base64.StdEncoding.EncodeToString(dish.Image),
			Name:            dish.Name,
			UUID:            dish.UUID,
			Steps:           dish.Steps,
			CustomStepsList: dish.CustomStepsList,
			Cuisine:         dish.Cuisine,
			IsOfficial:      dish.IsOfficial,
			IsShared:        dish.IsShared,
			IsMarked:        dish.IsMarked,
		},
	}

	response.SuccessData(c, getDishResponse)
}

func (api *DishApi) UpdateCustomSteps(c *gin.Context) {
	var updateDishCustomStepsRequest request.UpdateDishCustomSteps
	if err := request.ShouldBindJSON(c, &updateDishCustomStepsRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	dish := model.SysDish{
		FXModel:         global.FXModel{Id: updateDishCustomStepsRequest.Id},
		CustomStepsList: updateDishCustomStepsRequest.CustomStepsList,
	}

	customStepsListBytes, _ := json.Marshal(updateDishCustomStepsRequest.CustomStepsList)

	if err := global.FXDb.Model(&dish).Update("custom_steps_list", string(customStepsListBytes)).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	response.SuccessMessage(c, "更新成功")
}

func (api *DishApi) AddCustomSteps(c *gin.Context) {
	var addDishCustomStepsRequest request.AddDishCustomSteps
	if err := request.ShouldBindJSON(c, &addDishCustomStepsRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	dish := model.SysDish{
		FXModel: global.FXModel{Id: addDishCustomStepsRequest.Id},
	}

	if err := global.FXDb.First(&dish).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	customStepsList := dish.CustomStepsList
	if customStepsList == nil {
		customStepsList = make(map[string][]map[string]interface{})
	}
	customUUID := uuid.New()
	customStepsList[customUUID.String()] = dish.Steps

	customStepsListBytes, _ := json.Marshal(customStepsList)

	if err := global.FXDb.Model(&dish).Update("custom_steps_list", string(customStepsListBytes)).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	addDishCustomStepsResponse := response.AddDishCustomSteps{
		CustomSteps: dish.Steps,
		UUID:        customUUID,
	}

	response.SuccessMessageData(c, addDishCustomStepsResponse, "添加成功")
}

func (api *DishApi) DeleteCustomSteps(c *gin.Context) {
	var deleteDishCustomStepsRequest request.DeleteDishCustomSteps
	if err := request.ShouldBindJSON(c, &deleteDishCustomStepsRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	dish := model.SysDish{
		FXModel: global.FXModel{Id: deleteDishCustomStepsRequest.Id},
	}

	if err := global.FXDb.First(&dish).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	customStepsList := dish.CustomStepsList
	delete(customStepsList, deleteDishCustomStepsRequest.UUID.String())

	customStepsListBytes, _ := json.Marshal(customStepsList)

	if err := global.FXDb.Model(&dish).Update("custom_steps_list", string(customStepsListBytes)).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	response.SuccessMessage(c, "删除成功")
}

func (api *DishApi) AddToPersonals(c *gin.Context) {
	var addDishToPersonalsRequest request.AddDishToPersonals
	if err := request.ShouldBindJSON(c, &addDishToPersonalsRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	dish := model.SysDish{
		FXModel: global.FXModel{Id: addDishToPersonalsRequest.Id},
	}

	if err := global.FXDb.First(&dish).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	personalDish := model.SysDish{
		Name:            dish.Name,
		UUID:            uuid.New(),
		Steps:           dish.Steps,
		CustomStepsList: dish.CustomStepsList,
		Image:           dish.Image,
		Cuisine:         dish.Cuisine,
		IsOfficial:      false,
		IsShared:        false,
		IsMarked:        false,
		Owner:           global.FXSoftwareInfo.SerialNumber,
	}

	if err := global.FXDb.Create(&personalDish).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	response.SuccessMessage(c, "已添加至我的菜品")
}
