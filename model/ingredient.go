package model

import (
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/global"
)

type SysIngredient struct {
	global.FXModel
	Name   string `json:"name" gorm:"comment:食材名称;"`
	NameEn string `json:"nameEn" gorm:"comment:食材英文名称;"`
	NameTw string `json:"nameTw" gorm:"comment:食材繁体名称;"`
	Type   uint   `json:"pump" gorm:"comment:类别（肉类、蔬菜类）"`
}

type IngredientInfo struct {
	Id     uint   `json:"id"`
	Sort   uint   `json:"sort"`
	Name   string `json:"name"`
	NameEn string `json:"nameEn"`
	NameTw string `json:"nameTw"`
	Type   uint   `json:"type"`
}
