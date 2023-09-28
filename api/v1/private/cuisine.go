package private

import (
	"github.com/gin-gonic/gin"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/global"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/model"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/model/request"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/model/response"
)

type CuisineApi struct {
}

func (api *CuisineApi) List(c *gin.Context) {
	var cuisines []model.SysCuisine
	if err := global.FXDb.Order("sort").Find(&cuisines).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	cuisinesInfo := []model.CuisineInfo{}
	for _, cuisine := range cuisines {
		cuisinesInfo = append(cuisinesInfo, model.CuisineInfo{
			Id:          cuisine.Id,
			Sort:        cuisine.Sort,
			Name:        cuisine.Name,
			UnDeletable: cuisine.UnDeletable,
		})
	}

	listCuisinesResponse := response.ListCuisines{
		Cuisines: cuisinesInfo,
	}

	response.SuccessData(c, listCuisinesResponse)
}

func (api *CuisineApi) UpdateName(c *gin.Context) {
	var updateCuisineNameRequest request.UpdateCuisineName
	if err := request.ShouldBindJSON(c, &updateCuisineNameRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	cuisine := model.SysCuisine{
		FXModel: global.FXModel{
			Id: updateCuisineNameRequest.Id,
		},
		Name: updateCuisineNameRequest.Name,
	}

	if err := global.FXDb.Model(&cuisine).Update("name", cuisine.Name).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	response.SuccessMessage(c, "更新名称成功")
}

func (api *CuisineApi) UpdateUnDeletable(c *gin.Context) {
	var updateCuisineUnDeletableRequest request.UpdateCuisineUnDeletable
	if err := request.ShouldBindJSON(c, &updateCuisineUnDeletableRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	cuisine := model.SysCuisine{
		FXModel: global.FXModel{
			Id: updateCuisineUnDeletableRequest.Id,
		},
		UnDeletable: updateCuisineUnDeletableRequest.UnDeletable,
	}

	if err := global.FXDb.Model(&cuisine).Update("un_deletable", cuisine.UnDeletable).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	response.SuccessMessage(c, "更新不可删除属性成功")
}

func (api *CuisineApi) UpdateSorts(c *gin.Context) {
	var updateCuisinesSortsRequest request.UpdateCuisinesSorts
	if err := request.ShouldBindJSON(c, &updateCuisinesSortsRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	tx := global.FXDb.Begin()
	for _, cuisine := range updateCuisinesSortsRequest.Cuisines {
		sysCuisine := model.SysCuisine{
			FXModel: global.FXModel{
				Id:   cuisine.Id,
				Sort: cuisine.Sort,
			},
		}
		if err := tx.Model(&sysCuisine).Update("sort", sysCuisine.Sort).Error; err != nil {
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

func (api *CuisineApi) Delete(c *gin.Context) {
	// 删除菜系的同时，需要将属于该菜系的菜品划归到‘其他’类别下
	var deleteCuisineRequest request.DeleteCuisine
	if err := request.ShouldBindJSON(c, &deleteCuisineRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	cuisine := model.SysCuisine{
		FXModel: global.FXModel{
			Id: deleteCuisineRequest.Id,
		},
	}
	if err := global.FXDb.First(&cuisine).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	if cuisine.UnDeletable {
		response.ErrorMessage(c, "不允许删除该类别菜系")
		return
	}

	var unDeletableCuisine model.SysCuisine
	result := global.FXDb.Where("un_deletable", true).First(&unDeletableCuisine)
	if result.Error != nil {
		response.ErrorMessage(c, result.Error.Error()+"：当前未设置不允许删除的菜系，无法删除")
		return
	}

	// 划归菜品到‘其他’类别下
	if err := global.FXDb.Model(model.SysDish{}).Where("cuisine", cuisine.Id).
		Updates(model.SysDish{Cuisine: unDeletableCuisine.Id}).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	// 删除该菜系
	if err := global.FXDb.Delete(&cuisine).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	response.SuccessMessage(c, "删除成功")
}

func (api *CuisineApi) Add(c *gin.Context) {
	var addCuisineRequest request.AddCuisine
	if err := request.ShouldBindJSON(c, &addCuisineRequest); err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	cuisine := model.SysCuisine{
		FXModel: global.FXModel{
			Sort: 0,
		},
		Name:        addCuisineRequest.Name,
		UnDeletable: true,
	}

	if err := global.FXDb.Create(&cuisine).Error; err != nil {
		response.ErrorMessage(c, err.Error())
		return
	}

	addCuisineResponse := response.AddCuisine{
		Cuisine: model.CuisineInfo{
			Id:          cuisine.Id,
			Sort:        cuisine.Sort,
			Name:        cuisine.Name,
			UnDeletable: cuisine.UnDeletable,
		},
	}

	response.SuccessMessageData(c, addCuisineResponse, "添加菜系成功")
}
