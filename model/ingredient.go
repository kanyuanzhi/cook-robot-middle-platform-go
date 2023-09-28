package model

import (
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/global"
)

type SysIngredient struct {
	global.FXModel
	Name string `json:"name" gorm:"comment:食材名称;"`
	Type uint   `json:"pump" gorm:"comment:类别（肉类、蔬菜类）"`
}

type IngredientInfo struct {
	Id   uint   `json:"id"`
	Sort uint   `json:"sort"`
	Name string `json:"name"`
	Type uint   `json:"type"`
}
