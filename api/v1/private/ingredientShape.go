package private

import (
	"github.com/gin-gonic/gin"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/global"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/model"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/model/request"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/model/response"
)

type IngredientShapeApi struct{}

func (api *IngredientShapeApi) Count(c *gin.Context) {
	var count int64
	if err := global.FXDb.Model(&model.SysIngredientShape{}).Count(&count).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	countIngredientShapesResponse := response.CountIngredientShapes{
		Count: count,
	}

	response.SuccessData(c, countIngredientShapesResponse)
}

func (api *IngredientShapeApi) List(c *gin.Context) {
	var ingredientShapes []model.SysIngredientShape
	if err := global.FXDb.Order("sort").Find(&ingredientShapes).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	ingredientShapesInfo := []model.IngredientShapeInfo{}
	for _, ingredientShape := range ingredientShapes {
		ingredientShapesInfo = append(ingredientShapesInfo, model.IngredientShapeInfo{
			Id:   ingredientShape.Id,
			Sort: ingredientShape.Sort,
			Name: ingredientShape.Name,
		})
	}

	listIngredientShapesResponse := response.ListIngredientShapes{
		IngredientShapes: ingredientShapesInfo,
	}

	response.SuccessData(c, listIngredientShapesResponse)
}

func (api *IngredientShapeApi) Add(c *gin.Context) {
	var addIngredientShapeRequest request.AddIngredientShape
	if err := request.ShouldBindJSON(c, &addIngredientShapeRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	ingredientShape := model.SysIngredientShape{
		FXModel: global.FXModel{
			Sort: 0,
		},
		Name: addIngredientShapeRequest.Name,
	}

	if err := global.FXDb.Create(&ingredientShape).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	addIngredientShapeResponse := response.AddIngredientShape{
		IngredientShape: model.IngredientShapeInfo{
			Id:   ingredientShape.Id,
			Sort: ingredientShape.Sort,
			Name: ingredientShape.Name,
		},
	}

	response.SuccessMessageData(c, addIngredientShapeResponse, "添加成功")
}

func (api *IngredientShapeApi) Update(c *gin.Context) {
	var updateIngredientShapeRequest request.UpdateIngredientShape
	if err := request.ShouldBindJSON(c, &updateIngredientShapeRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	ingredientShape := model.SysIngredientShape{
		FXModel: global.FXModel{
			Id: updateIngredientShapeRequest.Id,
		},
		Name: updateIngredientShapeRequest.Name,
	}

	if err := global.FXDb.Model(&ingredientShape).Select("name", "un_deletable").Updates(ingredientShape).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	response.SuccessMessage(c, "更新成功")
}

func (api *IngredientShapeApi) UpdateSorts(c *gin.Context) {
	var updateIngredientShapeSortsRequest request.UpdateIngredientShapeSorts
	if err := request.ShouldBindJSON(c, &updateIngredientShapeSortsRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	tx := global.FXDb.Begin()
	for _, ingredientShape := range updateIngredientShapeSortsRequest.IngredientShapes {
		sysIngredientShape := model.SysIngredientShape{
			FXModel: global.FXModel{
				Id:   ingredientShape.Id,
				Sort: ingredientShape.Sort,
			},
		}
		if err := tx.Model(&sysIngredientShape).Select("sort").Updates(sysIngredientShape).Error; err != nil {
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

func (api *IngredientShapeApi) Delete(c *gin.Context) {
	var deleteIngredientShapeRequest request.DeleteIngredientShape
	if err := request.ShouldBindJSON(c, &deleteIngredientShapeRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	ingredientShape := model.SysIngredientShape{
		FXModel: global.FXModel{
			Id: deleteIngredientShapeRequest.Id,
		},
	}

	if err := global.FXDb.Delete(&ingredientShape).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	response.SuccessMessage(c, "删除成功")
}
