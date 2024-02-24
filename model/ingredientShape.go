package model

import "github.com/kanyuanzhi/cook-robot-middle-platform-go/global"

type SysIngredientShape struct {
	global.FXModel
	Name   string `json:"name" gorm:"comment:食材形状;"`
	NameEn string `json:"nameEn" gorm:"comment:食材形状英文名称;"`
	NameTw string `json:"nameTw" gorm:"comment:食材形状繁体名称;"`
}

type IngredientShapeInfo struct {
	Id     uint   `json:"id"`
	Sort   uint   `json:"sort"`
	Name   string `json:"name"`
	NameEn string `json:"nameEn"`
	NameTw string `json:"nameTw"`
}
