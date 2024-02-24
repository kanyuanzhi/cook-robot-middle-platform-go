package model

import (
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/global"
)

type SysIngredientType struct {
	global.FXModel
	Name        string `json:"name" gorm:"comment:食材种类名称;"`
	NameEn      string `json:"nameEn" gorm:"comment:食材种类英文名称;"`
	NameTw      string `json:"nameTw" gorm:"comment:食材种类繁体名称;"`
	UnDeletable bool   `json:"unDeletable" gorm:"comment:是否允许删除，‘其他’类型不允许删除;default:false"`
}

type IngredientTypeInfo struct {
	Id          uint   `json:"id"`
	Sort        uint   `json:"sort"`
	Name        string `json:"name"`
	NameEn      string `json:"nameEn"`
	NameTw      string `json:"nameTw"`
	UnDeletable bool   `json:"unDeletable"`
}
