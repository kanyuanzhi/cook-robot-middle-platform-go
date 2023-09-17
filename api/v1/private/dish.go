package private

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kanyuanzhi/middle-platform/global"
	"github.com/kanyuanzhi/middle-platform/model"
	"github.com/kanyuanzhi/middle-platform/model/request"
	"github.com/kanyuanzhi/middle-platform/model/response"
	"github.com/kanyuanzhi/middle-platform/utils"
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
		strings.Split(countDishesRequest.CuisineFilter, ","))
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
		strings.Split(listDishesRequest.CuisineFilter, ","))
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
			Id:          dish.Id,
			Image:       "data:image/png;base64," + base64.StdEncoding.EncodeToString(dish.Image),
			Name:        dish.Name,
			UUID:        dish.UUID,
			Steps:       dish.Steps,
			CustomSteps: dish.CustomSteps,
			Cuisine:     dish.Cuisine,
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
	}

	if err := global.FXDb.Model(&dish).Select("name", "cuisine", "steps").Updates(dish).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	updateDishWithStepsResponse := response.UpdateDishWithSteps{
		Dish: model.DishInfo{
			Id:          dish.Id,
			Image:       "data:image/png;base64," + base64.StdEncoding.EncodeToString(dish.Image),
			Name:        dish.Name,
			UUID:        dish.UUID,
			Steps:       dish.Steps,
			CustomSteps: dish.CustomSteps,
			Cuisine:     dish.Cuisine,
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

	if err := global.FXDb.Delete(&dishes).Error; err != nil {
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
		CustomSteps: map[string][]map[string]interface{}{
			uuid.New().String(): addDishRequest.Steps,
			uuid.New().String(): addDishRequest.Steps,
			uuid.New().String(): addDishRequest.Steps,
		},
		Image: imageData,
	}

	if err := global.FXDb.Create(&dish).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	addDishResponse := response.AddDish{
		Dish: model.DishInfo{
			Id:          dish.Id,
			Image:       "data:image/png;base64," + base64.StdEncoding.EncodeToString(dish.Image),
			Name:        dish.Name,
			UUID:        dish.UUID,
			Steps:       dish.Steps,
			CustomSteps: dish.CustomSteps,
			Cuisine:     dish.Cuisine,
		},
	}

	response.SuccessMessageData(c, addDishResponse, "添加成功")
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

	response.SuccessMessageData(c, updatedImage, "更新菜谱图片成功")
}
