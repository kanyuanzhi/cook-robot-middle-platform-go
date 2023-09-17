package private

import (
	"github.com/gin-gonic/gin"
	"github.com/kanyuanzhi/middle-platform/global"
	"github.com/kanyuanzhi/middle-platform/model"
	"github.com/kanyuanzhi/middle-platform/model/request"
	"github.com/kanyuanzhi/middle-platform/model/response"
	"strings"
)

type IngredientApi struct{}

func (api *IngredientApi) Count(c *gin.Context) {
	var countIngredientsRequest request.CountIngredients
	if err := request.ShouldBindQuery(c, &countIngredientsRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	filterDb, err := request.GenerateIngredientQueryCondition(countIngredientsRequest.EnableTypeFilter,
		strings.Split(countIngredientsRequest.TypeFilter, ","))
	if err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	var count int64
	if err := filterDb.Count(&count).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	countIngredientsResponse := response.CountIngredients{
		Count: count,
	}

	response.SuccessData(c, countIngredientsResponse)
}

func (api *IngredientApi) List(c *gin.Context) {
	var listIngredientsRequest request.ListIngredients
	if err := request.ShouldBindQuery(c, &listIngredientsRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	filterDb, err := request.GenerateIngredientQueryCondition(listIngredientsRequest.EnableTypeFilter,
		strings.Split(listIngredientsRequest.TypeFilter, ","))
	if err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	var ingredients []model.SysIngredient
	if err := filterDb.Order("sort").Find(&ingredients).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	ingredientsInfo := []model.IngredientInfo{}
	for _, ingredient := range ingredients {
		ingredientsInfo = append(ingredientsInfo, model.IngredientInfo{
			Id:   ingredient.Id,
			Sort: ingredient.Sort,
			Name: ingredient.Name,
			Type: ingredient.Type,
		})
	}

	listIngredientsResponse := response.ListIngredients{
		Ingredients: ingredientsInfo,
	}

	response.SuccessData(c, listIngredientsResponse)
}

func (api *IngredientApi) Add(c *gin.Context) {
	var addIngredientRequest request.AddIngredient
	if err := request.ShouldBindJSON(c, &addIngredientRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	ingredient := model.SysIngredient{
		FXModel: global.FXModel{
			Sort: 0,
		},
		Name: addIngredientRequest.Name,
		Type: 0,
	}

	if err := global.FXDb.Create(&ingredient).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	addIngredientResponse := response.AddIngredient{
		Ingredient: model.IngredientInfo{
			Id:   ingredient.Id,
			Sort: ingredient.Sort,
			Name: ingredient.Name,
			Type: ingredient.Type,
		},
	}

	response.SuccessMessageData(c, addIngredientResponse, "添加成功")
}

func (api *IngredientApi) Update(c *gin.Context) {
	var updateIngredientRequest request.UpdateIngredient
	if err := request.ShouldBindJSON(c, &updateIngredientRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	ingredient := model.SysIngredient{
		FXModel: global.FXModel{
			Id: updateIngredientRequest.Id,
		},
		Name: updateIngredientRequest.Name,
		Type: updateIngredientRequest.Type,
	}

	if err := global.FXDb.Model(&ingredient).Select("name", "type").Updates(ingredient).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	response.SuccessMessage(c, "更新成功")
}

func (api *IngredientApi) Delete(c *gin.Context) {
	var deleteIngredientRequest request.DeleteIngredient
	if err := request.ShouldBindJSON(c, &deleteIngredientRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	ingredient := model.SysIngredient{
		FXModel: global.FXModel{
			Id: deleteIngredientRequest.Id,
		},
	}

	if err := global.FXDb.Delete(&ingredient).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	response.SuccessMessage(c, "删除成功")
}
