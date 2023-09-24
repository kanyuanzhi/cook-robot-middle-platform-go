package request

import (
	"github.com/gin-gonic/gin"
	"github.com/kanyuanzhi/middle-platform/global"
	"github.com/kanyuanzhi/middle-platform/model"
	"gorm.io/gorm"
	"strconv"
)

func ShouldBindJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		//global.GqaLogger.Error(global.GqaConfig.System.BindError, zap.Any("err", err))
		//response.ErrorMessage(c, err.Error())
		return err
	}
	return nil
}

func ShouldBindQuery(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		//global.GqaLogger.Error(global.GqaConfig.System.BindError, zap.Any("err", err))
		//response.ErrorMessage(c, err.Error())
		return err
	}
	return nil
}

func GenerateDishQueryCondition(filter string, enableCuisineFilter bool, cuisineFilter []string, isOfficial bool) (*gorm.DB, error) {
	filterDb := global.FXDb.Model(&model.SysDish{})

	if enableCuisineFilter {
		var cuisineFilterUint []uint
		for _, cuisine := range cuisineFilter {
			cuisineId, _ := strconv.ParseUint(cuisine, 10, 32)
			cuisineFilterUint = append(cuisineFilterUint, uint(cuisineId))
		}
		filterDb = filterDb.Where("cuisine in ?", cuisineFilterUint)
	}

	if filter != "" {
		likeParam := "%" + filter + "%"
		filterDb = filterDb.Where("name LIKE ?", likeParam)
	}

	filterDb = filterDb.Where("is_official", isOfficial)

	return filterDb, nil
}

func GenerateIngredientQueryCondition(enableTypeFilter bool, typeFilter []string) (*gorm.DB, error) {
	filerDb := global.FXDb.Model(&model.SysIngredient{})

	if enableTypeFilter {
		var typeFilterUint []uint
		for _, ingredientType := range typeFilter {
			ingredientTypeId, _ := strconv.ParseUint(ingredientType, 10, 32)
			typeFilterUint = append(typeFilterUint, uint(ingredientTypeId))
		}
		filerDb = filerDb.Where("type in ?", typeFilterUint)
	}

	return filerDb, nil
}
