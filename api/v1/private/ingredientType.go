package private

import (
	"github.com/gin-gonic/gin"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/global"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/model"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/model/request"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/model/response"
)

type IngredientTypeApi struct{}

func (api *IngredientTypeApi) Count(c *gin.Context) {
	var count int64
	if err := global.FXDb.Model(&model.SysIngredientType{}).Count(&count).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	countIngredientTypesResponse := response.CountIngredientTypes{
		Count: count,
	}

	response.SuccessData(c, countIngredientTypesResponse)
}

func (api *IngredientTypeApi) List(c *gin.Context) {
	var ingredientTypes []model.SysIngredientType
	if err := global.FXDb.Order("sort").Find(&ingredientTypes).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	ingredientTypesInfo := []model.IngredientTypeInfo{}
	for _, ingredientType := range ingredientTypes {
		ingredientTypesInfo = append(ingredientTypesInfo, model.IngredientTypeInfo{
			Id:          ingredientType.Id,
			Sort:        ingredientType.Sort,
			Name:        ingredientType.Name,
			NameEn:      ingredientType.NameEn,
			NameTw:      ingredientType.NameTw,
			UnDeletable: ingredientType.UnDeletable,
		})
	}

	listIngredientTypesResponse := response.ListIngredientTypes{
		IngredientTypes: ingredientTypesInfo,
	}

	response.SuccessData(c, listIngredientTypesResponse)
}

func (api *IngredientTypeApi) Add(c *gin.Context) {
	var addIngredientTypeRequest request.AddIngredientType
	if err := request.ShouldBindJSON(c, &addIngredientTypeRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	ingredientType := model.SysIngredientType{
		FXModel: global.FXModel{
			Sort: 0,
		},
		Name:        addIngredientTypeRequest.Name,
		UnDeletable: true,
	}

	if err := global.FXDb.Create(&ingredientType).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	addIngredientTypeResponse := response.AddIngredientType{
		IngredientType: model.IngredientTypeInfo{
			Id:          ingredientType.Id,
			Sort:        ingredientType.Sort,
			Name:        ingredientType.Name,
			UnDeletable: ingredientType.UnDeletable,
		},
	}

	response.SuccessMessageData(c, addIngredientTypeResponse, "添加成功")
}

func (api *IngredientTypeApi) Update(c *gin.Context) {
	var updateIngredientTypeRequest request.UpdateIngredientType
	if err := request.ShouldBindJSON(c, &updateIngredientTypeRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	ingredientType := model.SysIngredientType{
		FXModel: global.FXModel{
			Id: updateIngredientTypeRequest.Id,
		},
		Name:        updateIngredientTypeRequest.Name,
		UnDeletable: updateIngredientTypeRequest.UnDeletable,
	}

	if err := global.FXDb.Model(&ingredientType).Select("name", "un_deletable").Updates(ingredientType).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	response.SuccessMessage(c, "更新成功")
}

func (api *IngredientTypeApi) UpdateSorts(c *gin.Context) {
	var updateIngredientTypeSortsRequest request.UpdateIngredientTypeSorts
	if err := request.ShouldBindJSON(c, &updateIngredientTypeSortsRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	tx := global.FXDb.Begin()
	for _, ingredientType := range updateIngredientTypeSortsRequest.IngredientTypes {
		sysIngredientType := model.SysIngredientType{
			FXModel: global.FXModel{
				Id:   ingredientType.Id,
				Sort: ingredientType.Sort,
			},
		}
		if err := tx.Model(&sysIngredientType).Select("sort").Updates(sysIngredientType).Error; err != nil {
			tx.Rollback()
			response.ErrorMessage(c, err.Error())
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	response.SuccessMessage(c, "更新顺序成功")
}

func (api *IngredientTypeApi) Delete(c *gin.Context) {
	var deleteIngredientTypeRequest request.DeleteIngredientType
	if err := request.ShouldBindJSON(c, &deleteIngredientTypeRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	ingredientType := model.SysIngredientType{
		FXModel: global.FXModel{
			Id: deleteIngredientTypeRequest.Id,
		},
	}
	if err := global.FXDb.First(&ingredientType).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	if ingredientType.UnDeletable {
		response.ErrorMessage(c, "不允许删除该食材类别")
		return
	}

	var unDeletableIngredientType model.SysIngredientType
	result := global.FXDb.Where("un_deletable", true).First(&unDeletableIngredientType)
	if result.Error != nil {
		response.ErrorMessage(c, result.Error.Error()+"：当前未设置不允许删除的食材类别，无法删除")
		return
	}

	// 划归食材到‘其他’类别下
	if err := global.FXDb.Model(model.SysIngredient{}).Where("type", ingredientType.Id).
		Updates(model.SysIngredient{Type: unDeletableIngredientType.Id}).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	// 删除该食材类别
	if err := global.FXDb.Delete(&ingredientType).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	response.SuccessMessage(c, "删除成功")
}
